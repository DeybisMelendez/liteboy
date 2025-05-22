package main

import (
	"fmt"

	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 160
	ScreenHeight = 144
	Scale        = 4
)

type Liteboy struct {
	cpu         *cpu.CPU
	ppu         *ppu.PPU
	bus         *bus.Bus
	cycles      int
	targetTPS   int
	tpsMode     []int
	fastForward int
	image       *ebiten.Image
}

func NewLiteboy(cpu *cpu.CPU, ppu *ppu.PPU, bus *bus.Bus) *Liteboy {
	return &Liteboy{
		cpu:         cpu,
		ppu:         ppu,
		bus:         bus,
		image:       ebiten.NewImage(ScreenWidth, ScreenHeight),
		tpsMode:     []int{70224, 70224 * 2, 70224 * 3, 70224 * 4},
		fastForward: 1,
	}
}

func (liteboy *Liteboy) Update() error {
	for liteboy.cycles < liteboy.tpsMode[liteboy.targetTPS] {
		liteboy.cycles += liteboy.cpu.Step()
		// Pasamos ciclos reales transcurridos
		liteboy.handleGamepad()
	}

	liteboy.handleKeyboard()

	// Renderizado
	liteboy.image.WritePixels(liteboy.ppu.Framebuffer)
	liteboy.cycles -= liteboy.tpsMode[liteboy.targetTPS]

	return nil
}

func (liteboy *Liteboy) Draw(screen *ebiten.Image) {
	// Escalar la imagen a la ventana (multiplicando el tamaño)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(Scale, Scale)
	screen.DrawImage(liteboy.image, op)

	// Mostrar FPS en pantalla
	msg := fmt.Sprintf("LiteBoy Emulator - Press ESC to quit\nFPS: %.2f TPS: %.2f Target TPS: %d", ebiten.ActualFPS(), ebiten.ActualTPS()*float64(liteboy.tpsMode[liteboy.targetTPS])/60, liteboy.tpsMode[liteboy.targetTPS])
	ebitenutil.DebugPrint(screen, msg)
}

func (liteboy *Liteboy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth * Scale, ScreenHeight * Scale
}

func (liteboy *Liteboy) handleKeyboard() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		liteboy.targetTPS = liteboy.fastForward
	} else {
		liteboy.targetTPS = 0
	}
}

func (liteboy *Liteboy) handleGamepad() {
	// Leer el valor del registro P1 (0xFF00)
	p1 := liteboy.bus.Read(0xFF00)

	// Bit 4: dirección (0=activado), Bit 5: botones
	directionKeys := (p1 & (1 << 4)) == 0
	buttonKeys := (p1 & (1 << 5)) == 0

	var input byte = 0x0F // bits 0-3: todos presionados (1 = no presionado)

	if directionKeys {
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			input &= ^byte(1 << 0) // Bit 0 - Derecha
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			input &= ^byte(1 << 1) // Bit 1 - Izquierda
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			input &= ^byte(1 << 2) // Bit 2 - Arriba
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			input &= ^byte(1 << 3) // Bit 3 - Abajo
		}
	}

	if buttonKeys {
		if ebiten.IsKeyPressed(ebiten.KeyZ) {
			input &= ^byte(1 << 0) // Bit 0 - A
		}
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			input &= ^byte(1 << 1) // Bit 1 - B
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			input &= ^byte(1 << 2) // Bit 2 - Select
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			input &= ^byte(1 << 3) // Bit 3 - Start
		}
	}

	// Escribir bits 0-3 en el registro FF00 sin tocar bits 4-7
	liteboy.bus.Write(0xFF00, (p1&0xF0)|input)
}
