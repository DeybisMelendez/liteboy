package main

import (
	"fmt"
	"log"
	"os"
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
	TargetCycle  = 70224
)

func initSDL() (*sdl.Window, *sdl.Renderer, *sdl.Texture) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Error al iniciar SDL: %v", err)
	}

	window, err := sdl.CreateWindow("LiteBoy Emulator",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		ScreenWidth*Scale, ScreenHeight*Scale,
		sdl.WINDOW_SHOWN)
	if err != nil {
		sdl.Quit()
		log.Fatalf("Error al crear ventana: %v", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		sdl.Quit()
		log.Fatalf("Error al crear renderer: %v", err)
	}

	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING,
		ScreenWidth, ScreenHeight,
	)
	if err != nil {
		renderer.Destroy()
		window.Destroy()
		sdl.Quit()
		log.Fatalf("Error al crear textura: %v", err)
	}

	return window, renderer, texture
}

func loadROM(path string) *cartridge.Cartridge {
	cart := cartridge.NewCartridge(path)
	if cart == nil {
		log.Fatalf("Error al cargar ROM: %s", path)
	}
	return cart
}

func mainLoop(cpu *cpu.CPU, ppu *ppu.PPU, renderer *sdl.Renderer, texture *sdl.Texture) {
	cycleCount := 0
	running := true

	for running {
		cycleCount += cpu.Step()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if _, ok := event.(*sdl.QuitEvent); ok {
				running = false
			}
		}

		if cycleCount >= TargetCycle {
			cycleCount -= TargetCycle

			err := texture.Update(nil, unsafe.Pointer(&ppu.Framebuffer[0]), ScreenWidth*4)
			if err != nil {
				log.Printf("Error al actualizar textura: %v", err)
			}

			renderer.Clear()
			if err := renderer.Copy(texture, nil, nil); err != nil {
				log.Printf("Error al copiar textura: %v", err)
			}
			renderer.Present()
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <path_a_la_rom.gb>")
		return
	}

	romPath := os.Args[1]
	cart := loadROM(romPath)
	gameBus := bus.NewBus(cart)
	gamePPU := ppu.NewPPU(gameBus)
	gameCPU := cpu.NewCPU(gameBus, gamePPU)

	window, renderer, texture := initSDL()
	defer func() {
		texture.Destroy()
		renderer.Destroy()
		window.Destroy()
		sdl.Quit()
	}()

	mainLoop(gameCPU, gamePPU, renderer, texture)
}
