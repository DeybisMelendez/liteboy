package ppu

type Pixel struct {
	ColorIndex byte // 0-3
	Palette    byte // BG = 0, OBJ0 = 1, OBJ1 = 2
	Priority   byte // 0 = OBJ above BG, 1 = BG priority
	IsSprite   bool
}

func newPixel(colorIndex, palette, priority byte, isSprite bool) *Pixel {
	return &Pixel{
		ColorIndex: colorIndex,
		Palette:    palette,
		Priority:   priority,
		IsSprite:   isSprite,
	}
}
