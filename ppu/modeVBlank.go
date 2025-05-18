package ppu

func (ppu *PPU) runVBlank() {
	if ppu.isVBlankInterruptEnabled() {
		ppu.requestInterrupt(InterruptSTAT)
	}
	if ppu.cycles < 456 {
		return
	}
	ppu.cycles -= 456
	ly := ppu.bus.Read(LYRegister)
	if ly == 154 {
		ppu.bus.Write(LYRegister, 0)
		ppu.updateCoincidenceFlag()
		ppu.setMode(ModeOAM)
	} else {
		ppu.bus.Write(LYRegister, ly+1)
		ppu.updateCoincidenceFlag()
	}
}
