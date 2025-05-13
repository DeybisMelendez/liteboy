package ppu

import (
	"github.com/deybismelendez/liteboy/bus"
)

const (
	LCDCRegister = 0xFF40
	STATRegister = 0xFF41
	SCYRegister  = 0xFF42
	SCXRegister  = 0xFF43
	LYRegister   = 0xFF44
	LYCRegister  = 0xFF45

	ScreenWidth  = 160
	ScreenHeight = 144

	LCDCFlagBGDisplay     = 1 << 0
	LCDCFlagOBJDisplay    = 1 << 1
	LCDCFlagOBJSize       = 1 << 2
	LCDCFlagBGTileMap     = 1 << 3
	LCDCFlagBGTileData    = 1 << 4
	LCDCFlagWindowEnable  = 1 << 5
	LCDCFlagWindowTileMap = 1 << 6
	LCDCFlagLCDEnable     = 1 << 7
)

type PPU struct {
	//lcdc                 LCDC
	//stat                 STAT
	bus                  *bus.Bus
	Framebuffer          []uint32
	cycles               int
	spritesOnCurrentLine []*Sprite
	pixelFIFO            []uint32 // FIFO para los píxeles
	fifoSize             int
}

func NewPPU(b *bus.Bus) *PPU {
	return &PPU{
		bus:         b,
		Framebuffer: make([]uint32, ScreenWidth*ScreenHeight),
		pixelFIFO:   make([]uint32, 0, ScreenWidth),
		fifoSize:    ScreenWidth}
}

// Función para agregar un píxel a la FIFO
func (ppu *PPU) addPixelToFIFO(pixel uint32) {
	if len(ppu.pixelFIFO) < ppu.fifoSize {
		ppu.pixelFIFO = append(ppu.pixelFIFO, pixel)
		//log.Printf("[PPU] Pixel agregado a FIFO (len=%d)", len(ppu.pixelFIFO))
	}
}

// Función para extraer un píxel de la FIFO y colocarlo en el framebuffer
func (ppu *PPU) popPixelFromFIFO(x, y int) {
	if len(ppu.pixelFIFO) > 0 {
		pixel := ppu.pixelFIFO[0]
		ppu.pixelFIFO = ppu.pixelFIFO[1:] // Elimina el primer píxel de la cola
		ppu.Framebuffer[getFramebufferIndex(x, y)] = pixel
		//log.Printf("[PPU] Pixel extraído de FIFO y escrito en (%d, %d)", x, y)
	}
}

func (ppu *PPU) setMode(mode byte) {
	ppu.bus.Write(STATRegister, (ppu.bus.Read(STATRegister)&^0x03)|mode&0x03)
}

func (ppu *PPU) getMode() byte {
	return ppu.bus.Read(STATRegister) & 0x03
}

func (ppu *PPU) Step(tCycles int) {
	if !ppu.isLCDEnabled() {
		ppu.bus.Write(LYRegister, 0)
		ppu.setMode(ModeHBlank)
		return
	}

	ppu.cycles += tCycles

	switch ppu.getMode() {
	case ModeOAM:
		//log.Printf("[PPU] Modo OAM - Ciclos acumulados: %d", ppu.cycles) // LOG
		ppu.scanOAM()
	case ModeVRAM:
		//log.Printf("[PPU] Modo VRAM (Modo 3) - Ciclos acumulados: %d", ppu.cycles) // LOG
		ppu.runVRAM()
	case ModeHBlank:
		//log.Printf("[PPU] Modo HBlank - LY: %d", ppu.bus.Read(LYRegister)) // LOG
		ppu.runHBlank()
	case ModeVBlank:
		//log.Printf("[PPU] Modo VBlank - LY: %d", ppu.bus.Read(LYRegister)) // LOG
		ppu.runVBlank()
	}
}

func (ppu *PPU) scanOAM() {
	if ppu.cycles >= 80 {
		//log.Printf("[PPU] Escaneando OAM - LY: %d", ppu.bus.Read(LYRegister)) // LOG
		spriteHeight := byte(8)
		if !ppu.isObj8x8() {
			spriteHeight = 16
		}
		var result []*Sprite
		ly := ppu.bus.Read(LYRegister)
		for i := uint16(0); i < 40; i++ {
			var index uint16 = i * 4
			y := ppu.bus.Read(0xFE00 + index)        // OAM Y
			x := ppu.bus.Read(0xFE00 + index + 1)    // OAM X
			tile := ppu.bus.Read(0xFE00 + index + 2) // Tile ID
			attr := ppu.bus.Read(0xFE00 + index + 3) // Attributes

			// Posición real de y es y - 16
			if ly >= y-16 && ly < (y-16)+spriteHeight {
				sprite := newSprite(x, y, tile, attr, i)
				result = append(result, sprite)

				if len(result) == MaxSpritesPerLine {
					break
				}
			}
		}
		//log.Printf("[PPU] Sprites en la línea actual: %d", len(result)) // LOG
		ppu.spritesOnCurrentLine = result
		ppu.cycles -= 80
		ppu.setMode(ModeVRAM)
	}
}

func (ppu *PPU) runHBlank() {
	if ppu.isHBlankInterruptEnabled() {
		//TODO: Request Interrupt HBlank
		ppu.requestInterrupt(InterruptSTAT)
	}
	if ppu.cycles >= 204 {

		if ppu.bus.Read(LYRegister) == 144 {
			ppu.setMode(ModeVBlank)
			ppu.requestInterrupt(InterruptVBlank)
			ppu.requestInterrupt(InterruptSTAT)
		} else {
			ppu.setMode(ModeOAM)
		}
		ppu.updateCoincidenceFlag()
		ppu.bus.Write(LYRegister, ppu.bus.Read(LYRegister)+1)
		ppu.cycles -= 204
	}
}

func (ppu *PPU) runVBlank() {
	if ppu.isVBlankInterruptEnabled() {
		ppu.requestInterrupt(InterruptSTAT)
	}
	if ppu.cycles >= 456 {
		ppu.cycles -= 456
		if ppu.bus.Read(LYRegister) == 154 {
			ppu.bus.Write(LYRegister, 0)
			ppu.setCoincidenceFlag(ppu.bus.Read(LYRegister) == ppu.bus.Read(LYCRegister))
			ppu.setMode(ModeOAM)
		} else {
			ppu.updateCoincidenceFlag()
			ppu.bus.Write(LYRegister, ppu.bus.Read(LYRegister)+1)
		}
	}
}

func (ppu *PPU) updateCoincidenceFlag() {
	ly := ppu.bus.Read(LYRegister)
	lyc := ppu.bus.Read(LYCRegister)
	match := ly == lyc
	ppu.setCoincidenceFlag(match)
	if match && ppu.isLYCInterruptEnabled() {
		ppu.requestInterrupt(InterruptSTAT)
	}
}
func (ppu *PPU) runVRAM() {
	// Calculamos el número de ciclos basados en la posición de SCX y los sprites
	var baseCycles = 172                              // valor base para el Modo 3
	spriteCycles := len(ppu.spritesOnCurrentLine) * 2 // Cada sprite puede requerir más ciclos para la transferencia

	if ppu.cycles < baseCycles+spriteCycles {
		return
	}

	// Procedemos con el renderizado de VRAM...
	ly := ppu.bus.Read(LYRegister)
	scx := ppu.bus.Read(SCXRegister)
	scy := ppu.bus.Read(SCYRegister)
	lcdc := ppu.bus.Read(LCDCRegister)

	if int(ly) >= ScreenHeight {
		//log.Printf("[PPU] Línea fuera de rango: LY=%d", ly)
		ppu.cycles -= baseCycles + spriteCycles
		ppu.setMode(ModeHBlank)
		return
	}
	//log.Printf("[PPU] Renderizando línea LY=%d", ly)

	bgTileMapAddr := uint16(0x9800)
	if lcdc&LCDCFlagBGTileMap != 0 {
		bgTileMapAddr = 0x9C00
	}

	tileDataAddr := uint16(0x8800)
	useSigned := true
	if lcdc&LCDCFlagBGTileData != 0 {
		tileDataAddr = 0x8000
		useSigned = false
	}

	for x := 0; x < ScreenWidth; x++ {
		scrollX := uint16((uint16(x) + uint16(scx)) & 0xFF)
		scrollY := uint16((uint16(ly) + uint16(scy)) & 0xFF)

		tileX := scrollX / 8
		tileY := scrollY / 8
		tileIndexOffset := tileY*32 + tileX
		tileIndex := ppu.bus.Read(bgTileMapAddr + tileIndexOffset)

		var tileAddr uint16
		if useSigned {
			tileAddr = tileDataAddr + uint16(int8(tileIndex))*16
		} else {
			tileAddr = tileDataAddr + uint16(tileIndex)*16
		}

		row := (scrollY % 8) * 2
		byte1 := ppu.bus.Read(tileAddr + uint16(row))
		byte2 := ppu.bus.Read(tileAddr + uint16(row) + 1)
		bit := 7 - (scrollX % 8)
		//log.Printf("[PPU] ColorID: %02X, %02X bit %04X", byte1, byte2, bit)

		colorID := ((byte2 >> bit) & 1 << 1) | ((byte1 >> bit) & 1)

		palette := ppu.bus.Read(0xFF47)
		color := (palette >> (colorID * 2)) & 0x03
		//log.Printf("[PPU] ColorID: %d, Color: %02X (BG Palette: %02X)", colorID, color, palette)
		ppu.addPixelToFIFO(getColorFromPalette(color))
		//ppu.Framebuffer[getFramebufferIndex(x, int(ly))] = getColorFromPalette(color)
	}

	// Transferimos los píxeles de la FIFO al framebuffer
	for x := 0; x < ScreenWidth; x++ {
		ppu.popPixelFromFIFO(x, int(ly))
	}
	//log.Printf("[PPU] Línea %d renderizada y enviada a framebuffer", ly)
	ppu.cycles -= baseCycles + spriteCycles
	ppu.setMode(ModeHBlank)
}

func getFramebufferIndex(x, y int) int {
	return y*ScreenWidth + x
}

func getColorFromPalette(color byte) uint32 {
	switch color {
	case 0:
		return 0xEEEEEEEE // blanco
	case 1:
		return 0xAAAAAAFF // gris claro
	case 2:
		return 0x555555FF // gris oscuro
	case 3:
		return 0x000000FF // negro
	default:
		panic("No se reconoce color")
	}
}
