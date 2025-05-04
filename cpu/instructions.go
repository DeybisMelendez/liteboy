package cpu

// NOP
func (cpu *CPU) nop(cycles int) {
	cpu.cycles += cycles
}

// INC 8 bits
func (cpu *CPU) inc8(value *byte, cycles int) {
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
	cpu.cycles += cycles
}

// INC 16 bits
func (cpu *CPU) inc16(set func(uint16), get func() uint16, cycles int) {
	set(get() + 1)
	cpu.cycles += cycles
}

func (cpu *CPU) dec8(value *byte, cycles int) {
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
	cpu.cycles += cycles
}
func (cpu *CPU) rlca(value byte) byte {
	carry := value >> 7
	result := (value << 1) | carry
	cpu.regs.f = 0
	if carry != 0 {
		cpu.regs.f |= FlagC
	}
	return result
}
func (cpu *CPU) rrca(value byte) byte {
	carry := value & 0x01
	result := (value >> 1) | (carry << 7)
	cpu.regs.f = 0
	if carry != 0 {
		cpu.regs.f |= FlagC
	}
	return result
}

func (cpu *CPU) updateAdd16Flags(a, b uint16) {
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

// Load Data 8 bits
func (cpu *CPU) ld8(set *byte, value byte, cycles int) {
	*set = value
	cpu.cycles += cycles
}

// Load Data 16 bits
func (cpu *CPU) ld16(set func(uint16), value uint16, cycles int) {
	set(value)
	cpu.cycles += cycles
	cpu.regs.pc += 2
}
