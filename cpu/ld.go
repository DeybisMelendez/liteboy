package cpu

// Load Data 8 bits puede ser implementado directamente en Step

// Load Data 16 bits
/*func (cpu *CPU) ld16(set func(uint16), value uint16) {
	set(value)
}*/

func (cpu *CPU) ldAF(value uint16) {
	cpu.a = byte(value >> 8)
	cpu.f = byte(value & 0xF0) // Solo se permiten los bits de flags (Z, N, H, C)
}

func (cpu *CPU) ldBC(value uint16) {
	cpu.b = byte(value >> 8)
	cpu.c = byte(value & 0xFF)
}
func (cpu *CPU) ldDE(value uint16) {
	cpu.d = byte(value >> 8)
	cpu.e = byte(value & 0xFF)
}

func (cpu *CPU) ldHL(value uint16) {
	cpu.h = byte(value >> 8)
	cpu.l = byte(value & 0xFF)
}

// LD HL, SP+e8
func (cpu *CPU) ld_HL_SP_e8() {
	e8 := cpu.getE8()
	sp := cpu.sp
	result := uint16(int32(int16(sp) + int16(e8)))

	cpu.ldHL(result)
	cpu.f = 0

	// Detectar Half-Carry y Carry basados en los 8 bits bajos
	if ((sp & 0x0F) + (uint16(e8) & 0x0F)) > 0x0F {
		cpu.f |= FlagH
	}
	if ((sp & 0xFF) + (uint16(e8) & 0xFF)) > 0xFF {
		cpu.f |= FlagC
	}
}
