package ppu

func (ppu *PPU) runHBlank() {
	if ppu.isHBlankInterruptEnabled() {
		ppu.requestInterrupt(InterruptSTAT)
	}
	if ppu.cycles < 204 {
		return
	}
	ppu.cycles -= 204
	ly := ppu.bus.Read(LYRegister)
	if ly == 144 {
		ppu.setMode(ModeVBlank)
		ppu.requestInterrupt(InterruptVBlank)
		ppu.requestInterrupt(InterruptSTAT)
	} else {
		ppu.setMode(ModeOAM)
	}
	ppu.bus.Write(LYRegister, ly+1)
	ppu.updateCoincidenceFlag()

}
