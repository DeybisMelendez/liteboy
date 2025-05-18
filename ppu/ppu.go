package ppu

import (
	"github.com/deybismelendez/liteboy/bus"
)

const (
	LCDCRegister     = 0xFF40
	STATRegister     = 0xFF41
	SCYRegister      = 0xFF42
	SCXRegister      = 0xFF43
	LYRegister       = 0xFF44
	LYCRegister      = 0xFF45
	WYRegister       = 0xFF4A // Window Y Position
	WXRegister       = 0xFF4B // Window X Position (el valor real en pantalla es WX - 7)
	ScreenWidth      = 160
	ScreenHeight     = 144
	TransparentPixel = 0xFFFFFFFF
)

type PPU struct {
	bus                  *bus.Bus
	Framebuffer          []uint32
	cycles               int
	spritesOnCurrentLine []*Sprite
	pixelFIFO            []uint32 // FIFO para los píxeles
	fifoSize             int
	windowLineCounter    uint16
}

func NewPPU(b *bus.Bus) *PPU {
	return &PPU{
		bus:         b,
		Framebuffer: make([]uint32, ScreenWidth*ScreenHeight),
		pixelFIFO:   make([]uint32, 0, ScreenWidth),
		fifoSize:    ScreenWidth}
}

// Función para agregar un píxel a la FIFO
func (ppu *PPU) addPixelToFIFO(pixel uint32) {
	if len(ppu.pixelFIFO) < ppu.fifoSize {
		ppu.pixelFIFO = append(ppu.pixelFIFO, pixel)
	}
}

// Función para extraer un píxel de la FIFO y colocarlo en el framebuffer
func (ppu *PPU) popPixelFromFIFO(x, y int) {
	if len(ppu.pixelFIFO) > 0 {
		pixel := ppu.pixelFIFO[0]
		ppu.pixelFIFO = ppu.pixelFIFO[1:] // Elimina el primer píxel de la cola
		ppu.Framebuffer[getFramebufferIndex(x, y)] = pixel
	}
}

func getFramebufferIndex(x, y int) int {
	return y*ScreenWidth + x
}

func getColorFromPalette(color byte) uint32 {
	switch color {
	case 0:
		return 0xEEEEEEFF // blanco
	case 1:
		return 0xAAAAAAFF // gris claro
	case 2:
		return 0x555555FF // gris oscuro
	case 3:
		return 0x000000FF // negro
	default:
		panic("No se reconoce color")
	}
}
