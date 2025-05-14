package main

import (
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
	TargetCycle  = 70224
)

/*func main() {
	cart := cartridge.NewCartridge("roms/tetris.gb") // Usá una ROM de test si tenés
	b := bus.NewBus(cart)
	c := cpu.NewCPU(b)

	reader := bufio.NewReader(os.Stdin)
	step := 0

	for {
		c.Step()

		fmt.Print("\nPresioná ENTER para continuar o 'q' + ENTER para salir: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			fmt.Println("Saliendo del paso a paso.")
			break
		}

		step++
	}
}*/

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("No se pudo inicializar SDL: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("LiteBoy Emulator",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		ScreenWidth*Scale, ScreenHeight*Scale,
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("No se pudo crear la ventana: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("No se pudo crear el renderer: %v", err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING, ScreenWidth, ScreenHeight)
	if err != nil {
		log.Fatalf("No se pudo crear la textura: %v", err)
	}
	defer texture.Destroy()
	cart := cartridge.NewCartridge("roms/blargg/interrupt_time/interrupt_time.gb")

	gameBus := bus.NewBus(cart)
	gameCPU := cpu.NewCPU(gameBus)
	gamePPU := ppu.NewPPU(gameBus)

	cycleCount := 0
	//lastFPSUpdate := time.Now()
	//frames := 0
	//frameDelay := time.Second / TargetFPS

	running := true
	for running {
		//start := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		if cycleCount >= TargetCycle {
			cycleCount -= TargetCycle
			//frames++
			//fmt.Println("Frame:", frames)

			err := texture.Update(nil, unsafe.Pointer(&gamePPU.Framebuffer[0]), ScreenWidth*4)
			if err != nil {
				log.Printf("Error al actualizar la textura: %v", err)
			}

			renderer.Clear()
			err = renderer.Copy(texture, nil, nil)
			if err != nil {
				log.Printf("Error al copiar la textura: %v", err)
			}
			renderer.Present()
		} else {
			cycles := gameCPU.Step()
			cycleCount += cycles
			gamePPU.Step(cycles)
		}

		/*if time.Since(lastFPSUpdate) >= time.Second {
			log.Printf("[FPS] %d frames/s", frames)
			frames = 0
			lastFPSUpdate = time.Now()
		}

		elapsed := time.Since(start)
		if elapsed < frameDelay {
			sdl.Delay(uint32((frameDelay - elapsed).Milliseconds()))
		}*/
	}
}
