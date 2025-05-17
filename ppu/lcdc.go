package ppu

type LCDC byte

const (
	LCDCFlagBGEnablePriority = 1 << 0
	LCDCFlagOBJDisplay       = 1 << 1
	LCDCFlagOBJSize          = 1 << 2
	LCDCFlagBGTileMap        = 1 << 3
	LCDCFlagBGWindowTileData = 1 << 4
	LCDCFlagWindowEnable     = 1 << 5
	LCDCFlagWindowTileMap    = 1 << 6
	LCDCFlagEnable           = 1 << 7
)

func (ppu *PPU) isLCDEnabled() bool {
	return ppu.bus.Read(LCDCRegister)&LCDCFlagEnable != 0
}

// 9800; 9C00
func (ppu *PPU) getWindowTileMapArea() uint16 {
	if ppu.bus.Read(LCDCRegister)&LCDCFlagWindowTileMap != 0 {
		return 0x9C00
	} else {
		return 0x9800
	}
}
func (ppu *PPU) isWindowEnabled() bool {
	return ppu.bus.Read(LCDCRegister)&LCDCFlagWindowEnable != 0
}

// 8800; 8000
func (ppu *PPU) getBGAndWindowTileDataArea() uint16 {
	if ppu.bus.Read(LCDCRegister)&LCDCFlagBGWindowTileData != 0 {
		return 0x8000 // unsigned addressing
	} else {
		return 0x8800 // signed addressing
	}
}

// 9800; 9C00
func (ppu *PPU) getBGTileMapArea() uint16 {
	if ppu.bus.Read(LCDCRegister)&LCDCFlagBGTileMap != 0 {
		return 0x9C00
	} else {
		return 0x9800
	}
}

func (ppu *PPU) getObjHeight() byte {
	if ppu.bus.Read(LCDCRegister)&LCDCFlagOBJSize == 0 {
		return 8
	}
	return 16
}

// false = Off; true = On
func (ppu *PPU) isObjEnabled() bool {
	return ppu.bus.Read(LCDCRegister)&LCDCFlagOBJDisplay != 0
}

// false = Off; true = On
func (ppu *PPU) isBGAndWindowEnabledPriority() bool {
	return ppu.bus.Read(LCDCRegister)&LCDCFlagBGEnablePriority != 0
}
