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
func (ppu *PPU) getMode() byte {
	return ppu.bus.Read(STATRegister) & 0x03
}

func (ppu *PPU) setMode(mode byte) {
	stat := ppu.bus.Read(STATRegister)
	stat = (stat &^ 0x03) | (mode & 0x03) // Bits 0-1 del STAT: modo actual
	ppu.bus.Write(STATRegister, stat)

	switch mode {
	case ModeHBlank:
		if ppu.isHBlankInterruptEnabled() {
			ppu.requestInterrupt(InterruptSTAT)
		}
	case ModeVBlank:
		ppu.requestInterrupt(InterruptVBlank)
		if ppu.isVBlankInterruptEnabled() {
			ppu.requestInterrupt(InterruptSTAT)
		}
	case ModeOAM:
		if ppu.isOAMInterruptEnabled() {
			ppu.requestInterrupt(InterruptSTAT)
		}
	}
}
