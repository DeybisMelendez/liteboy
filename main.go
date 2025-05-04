package main

import "github.com/deybismelendez/liteboy/cartridge"

func main() {
	cart := cartridge.NewCartridge("roms/tetris.gb")
	cart.PrintHeaderInfo()
}
