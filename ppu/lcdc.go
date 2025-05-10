package ppu

type LCDC byte

// 0 = Off; 1 = On
func (l LCDC) IsLCDEnabled() bool {
	return l&0x80 != 0
}

// 0 = 9800–9BFF; 1 = 9C00–9FFF
func (l LCDC) WindowTileMapArea() uint16 {
	if l&0x40 != 0 {
		return 0x9C00
	} else {
		return 0x9800
	}
}

// 0 = Off; 1 = On
func (l LCDC) IsWindowEnabled() bool {
	return l&0x20 != 0
}

// 0 = 8800–97FF; 1 = 8000–8FFF
func (l LCDC) BGAndWindowTileArea() uint16 {
	if l&0x10 != 0 {
		return 0x8000 // unsigned addressing
	} else {
		return 0x8800 // signed addressing
	}
}

// 0 = 9800–9BFF; 1 = 9C00–9FFF
func (l LCDC) BGTileMapArea() uint16 {
	if l&0x08 != 0 {
		return 0x9C00
	} else {
		return 0x9800
	}
}

// 0 = 8×8; 1 = 8×16
func (l LCDC) ObjIs8x8() bool {
	return l&0x04 == 0
}

// 0 = Off; 1 = On
func (l LCDC) IsObjEnabled() bool {
	return l&0x02 != 0
}

// 0 = Off; 1 = On
func (l LCDC) IsBGWindowEnabled() bool {
	return l&0x01 != 0
}
