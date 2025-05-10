package main

import (
	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
)

const (
	ScreenWidth  = 160
	ScreenHeight = 144
	Scale        = 4
	TargetFPS    = 60
)

func main() {
	//start := time.Now()
	/*var cycleCount int
	// Inicializar SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("No se pudo inicializar SDL: %v", err)
	}
	defer sdl.Quit()

	// Crear ventana
	window, err := sdl.CreateWindow("LiteBoy Emulator",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		ScreenWidth*Scale, ScreenHeight*Scale,
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("No se pudo crear la ventana: %v", err)
	}
	defer window.Destroy()

	// Crear renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("No se pudo crear el renderer: %v", err)
	}
	defer renderer.Destroy()

	// Crear textura
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING, ScreenWidth, ScreenHeight)
	if err != nil {
		log.Fatalf("No se pudo crear la textura: %v", err)
	}
	defer texture.Destroy()*/

	// Cargar ROM
	cart := cartridge.NewCartridge("roms/tetris.gb")

	// Inicializar componentes
	gameBus := bus.NewBus(cart)
	gameCPU := cpu.NewCPU(gameBus)
	gamePPU := ppu.NewPPU(gameBus)
	steps := 200000
	//frameDelay := time.Second / TargetFPS
	for steps != 0 {
		//frameStart := time.Now()

		// Manejar eventos
		/*for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}*/

		// Ejecutar CPU y PPU
		cycles := gameCPU.Step()
		//cycleCount += cycles
		gamePPU.Step(cycles)
		steps--
		// Actualizar textura con el framebuffer del PPU
		/*pixels := make([]uint32, ScreenWidth*ScreenHeight)
		for y := 0; y < ScreenHeight; y++ {
			for x := 0; x < ScreenWidth; x++ {
				pixels[y*ScreenWidth+x] = gamePPU.Framebuffer[y][x]
			}
		}
		err = texture.Update(nil, unsafe.Pointer(&pixels[0]), ScreenWidth*4)
		if err != nil {
			log.Printf("Error al actualizar la textura: %v", err)
		}

		// Renderizar
		renderer.Clear()
		err = renderer.Copy(texture, nil, nil)
		if err != nil {
			log.Printf("Error al copiar la textura: %v", err)
		}
		renderer.Present()
		*/
		// Controlar FPS
		/*frameTime := time.Since(frameStart)
		if frameDelay > frameTime {
			sdl.Delay(uint32((frameDelay - frameTime).Milliseconds()))
		}*/
		// Calcula velocidad cada segundo
		/*if time.Since(start) >= time.Second {
			fmt.Printf("CPU Speed: %d cycles/sec\n", cycleCount)
			cycleCount = 0
			start = time.Now()
		}*/
	}
}
