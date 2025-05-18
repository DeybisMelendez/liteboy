package ppu

func (ppu *PPU) scanOAM() {
	if ppu.cycles < 80 {
		return
	}
	ppu.cycles -= 80

	spriteHeight := ppu.getObjHeight()

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
			sprite := newSprite(x, y, tile, attr)
			result = append(result, sprite)

			if len(result) == MaxSpritesPerLine {
				break
			}
		}
	}
	ppu.spritesOnCurrentLine = result
	ppu.setMode(ModeVRAM)
}
