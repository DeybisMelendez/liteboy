package ppu

type STAT byte

const (
	ModeHBlank = 0
	ModeVBlank = 1
	ModeOAM    = 2
	ModeVRAM   = 3
)

func (ppu *PPU) isLYCInterruptEnabled() bool {
	return ppu.bus.Read(STATRegister)&0x40 != 0
}

func (ppu *PPU) isOAMInterruptEnabled() bool {
	return ppu.bus.Read(STATRegister)&0x20 != 0
}

func (ppu *PPU) isVBlankInterruptEnabled() bool {
	return ppu.bus.Read(STATRegister)&0x10 != 0
}

func (ppu *PPU) isHBlankInterruptEnabled() bool {
	return ppu.bus.Read(STATRegister)&0x08 != 0
}

func (ppu *PPU) setCoincidenceFlag(set bool) {
	if set {
		ppu.bus.Write(STATRegister, ppu.bus.Read(STATRegister)|0x04) // Set bit 2
	} else {
		ppu.bus.Write(STATRegister, ppu.bus.Read(STATRegister)&^0x04) // Clear bit 2 (bitwise AND NOT)
	}
}
func (ppu *PPU) updateCoincidenceFlag() {
	ly := ppu.bus.Read(LYRegister)
	lyc := ppu.bus.Read(LYCRegister)
	match := ly == lyc
	ppu.setCoincidenceFlag(match)
	if match && ppu.isLYCInterruptEnabled() {
		ppu.requestInterrupt(InterruptSTAT)
	}
}
