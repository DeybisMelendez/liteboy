package ppu

type LCDC byte

// false = Off; true = On
func (l LCDC) IsLCDEnabled() bool {
	return l&0x80 != 0
}

// 9800; 9C00
func (l LCDC) GetWindowTileMapArea() uint16 {
	if l&0x40 != 0 {
		return 0x9C00
	} else {
		return 0x9800
	}
}

// false = Off; true = On
func (l LCDC) IsWindowEnabled() bool {
	return l&0x20 != 0
}

// 8800; 8000
func (l LCDC) GetBGAndWindowTileDataArea() uint16 {
	if l&0x10 != 0 {
		return 0x8000 // unsigned addressing
	} else {
		return 0x8800 // signed addressing
	}
}

// 9800; 9C00
func (l LCDC) GetBGTileMapArea() uint16 {
	if l&0x08 != 0 {
		return 0x9C00
	} else {
		return 0x9800
	}
}

func (ppu *PPU) isLCDEnabled() bool {
	return ppu.bus.Read(LCDCRegister)&0x80 != 0
}
func (ppu *PPU) isObj8x8() bool {
	return ppu.bus.Read(LCDCRegister)&0x04 == 0
}

// false = Off; true = On
func (l LCDC) IsObjEnabled() bool {
	return l&0x02 != 0
}

// false = Off; true = On
func (l LCDC) IsBGAndWindowEnabled() bool {
	return l&0x01 != 0
}
