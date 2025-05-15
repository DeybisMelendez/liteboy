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
	WYRegister   = 0xFF4A // Window Y Position
	WXRegister   = 0xFF4B // Window X Position (el valor real en pantalla es WX - 7)

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
	}
}

// Función para extraer un píxel de la FIFO y colocarlo en el framebuffer
func (ppu *PPU) popPixelFromFIFO(x, y int) {
	if len(ppu.pixelFIFO) > 0 {
		pixel := ppu.pixelFIFO[0]
		ppu.pixelFIFO = ppu.pixelFIFO[1:] // Elimina el primer píxel de la cola
		ppu.Framebuffer[getFramebufferIndex(x, y)] = pixel
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
		ppu.scanOAM()
	case ModeVRAM:
		ppu.runVRAM()
	case ModeHBlank:
		ppu.runHBlank()
	case ModeVBlank:
		ppu.runVBlank()
	}
}

func (ppu *PPU) scanOAM() {
	if ppu.IsOAMInterruptEnabled() {
		ppu.requestInterrupt(InterruptSTAT)
	}
	if ppu.cycles >= 80 {
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
		ppu.spritesOnCurrentLine = result
		ppu.cycles -= 80
		ppu.setMode(ModeVRAM)
	}
}

func (ppu *PPU) runHBlank() {
	if ppu.isHBlankInterruptEnabled() {
		ppu.requestInterrupt(InterruptSTAT)
	}
	if ppu.cycles >= 204-len(ppu.spritesOnCurrentLine)*2 {
		ppu.bus.Write(LYRegister, ppu.bus.Read(LYRegister)+1)
		if ppu.bus.Read(LYRegister) == 144 {
			ppu.setMode(ModeVBlank)
			ppu.requestInterrupt(InterruptVBlank)
			ppu.requestInterrupt(InterruptSTAT)
		} else {
			ppu.setMode(ModeOAM)
		}
		ppu.updateCoincidenceFlag()
		ppu.cycles -= 204 - len(ppu.spritesOnCurrentLine)*2
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
	// Ciclos mínimos para entrar a Modo 3
	const baseCycles = 172
	spriteCycles := len(ppu.spritesOnCurrentLine) * 2
	totalCycles := baseCycles + spriteCycles

	// No hagas nada hasta completar Modo 3
	if ppu.cycles < totalCycles {
		return
	}

	ly := int(ppu.bus.Read(LYRegister))   // Línea actual 0–143
	scx := int(ppu.bus.Read(SCXRegister)) // Scroll X 0–255
	scy := int(ppu.bus.Read(SCYRegister)) // Scroll Y 0–255
	lcdc := ppu.bus.Read(LCDCRegister)
	wx := int(ppu.bus.Read(WXRegister)) // Window X+7
	wy := int(ppu.bus.Read(WYRegister)) // Window Y
	bgTileMap := uint16(0x9800)
	if lcdc&LCDCFlagBGTileMap != 0 {
		bgTileMap = 0x9C00
	}
	useSigned := (lcdc & LCDCFlagBGTileData) == 0 // modo 8800 o 8000

	// Para cada píxel de la línea...
	for px := 0; px < ScreenWidth; px++ {
		var tileX, tileY, tileMapAddr uint16
		var inWindow bool

		// ¿Se dibuja window aquí?
		if lcdc&LCDCFlagWindowEnable != 0 &&
			ly >= wy && px >= wx-7 {
			inWindow = true
			// Coordenadas dentro de la ventana
			tileX = uint16(px-(wx-7)) / 8
			// Contador Y interno de ventana
			// (se reinicia en VBlank y sólo crece cuando entra)
			windowLine := ly - wy
			tileY = uint16(windowLine) / 8
			tileMapAddr = 0x9800
			if lcdc&LCDCFlagWindowTileMap != 0 {
				tileMapAddr = 0x9C00
			}
		} else {
			// Background scrolling normal
			scrollX := (px + scx) & 0xFF // wrap-around 256 px
			scrollY := (ly + scy) & 0xFF
			tileX = uint16(scrollX) / 8
			tileY = uint16(scrollY) / 8
			tileMapAddr = bgTileMap
		}

		// Offset dentro del tile map (32 cols)
		mapOffset := tileY*32 + tileX
		tileIndex := ppu.bus.Read(tileMapAddr + mapOffset)

		// Dirección base de datos de tile
		var tileAddr uint16
		if useSigned {
			tileAddr = 0x9000 + uint16(int8(tileIndex))*16
		} else {
			tileAddr = 0x8000 + uint16(tileIndex)*16
		}

		// Cálculo de fila dentro del tile
		var rowY int
		if inWindow {
			rowY = (ly - wy) % 8
		} else {
			rowY = ((ly + scy) & 0xFF) % 8
		}
		addr := tileAddr + uint16(rowY*2)
		byte1 := ppu.bus.Read(addr)
		byte2 := ppu.bus.Read(addr + 1)

		// bit de píxel dentro del byte
		var bitX int
		if inWindow {
			bitX = 7 - ((px - (wx - 7)) % 8)
		} else {
			bitX = 7 - ((px + scx) & 7)
		}

		// Construcción del color final
		colorID := (((byte2 >> bitX) & 1) << 1) | ((byte1 >> bitX) & 1)
		palette := ppu.bus.Read(0xFF47)
		color := (palette >> (colorID * 2)) & 0x03
		ppu.addPixelToFIFO(getColorFromPalette(color))
	}

	// Vuelca FIFO al framebuffer y renderiza sprites
	for x := 0; x < ScreenWidth; x++ {
		ppu.popPixelFromFIFO(x, int(ppu.bus.Read(LYRegister)))
	}
	ppu.renderSprites()

	// Ajusta ciclos y pasa a HBlank
	ppu.cycles -= totalCycles
	ppu.setMode(ModeHBlank)
}

func (ppu *PPU) renderSprites() {
	spriteHeight := byte(8)
	if !ppu.isObj8x8() {
		spriteHeight = 16
	}

	ly := ppu.bus.Read(LYRegister)

	for _, sprite := range ppu.spritesOnCurrentLine {
		spriteY := int(sprite.Y) - 16
		spriteX := int(sprite.X) - 8
		line := int(ly) - spriteY

		if sprite.Atributes&0x40 != 0 { // Y flip
			line = int(spriteHeight) - 1 - line
		}

		tileIndex := sprite.TileIndex
		if spriteHeight == 16 {
			tileIndex &= 0xFE // Ignorar bit 0 en modo 8x16
		}

		tileAddr := 0x8000 + uint16(tileIndex)*16 + uint16(line)*2
		byte1 := ppu.bus.Read(tileAddr)
		byte2 := ppu.bus.Read(tileAddr + 1)

		for x := 0; x < 8; x++ {
			bit := 7 - x
			if sprite.Atributes&0x20 != 0 { // X flip
				bit = x
			}

			colorID := (((byte2 >> bit) & 1) << 1) | ((byte1 >> bit) & 1)
			if colorID == 0 {
				continue // Transparente
			}

			var paletteAddr uint16 = 0xFF48
			if sprite.Atributes&0x10 != 0 {
				paletteAddr = 0xFF49
			}

			palette := ppu.bus.Read(paletteAddr)
			color := (palette >> (colorID * 2)) & 0x03
			screenX := spriteX + x

			if screenX < 0 || screenX >= ScreenWidth {
				continue
			}

			// Prioridad: fondo (bit 7)
			bgPriority := sprite.Atributes&0x80 != 0
			if bgPriority {
				// Omitimos dibujar si fondo no es color 0
				bgPixel := ppu.Framebuffer[getFramebufferIndex(screenX, int(ly))]
				if bgPixel != getColorFromPalette(0) {
					continue
				}
			}
			//ppu.addPixelToFIFO(getColorFromPalette(color))
			ppu.Framebuffer[getFramebufferIndex(screenX, int(ly))] = getColorFromPalette(color)
		}
	}

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
