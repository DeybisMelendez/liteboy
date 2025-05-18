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
)

var WhiteColor Pixel = Pixel{R: 0xEE, G: 0xEE, B: 0xEE, A: 0xFF}

type PPU struct {
	bus                  *bus.Bus
	Framebuffer          []byte
	cycles               int
	spritesOnCurrentLine []*Sprite
	pixelFIFO            []*Pixel // FIFO para los píxeles
	fifoSize             int
	windowLineCounter    uint16
}

func NewPPU(b *bus.Bus) *PPU {
	return &PPU{
		bus:         b,
		Framebuffer: make([]byte, ScreenWidth*ScreenHeight*4),
		pixelFIFO:   make([]*Pixel, 0, ScreenWidth),
		fifoSize:    ScreenWidth}
}

// Función para agregar un píxel a la FIFO
func (ppu *PPU) addPixelToFIFO(pixel *Pixel) {
	if len(ppu.pixelFIFO) < ppu.fifoSize {
		ppu.pixelFIFO = append(ppu.pixelFIFO, pixel)
	}
}

// Función para extraer un píxel de la FIFO y colocarlo en el framebuffer
func (ppu *PPU) popPixelFromFIFO(x, y int) {
	pixel := ppu.pixelFIFO[0]
	ppu.pixelFIFO = ppu.pixelFIFO[1:]
	for i := 1; i < 4; i++ {
		ppu.Framebuffer[getFramebufferIndex(x, y)] = pixel.R
		ppu.Framebuffer[getFramebufferIndex(x, y)+1] = pixel.G
		ppu.Framebuffer[getFramebufferIndex(x, y)+2] = pixel.B
		ppu.Framebuffer[getFramebufferIndex(x, y)+3] = pixel.A

	}
}

func getFramebufferIndex(x, y int) int {
	return (y*ScreenWidth + x) * 4
}

func getColorFromPalette(color byte) *Pixel {
	switch color {
	case 0:
		return newPixel(0xEE, 0xEE, 0xEE, 0xFF) // blanco
	case 1:
		return newPixel(0xAA, 0xAA, 0xAA, 0xFF) // gris claro
	case 2:
		return newPixel(0x55, 0x55, 0x55, 0xFF) // gris oscuro
	case 3:
		return newPixel(0x00, 0x00, 0x00, 0xFF) // negro
	default:
		panic("No se reconoce color")
	}
}
