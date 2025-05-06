package main

import (
	"time"

	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
)

func main() {
	cart := cartridge.NewCartridge("roms/tetris.gb")
	cart.PrintHeaderInfo()
	cpu := cpu.NewCPU()
	cpu.LoadMemory(cart.GetROM())
	for {
		cpu.Step()
		time.Sleep(10 * time.Millisecond)
	}
}
