package ppu

import (
	"github.com/deybismelendez/liteboy/bus.go"
)

type PPU struct {
	lcdc, stat       byte
	scy, scx         byte
	ly, lyc          byte
	bgp, obp0, obp1  byte
	windowY, windowX byte
	cycleCount       int
	mode             byte

	framebuffer [160 * 144]uint32
	bus         *bus.Bus
}

const (
	ModeHBlank = 0
	ModeVBlank = 1
	ModeOAM    = 2
	ModeVRAM   = 3
)

func NewPPU(bus *bus.Bus) *PPU {
	ppu := &PPU{bus: bus}
	return ppu
}

func (ppu *PPU) Step(cycles int) {
	ppu.cycleCount += cycles

	switch ppu.mode {
	case ModeOAM:
		if ppu.cycleCount >= 80 {
			ppu.cycleCount -= 80
			ppu.mode = ModeVRAM
		}
	case ModeVRAM:
		if ppu.cycleCount >= 172 {
			ppu.cycleCount -= 172
			ppu.renderScanline()
			ppu.mode = ModeHBlank
		}
	case ModeHBlank:
		if ppu.cycleCount >= 204 {
			ppu.cycleCount -= 204
			ppu.ly++
			if ppu.ly == 144 {
				ppu.mode = ModeVBlank
				ppu.requestInterrupt(0x01) // VBlank interrupt
			} else {
				ppu.mode = ModeOAM
			}
		}
	case ModeVBlank:
		if ppu.cycleCount >= 456 {
			ppu.cycleCount -= 456
			ppu.ly++
			if ppu.ly > 153 {
				ppu.ly = 0
				ppu.mode = ModeOAM
			}
		}
	}
}

func (ppu *PPU) renderScanline() {
	y := int(ppu.ly)
	scx, scy := int(ppu.scx), int(ppu.scy)
	tileMapBase := 0x9800
	if ppu.lcdc&0x08 != 0 {
		tileMapBase = 0x9C00
	}
	tileDataBase := 0x8800
	signedIndex := true
	if ppu.lcdc&0x10 != 0 {
		tileDataBase = 0x8000
		signedIndex = false
	}
	for x := 0; x < 160; x++ {
		px := (x + scx) & 0xFF
		py := (y + scy) & 0xFF
		tileX := px / 8
		tileY := py / 8
		tileMapAddr := uint16(tileMapBase + tileY*32 + tileX)
		tileIndex := ppu.bus.Read(tileMapAddr)
		var tileAddr uint16
		if signedIndex {
			tileAddr = uint16(int16(int8(tileIndex)))*16 + uint16(tileDataBase)
		} else {
			tileAddr = uint16(tileIndex)*16 + uint16(tileDataBase)
		}
		line := py % 8
		byte1 := ppu.bus.Read(tileAddr + uint16(line*2))
		byte2 := ppu.bus.Read(tileAddr + uint16(line*2) + 1)
		bit := 7 - (px % 8)
		hi := (byte2 >> bit) & 1
		lo := (byte1 >> bit) & 1
		color := (hi << 1) | lo
		ppu.framebuffer[y*160+x] = ppu.mapColor(color, ppu.bgp)
	}
	ppu.renderSprites(y)
}

func (ppu *PPU) renderSprites(y int) {
	spriteHeight := 8
	if ppu.lcdc&0x04 != 0 {
		spriteHeight = 16
	}
	for i := 0; i < 40; i++ {
		index := i * 4
		yPos := int(ppu.bus.Read(0xFE00+uint16(index))) - 16
		xPos := int(ppu.bus.Read(0xFE00+uint16(index+1))) - 8
		tile := ppu.bus.Read(0xFE00 + uint16(index+2))
		attr := ppu.bus.Read(0xFE00 + uint16(index+3))
		if y < yPos || y >= yPos+spriteHeight {
			continue
		}
		line := y - yPos
		if attr&0x40 != 0 {
			line = spriteHeight - 1 - line
		}
		addr := uint16(tile)*16 + uint16(line*2)
		byte1 := ppu.bus.Read(0x8000 + addr)
		byte2 := ppu.bus.Read(0x8000 + addr + 1)
		for x := 0; x < 8; x++ {
			bit := 7 - x
			if attr&0x20 != 0 {
				bit = x
			}
			hi := (byte2 >> bit) & 1
			lo := (byte1 >> bit) & 1
			color := (hi << 1) | lo
			if color == 0 {
				continue
			}
			pal := ppu.obp0
			if attr&0x10 != 0 {
				pal = ppu.obp1
			}
			finalX := xPos + x
			if finalX < 0 || finalX >= 160 {
				continue
			}
			ppu.framebuffer[y*160+finalX] = ppu.mapColor(color, pal)
		}
	}
}

func (ppu *PPU) mapColor(colorID byte, palette byte) uint32 {
	switch (palette >> (colorID * 2)) & 0x03 {
	case 0:
		return 0xFFFFFFFF // White
	case 1:
		return 0xAAAAAAFF // Light Gray
	case 2:
		return 0x555555FF // Dark Gray
	case 3:
		return 0x000000FF // Black
	default:
		return 0xFFFFFFFF
	}
}

func (ppu *PPU) requestInterrupt(flag byte) {
	val := ppu.bus.Read(0xFF0F)
	val |= flag
	ppu.bus.Write(0xFF0F, val)
}

func (ppu *PPU) Read(addr uint16) byte {
	switch addr {
	case 0xFF40:
		return ppu.lcdc
	case 0xFF41:
		return ppu.stat
	case 0xFF42:
		return ppu.scy
	case 0xFF43:
		return ppu.scx
	case 0xFF44:
		return ppu.ly
	case 0xFF45:
		return ppu.lyc
	case 0xFF47:
		return ppu.bgp
	case 0xFF48:
		return ppu.obp0
	case 0xFF49:
		return ppu.obp1
	case 0xFF4A:
		return ppu.windowY
	case 0xFF4B:
		return ppu.windowX
	default:
		return 0xFF
	}
}

func (ppu *PPU) Write(addr uint16, val byte) {
	switch addr {
	case 0xFF40:
		ppu.lcdc = val
	case 0xFF41:
		ppu.stat = val
	case 0xFF42:
		ppu.scy = val
	case 0xFF43:
		ppu.scx = val
	case 0xFF44:
		ppu.ly = 0 // LY is read-only
	case 0xFF45:
		ppu.lyc = val
	case 0xFF47:
		ppu.bgp = val
	case 0xFF48:
		ppu.obp0 = val
	case 0xFF49:
		ppu.obp1 = val
	case 0xFF4A:
		ppu.windowY = val
	case 0xFF4B:
		ppu.windowX = val
	}
}

func (ppu *PPU) GetFrameBuffer() *[160 * 144]uint32 {
	return &ppu.framebuffer
}

func (ppu *PPU) IsVBlankReady() bool {
	return ppu.ly == 144
}
