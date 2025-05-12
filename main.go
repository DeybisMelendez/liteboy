package main

import (
	"fmt"
	"log"
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
	TargetCycle  = 70224 // Ciclos por fotograma (la cantidad de ciclos para un fotograma de GameBoy)
)

func main() {
	steps := 0
	// Inicialización de SDL
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
	defer texture.Destroy()

	// Cargar ROM
	cart := cartridge.NewCartridge("roms/tetris.gb")

	// Inicializar componentes
	gameBus := bus.NewBus(cart)
	gameCPU := cpu.NewCPU(gameBus)
	gamePPU := ppu.NewPPU(gameBus)

	// Variables de sincronización de FPS y ciclo
	//var lastTime time.Time
	//var frameDelay = time.Second / TargetFPS
	var cycleCount int

	// Variables de ciclo de CPU y PPU
	running := true
	for running {
		//fmt.Println("Step:", steps)
		steps++
		// Manejar eventos SDL
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Calculamos el tiempo transcurrido desde el último fotograma
		//start := time.Now()

		// Ejecutar la CPU y la PPU
		cycles := gameCPU.Step() // Ejecutamos la CPU
		cycleCount += cycles     // Acumulamos los ciclos de CPU
		gamePPU.Step(cycles * 4) // Ejecutamos la PPU
		// Actualizar textura con el framebuffer del PPU
		if cycleCount > 70224 {
			fmt.Println("imprime")
			cycleCount -= 70224

			err := texture.Update(nil, unsafe.Pointer(&gamePPU.Framebuffer[0]), ScreenWidth*4)
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
		}

		// Controlar FPS: mantener la tasa de fotogramas constante
		/*elapsed := time.Since(start)
		if elapsed < frameDelay {
			sdl.Delay(uint32(frameDelay - elapsed)) // Ajustar la velocidad del juego según el FPS
		}*/

		// Reportar la velocidad de la CPU cada segundo
		/*if time.Since(lastTime) >= time.Second {
			log.Printf("CPU Speed: %d cycles/sec", cycleCount)
			cycleCount = 0
			lastTime = time.Now()
		}*/
	}
}
