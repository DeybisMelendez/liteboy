package ppu

import (
	"github.com/deybismelendez/liteboy/bus"
)

const (
	ScreenWidth  = 160
	ScreenHeight = 144

	ModeHBlank = 0
	ModeVBlank = 1
	ModeOAM    = 2
	ModeVRAM   = 3

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
	bus          *bus.Bus
	Framebuffer  [ScreenHeight][ScreenWidth]uint32
	mode         int
	cycleCounter int
	ly           byte
}

func NewPPU(bus *bus.Bus) *PPU {
	return &PPU{bus: bus}
}

func (p *PPU) Step(cycles int) {
	if !p.lcdEnabled() {
		return
	}
	//fmt.Println("Cycles", p.cycleCounter)
	p.cycleCounter += cycles * 4

	switch p.mode {
	case ModeOAM:
		if p.cycleCounter >= 80 {
			p.cycleCounter -= 80
			p.mode = ModeVRAM
		}
	case ModeVRAM:
		if p.cycleCounter >= 172 {
			p.cycleCounter -= 172
			p.renderScanline()
			p.mode = ModeHBlank
		}
	case ModeHBlank:
		if p.cycleCounter >= (456 - 80 - 172) {
			p.cycleCounter -= (456 - 80 - 172)
			p.ly++
			p.bus.IO[0x44] = p.ly
			if p.ly == 144 {
				p.mode = ModeVBlank
				// Generar interrupción VBlank aquí si es necesario
			} else {
				p.mode = ModeOAM
			}
		}
	case ModeVBlank:
		if p.cycleCounter >= 456 {
			p.cycleCounter -= 456
			p.ly++
			p.bus.IO[0x44] = p.ly
			if p.ly > 153 {
				p.ly = 0
				p.bus.IO[0x44] = p.ly
				p.mode = ModeOAM
			}
		}
	}
}

func (p *PPU) lcdEnabled() bool {
	return p.bus.IO[0x40]&LCDCFlagLCDEnable != 0
}

func (p *PPU) renderScanline() {
	lcdc := p.bus.IO[0x40]
	if lcdc&LCDCFlagBGDisplay == 0 {
		return
	}

	scy := p.bus.IO[0x42]
	scx := p.bus.IO[0x43]

	tileMapAddr := uint16(0x9800)
	if lcdc&LCDCFlagBGTileMap != 0 {
		tileMapAddr = 0x9C00
	}

	tileDataAddr := uint16(0x8800)
	signedIndex := true
	if lcdc&LCDCFlagBGTileData != 0 {
		tileDataAddr = 0x8000
		signedIndex = false
	}

	y := int(p.ly)
	for x := 0; x < ScreenWidth; x++ {
		pixelX := (uint16(x) + uint16(scx)) % 256
		pixelY := (uint16(y) + uint16(scy)) % 256
		tileX := pixelX / 8
		tileY := pixelY / 8
		tileIndexAddr := tileMapAddr + tileY*32 + tileX

		tileIndex := p.bus.Read(tileIndexAddr)
		var tileAddr uint16
		if signedIndex {
			tileAddr = tileDataAddr + uint16(int8(tileIndex))*16
		} else {
			tileAddr = tileDataAddr + uint16(tileIndex)*16
		}

		line := (pixelY % 8) * 2
		byte1 := p.bus.Read(tileAddr + uint16(line))
		byte2 := p.bus.Read(tileAddr + uint16(line) + 1)

		bit := 7 - (pixelX % 8)
		lo := (byte1 >> bit) & 1
		hi := (byte2 >> bit) & 1
		color := (hi << 1) | lo

		rgb := bgPalette(color)
		p.Framebuffer[y][x] = rgb
	}
}

func bgPalette(color byte) uint32 {
	switch color {
	case 0:
		return 0xFFFFFFFF // Blanco
	case 1:
		return 0xAAAAAAFF // Gris claro
	case 2:
		return 0x555555FF // Gris oscuro
	case 3:
		return 0x000000FF // Negro
	default:
		return 0xFF00FFFF // Magenta (color de error)
	}
}
