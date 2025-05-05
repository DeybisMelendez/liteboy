package cpu

// Instructions: https://gbdev.io/gb-opcodes/optables/
// Reference: https://rgbds.gbdev.io/docs/v0.9.2/gbz80.7

// next actualiza el estado del cpu saltando al siguiente ciclo
func (cpu *CPU) next(bytes uint16, cycles int) {
	cpu.cycles += cycles
	cpu.regs.pc += bytes
}

func (cpu *CPU) halt() {
	cpu.halted = true
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
	*value = result
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
	*value = result
}

// DEC 16 bits
func (cpu *CPU) dec16(set func(uint16), value uint16) {
	set(value - 1)
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

// CCF (Complement Carry Flag)
func (cpu *CPU) ccf() {
	if cpu.regs.f&FlagC != 0 {
		cpu.regs.f &^= FlagC
	} else {
		cpu.regs.f |= FlagC
	}
	cpu.regs.f &^= FlagN | FlagH
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

// Add 8 bits
func (cpu *CPU) add8(a *byte, b byte) {
	result := *a + b
	cpu.regs.f &^= FlagN // N = 0
	if ((*a & 0x0F) + (b & 0x0F)) > 0x0F {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	if uint16(*a)+uint16(b) > 0xFF {
		cpu.regs.f |= FlagC
	} else {
		cpu.regs.f &^= FlagC
	}
	*a = result
	if result == 0 {
		cpu.regs.f |= FlagZ
	} else {
		cpu.regs.f &^= FlagZ
	}
}

// Resta 8 Bits
func (cpu *CPU) sub8(a *byte, b byte) {
	cpu.regs.f |= FlagN // N = 1
	if (*a & 0x0F) < (b & 0x0F) {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	if *a < b {
		cpu.regs.f |= FlagC
	} else {
		cpu.regs.f &^= FlagC
	}
	result := *a - b
	*a = result
	if result == 0 {
		cpu.regs.f |= FlagZ
	} else {
		cpu.regs.f &^= FlagZ
	}
}

// SBC Sub with carry 8 bits
func (cpu *CPU) sbc8(a *byte, b byte) {
	carry := byte(0)
	if cpu.regs.f&FlagC != 0 {
		carry = 1
	}
	cpu.regs.f |= FlagN // N = 1
	if (*a & 0x0F) < ((b & 0x0F) + carry) {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	if uint16(*a) < uint16(b)+uint16(carry) {
		cpu.regs.f |= FlagC
	} else {
		cpu.regs.f &^= FlagC
	}
	result := *a - b - carry
	*a = result
	if result == 0 {
		cpu.regs.f |= FlagZ
	} else {
		cpu.regs.f &^= FlagZ
	}
}

// ADC Add with carry 8 bits
func (cpu *CPU) adc8(a *byte, b byte) {
	carry := byte(0)
	if cpu.regs.f&FlagC != 0 {
		carry = 1
	}
	sum := uint16(*a) + uint16(b) + uint16(carry)
	cpu.regs.f &^= FlagN // N = 0

	if ((*a & 0x0F) + (b & 0x0F) + carry) > 0x0F {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	if sum > 0xFF {
		cpu.regs.f |= FlagC
	} else {
		cpu.regs.f &^= FlagC
	}

	*a = byte(sum)
	if *a == 0 {
		cpu.regs.f |= FlagZ
	} else {
		cpu.regs.f &^= FlagZ
	}
}

// Or 8 bits
func (cpu *CPU) or8(a *byte, b byte) {
	*a |= b
	cpu.regs.f = 0
	if *a == 0 {
		cpu.regs.f |= FlagZ
	}
}

// Xor 8 bits
func (cpu *CPU) xor8(a *byte, b byte) {
	*a ^= b
	cpu.regs.f = 0
	if *a == 0 {
		cpu.regs.f |= FlagZ
	}
}

// cp Compare 8 bits
func (cpu *CPU) cp8(a byte, b byte) {
	cpu.regs.f |= FlagN // N = 1
	if (a & 0x0F) < (b & 0x0F) {
		cpu.regs.f |= FlagH
	} else {
		cpu.regs.f &^= FlagH
	}
	if a < b {
		cpu.regs.f |= FlagC
	} else {
		cpu.regs.f &^= FlagC
	}
	if a == b {
		cpu.regs.f |= FlagZ
	} else {
		cpu.regs.f &^= FlagZ
	}
}
func (cpu *CPU) and8(a *byte, b byte) {
	*a &= b
	cpu.regs.f = 0
	if *a == 0 {
		cpu.regs.f |= FlagZ
	}
	cpu.regs.f |= FlagH // H siempre se activa en AND
}

// ret return from subroutine
func (cpu *CPU) ret() {
	lo := cpu.memory[cpu.regs.sp]
	hi := cpu.memory[cpu.regs.sp+1]
	cpu.regs.sp += 2
	cpu.regs.pc = uint16(hi)<<8 | uint16(lo)
}

func (cpu *CPU) pop16(set func(uint16)) {
	lo := cpu.memory[cpu.regs.sp]
	hi := cpu.memory[cpu.regs.sp+1]
	cpu.regs.sp += 2
	value := uint16(hi)<<8 | uint16(lo)
	set(value)
}
func (cpu *CPU) call16(addr uint16) {
	cpu.regs.sp -= 2
	cpu.memory[cpu.regs.sp] = byte(cpu.regs.pc & 0xFF)
	cpu.memory[cpu.regs.sp+1] = byte(cpu.regs.pc >> 8)
	cpu.regs.pc = addr
}
func (cpu *CPU) push16(value uint16) {
	cpu.regs.sp -= 2
	cpu.memory[cpu.regs.sp] = byte(value & 0xFF)
	cpu.memory[cpu.regs.sp+1] = byte(value >> 8)
}

// reset, jump to fixed address
func (cpu *CPU) rst16(addr uint16) {
	cpu.regs.sp -= 2
	cpu.memory[cpu.regs.sp] = byte(cpu.regs.pc & 0xFF)
	cpu.memory[cpu.regs.sp+1] = byte(cpu.regs.pc >> 8)
	cpu.regs.pc = addr
}
func (cpu *CPU) ldh8(set *byte, value byte) {
	addr := 0xFF00 + uint16(cpu.getN8())
	cpu.ld8(cpu.getAddr(addr), value)
}
