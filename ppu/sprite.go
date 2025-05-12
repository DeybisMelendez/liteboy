package ppu

const MaxSpritesPerLine = 10

type Sprite struct {
	X, Y, TileIndex, Atributes byte
	OAMIndex                   uint16
}

func newSprite(x, y, tileIndex, atributes byte, OAMIndex uint16) *Sprite {
	return &Sprite{
		X:         x,
		Y:         y,
		TileIndex: tileIndex,
		Atributes: atributes,
		OAMIndex:  OAMIndex,
	}
}
