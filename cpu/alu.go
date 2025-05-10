package cpu

// INC 8 bits Z 0 H -
func (cpu *CPU) inc8(value *byte) {
	result := *value + 1
	cpu.f &^= FlagN // N = 0
	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
	if (*value&0x0F)+1 > 0x0F {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	*value = result
}

// INC Address 8 bits Z 0 H -
func (cpu *CPU) inc8Address(addr uint16) {
	value := cpu.bus.Read(addr)
	result := value + 1
	cpu.f &^= FlagN // N = 0
	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
	if (value&0x0F)+1 > 0x0F {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	cpu.bus.Write(addr, value)
}

// INC 16 bits - - - -
func (cpu *CPU) inc16(set func(uint16), value uint16) {
	set(value + 1)
}

func (cpu *CPU) dec8(r *byte) {
	old := *r
	*r--

	cpu.f &= FlagC // solo preserva Carry
	cpu.f |= FlagN // DEC siempre activa N

	if *r == 0 {
		cpu.f |= FlagZ
	}
	if (old & 0x0F) == 0 {
		cpu.f |= FlagH // half carry de 0x10 a 0x0F
	}
}

// DEC 8 bits Z 1 H -
func (cpu *CPU) dec8Address(addr uint16) {
	value := cpu.bus.Read(addr)
	result := value - 1
	cpu.f |= FlagN // N = 1
	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
	if (value & 0x0F) == 0x00 {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	cpu.bus.Write(addr, value)
}

// DEC 16 bits - - - -
func (cpu *CPU) dec16(set func(uint16), value uint16) {
	set(value - 1)
}

// Add 16 bits - 0 H C
func (cpu *CPU) add16(set func(uint16), a uint16, b uint16) {
	set(a + b)
	cpu.f &^= FlagN // N = 0
	if ((a & 0x0FFF) + (b & 0x0FFF)) > 0x0FFF {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	if uint32(a)+uint32(b) > 0xFFFF {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}
}

// Add 8 bits Z 0 H C
func (cpu *CPU) add8(a *byte, b byte) {
	result := *a + b
	cpu.f &^= FlagN // N = 0
	if ((*a & 0x0F) + (b & 0x0F)) > 0x0F {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	if uint16(*a)+uint16(b) > 0xFF {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}
	*a = result
	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
}

// Resta 8 Bits Z 1 H C
func (cpu *CPU) sub8(a *byte, b byte) {
	cpu.f |= FlagN // N = 1
	if (*a & 0x0F) < (b & 0x0F) {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	if *a < b {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}
	result := *a - b
	*a = result
	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
}

// SBC Sub with carry 8 bits Z 1 H C
func (cpu *CPU) sbc8(a *byte, b byte) {
	carry := byte(0)
	if cpu.f&FlagC != 0 {
		carry = 1
	}
	cpu.f |= FlagN // N = 1
	if (*a & 0x0F) < ((b & 0x0F) + carry) {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	if uint16(*a) < uint16(b)+uint16(carry) {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}
	result := *a - b - carry
	*a = result
	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
}

// ADC Add with carry 8 bits Z 0 H C
func (cpu *CPU) adc8(a *byte, b byte) {
	carry := byte(0)
	if cpu.f&FlagC != 0 {
		carry = 1
	}
	sum := uint16(*a) + uint16(b) + uint16(carry)
	cpu.f &^= FlagN // N = 0

	if ((*a & 0x0F) + (b & 0x0F) + carry) > 0x0F {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	if sum > 0xFF {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}

	*a = byte(sum)
	if *a == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
}

// Or 8 Z 0 0 0
func (cpu *CPU) or8(a *byte, b byte) {
	*a |= b
	cpu.f = 0
	if *a == 0 {
		cpu.f |= FlagZ
	}
}

// Xor 8 bits Z 0 0 0
func (cpu *CPU) xor8(a *byte, b byte) {
	*a ^= b
	cpu.f = 0
	if *a == 0 {
		cpu.f |= FlagZ
	}
}

// cp Compare 8 bits Z 1 H C
func (cpu *CPU) cp8(a byte, b byte) {
	cpu.f |= FlagN // N = 1
	if (a & 0x0F) < (b & 0x0F) {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	if a < b {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}
	if a == b {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
}

// and Z 0 1 0
func (cpu *CPU) and8(a *byte, b byte) {
	*a &= b
	cpu.f = 0
	if *a == 0 {
		cpu.f |= FlagZ
	}
	cpu.f |= FlagH // H siempre se activa en AND
}

// DAA (Decimal Adjust Accumulator) Z - 0 C
func (cpu *CPU) daa() {
	a := cpu.a
	var adjust byte = 0
	carry := false

	if (cpu.f&FlagH) != 0 || ((cpu.f&FlagN) == 0 && (a&0x0F) > 9) {
		adjust |= 0x06
	}
	if (cpu.f&FlagC) != 0 || ((cpu.f&FlagN) == 0 && a > 0x99) {
		adjust |= 0x60
		carry = true
	}

	if (cpu.f & FlagN) != 0 {
		a -= adjust
	} else {
		a += adjust
	}

	cpu.a = a
	cpu.f &^= FlagZ | FlagH | FlagC

	if a == 0 {
		cpu.f |= FlagZ
	}
	if carry {
		cpu.f |= FlagC
	}
}

// SCF (Set Carry Flag)
func (cpu *CPU) scf() {
	cpu.f &^= FlagN | FlagH
	cpu.f |= FlagC
}
