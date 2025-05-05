package cpu

// Instructions: https://gbdev.io/gb-opcodes/optables/
// Reference: https://rgbds.gbdev.io/docs/v0.9.2/gbz80.7

// next actualiza el estado del cpu saltando al siguiente ciclo
func (cpu *CPU) next(bytes uint16, cycles int) {
	cpu.cycles += cycles
	cpu.regs.pc += bytes
}

// INC 8 bits
func (cpu *CPU) inc8(value *byte) {
	result := *value + 1
	cpu.regs.f &^= FlagN // N = 0
	if result == 0 {
		cpu.regs.f |= FlagZ
	} else {
		cpu.regs.f &^= FlagZ
	}
	if (*value&0x0F)+1 > 0x0F {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	*value++
}

// INC 16 bits
func (cpu *CPU) inc16(set func(uint16), value uint16) {
	set(value + 1)
}

// DEC 8 bits
func (cpu *CPU) dec8(value *byte) {
	result := *value - 1
	cpu.regs.f |= FlagN // N = 1
	if result == 0 {
		cpu.regs.f |= FlagZ
	} else {
		cpu.regs.f &^= FlagZ
	}
	if (*value & 0x0F) == 0x00 {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	*value--
}

// DEC 16 bits
func (cpu *CPU) dec16(set func(uint16), value uint16) {
	set(value + 1)
}

// RLCA
func (cpu *CPU) rlca(value *byte) {
	carry := *value >> 7
	result := (*value << 1) | carry
	cpu.regs.f = 0
	if carry != 0 {
		cpu.regs.f |= FlagC
	}
	*value = result
}

// RRCA
func (cpu *CPU) rrca(value *byte) {
	carry := *value & 0x01
	result := (*value >> 1) | (carry << 7)
	cpu.regs.f = 0
	if carry != 0 {
		cpu.regs.f |= FlagC
	}
	*value = result
}

// RLA
func (cpu *CPU) rla() {
	carry := byte(0)
	if cpu.regs.f&FlagC != 0 {
		carry = 1
	}
	newCarry := cpu.regs.a >> 7
	cpu.regs.a = (cpu.regs.a << 1) | carry
	cpu.regs.f = 0
	if newCarry != 0 {
		cpu.regs.f |= FlagC
	}
}

// RRA
func (cpu *CPU) rra() {
	oldCarry := byte(0)
	if cpu.regs.f&FlagC != 0 {
		oldCarry = 1
	}
	newCarry := cpu.regs.a & 0x01
	cpu.regs.a = (cpu.regs.a >> 1) | (oldCarry << 7)
	cpu.regs.f = 0
	if newCarry != 0 {
		cpu.regs.f |= FlagC
	}
}

// DAA (Decimal Adjust Accumulator)
func (cpu *CPU) daa() {
	a := cpu.regs.a
	var adjust byte = 0
	carry := false

	if (cpu.regs.f&FlagH) != 0 || ((cpu.regs.f&FlagN) == 0 && (a&0x0F) > 9) {
		adjust |= 0x06
	}
	if (cpu.regs.f&FlagC) != 0 || ((cpu.regs.f&FlagN) == 0 && a > 0x99) {
		adjust |= 0x60
		carry = true
	}

	if (cpu.regs.f & FlagN) != 0 {
		a -= adjust
	} else {
		a += adjust
	}

	cpu.regs.a = a
	cpu.regs.f &^= FlagZ | FlagH | FlagC

	if a == 0 {
		cpu.regs.f |= FlagZ
	}
	if carry {
		cpu.regs.f |= FlagC
	}
}

// Load Data 8 bits
func (cpu *CPU) ld8(set *byte, value byte) {
	*set = value
}

// Load Data 16 bits
func (cpu *CPU) ld16(set func(uint16), value uint16) {
	set(value)
}

// Add 16 bits
func (cpu *CPU) add16(set func(uint16), a uint16, b uint16) {
	set(a + b)
	cpu.regs.f &^= FlagN // N = 0
	if ((a & 0x0FFF) + (b & 0x0FFF)) > 0x0FFF {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	if uint32(a)+uint32(b) > 0xFFFF {
		cpu.regs.f |= FlagC
	} else {
		cpu.regs.f &^= FlagC
	}
}
