package ppu

type Pixel struct {
	R byte
	G byte
	B byte
	A byte
}

func newPixel(r, g, b, a byte) *Pixel {
	return &Pixel{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}
