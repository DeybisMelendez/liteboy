package main

import (
	"log"
	"time"
	"unsafe"

	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	ScreenWidth  = 160
	ScreenHeight = 144
	Scale        = 4
	TargetFPS    = 60
)

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Error al iniciar SDL: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("LiteBoy", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		ScreenWidth*Scale, ScreenHeight*Scale, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Error al crear ventana: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Error al crear renderer: %v", err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING,
		ScreenWidth, ScreenHeight)
	if err != nil {
		log.Fatalf("Error al crear textura: %v", err)
	}
	defer texture.Destroy()

	// Cargar ROM
	cart := cartridge.NewCartridge("roms/yakuman.gb")

	// Inicializar componentes
	gameBus := bus.NewBus(cart.GetROM())
	gameCPU := cpu.NewCPU(gameBus)
	gamePPU := ppu.NewPPU(gameBus)

	// Bucle principal
	running := true
	for running {
		start := time.Now()

		cycles := gameCPU.Step()
		gamePPU.Step(cycles)

		// Dibujar frame

		frame := *gamePPU.GetFrameBuffer()
		texture.Update(nil, unsafe.Pointer(&frame[0]), ScreenWidth*4)
		renderer.Clear()
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		// Control de FPS (~60 Hz)
		elapsed := time.Since(start)
		frameDuration := time.Second / TargetFPS
		if elapsed < frameDuration {
			time.Sleep(frameDuration - elapsed)
		}

		// Manejo de eventos
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
	}
}
