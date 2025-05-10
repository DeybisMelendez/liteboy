package cpu

func (cpu *CPU) updateZ0H_(result, value byte) {
	cpu.f &^= FlagN  // N = 0
	if result == 0 { // Z
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
	if (value&0x0F)+1 > 0x0F { // H
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
}
func (cpu *CPU) updateZ1H_(result, value byte) {
	cpu.f |= FlagN   // N = 1
	if result == 0 { // Z
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
	if (value & 0x0F) == 0x00 { // H
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
}
