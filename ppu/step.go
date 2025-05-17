package ppu

func (ppu *PPU) Step(tCycles int) {
	ppu.bus.Client = 1
	if !ppu.isLCDEnabled() {
		ppu.bus.Write(LYRegister, 0)
		ppu.setMode(ModeHBlank)
		return
	}

	ppu.cycles += tCycles

	switch ppu.getMode() {
	case ModeOAM:
		ppu.scanOAM()
	case ModeVRAM:
		ppu.runVRAM()
	case ModeHBlank:
		ppu.runHBlank()
	case ModeVBlank:
		ppu.runVBlank()
	}
}
