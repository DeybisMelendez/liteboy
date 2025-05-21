package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deybismelendez/liteboy/apu"
	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
	"github.com/deybismelendez/liteboy/timer"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <path_a_la_rom.gb>")
		return
	}

	romPath := os.Args[1]
	if len(os.Args) == 3 && os.Args[2] == "--info" {
		cart := cartridge.NewCartridge(romPath)
		cart.PrintHeaderInfo()
		os.Exit(0)
	}

	cart := cartridge.NewCartridge(romPath)
	gameBus := bus.NewBus(cart)
	gamePPU := ppu.NewPPU(gameBus)
	gameTimer := timer.NewTimer(gameBus)
	gameCPU := cpu.NewCPU(gameBus, gameTimer, gamePPU)
	gameAPU := apu.NewAPU(gameBus)
	game := NewLiteboy(gameCPU, gamePPU, gameBus, gameAPU)

	// Configurar ventana y correr el loop de Ebiten
	ebiten.SetWindowSize(ScreenWidth*Scale, ScreenHeight*Scale)
	ebiten.SetWindowTitle("LiteBoy Emulator")
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
