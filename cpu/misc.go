package cpu

func (cpu *CPU) halt() {
	cpu.halted = true
}

func (cpu *CPU) ccf() {
	cpu.f ^= FlagC          // Toggle Carry
	cpu.f &^= FlagN | FlagH // Clear N and H
}

func (cpu *CPU) di() {
	cpu.ime = false
}

func (cpu *CPU) ei() {
	cpu.enableIME = true
}
func (cpu *CPU) reti() {
	cpu.ime = true
	cpu.ret()
}
