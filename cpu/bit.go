package cpu

// RLCA
func (cpu *CPU) rlca() {
	carry := cpu.a >> 7
	result := (cpu.a << 1) | carry
	cpu.f = 0
	if carry != 0 {
		cpu.f |= FlagC
	}
	cpu.a = result
}

// RRCA
func (cpu *CPU) rrca() {
	carry := cpu.a & 0x01
	result := (cpu.a >> 1) | (carry << 7)
	cpu.f = 0
	if carry != 0 {
		cpu.f |= FlagC
	}
	cpu.a = result
}

// RLA
func (cpu *CPU) rla() {
	carry := byte(0)
	if cpu.f&FlagC != 0 {
		carry = 1
	}
	newCarry := cpu.a >> 7
	cpu.a = (cpu.a << 1) | carry
	cpu.f = 0
	if newCarry != 0 {
		cpu.f |= FlagC
	}
}

// RRA
func (cpu *CPU) rra() {
	oldCarry := byte(0)
	if cpu.f&FlagC != 0 {
		oldCarry = 1
	}
	newCarry := cpu.a & 0x01
	cpu.a = (cpu.a >> 1) | (oldCarry << 7)
	cpu.f = 0
	if newCarry != 0 {
		cpu.f |= FlagC
	}
}
