package cpu

// Load Data 8 bits puede ser implementado directamente en Step

// Load Data 16 bits
func (cpu *CPU) ld16(set func(uint16), value uint16) {
	set(value)
}

// LD HL, SP+e8
func (cpu *CPU) ld_HL_SP_e8() {
	e8 := cpu.getE8()
	sp := cpu.sp
	result := uint16(int32(int16(sp) + int16(e8)))

	cpu.setHL(result)
	cpu.f = 0

	// Detectar Half-Carry y Carry basados en los 8 bits bajos
	if ((sp & 0x0F) + (uint16(e8) & 0x0F)) > 0x0F {
		cpu.f |= FlagH
	}
	if ((sp & 0xFF) + (uint16(e8) & 0xFF)) > 0xFF {
		cpu.f |= FlagC
	}
}
