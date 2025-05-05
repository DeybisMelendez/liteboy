package main

import (
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
)

func main() {
	cart := cartridge.NewCartridge("roms/zelda.gb")
	cart.PrintHeaderInfo()
	cpu := cpu.NewCPU()
	cpu.LoadMemory(cart.GetROM())
	for {
		cpu.Step()
	}
}
