package cpu

// update actualiza el estado del cpu saltando al siguiente ciclo
func (cpu *CPU) update(bytes uint16, cycles int) int {
	cpu.pc += bytes
	cpu.cycles += cycles
	if cpu.enableIME {
		cpu.ime = true
		cpu.enableIME = false
	}

	return cycles
}

func (cpu *CPU) halt() {
	cpu.halted = true
}

// CCF (Complement Carry Flag)
func (cpu *CPU) ccf() {
	if cpu.f&FlagC != 0 {
		cpu.f &^= FlagC
	} else {
		cpu.f |= FlagC
	}
	cpu.f &^= FlagN | FlagH
}
func (cpu *CPU) di() {
	cpu.ime = false
}

func (cpu *CPU) ei() {
	cpu.enableIME = true
}
func (cpu *CPU) reti() {
	cpu.ret()
	cpu.ime = true
}
