package cpu

// INC Register Z 0 H -
func (cpu *CPU) incR(r *byte) {
	val := *r
	result := byte(val + 1)

	cpu.f &^= FlagN // N = 0

	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}

	if (val&0x0F)+1 > 0x0F {
		cpu.f |= FlagH // hubo acarreo en el nibble bajo
	} else {
		cpu.f &^= FlagH
	}

	*r = result
}

// DEC Register Z 1 H -
func (cpu *CPU) decR(r *byte) {
	val := *r
	result := val - 1

	cpu.f |= FlagN // N = 1

	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}

	if (val & 0x0F) == 0x00 {
		cpu.f |= FlagH // hubo préstamo en el nibble bajo
	} else {
		cpu.f &^= FlagH
	}

	*r = result
}

// INC (HL) 8 bits Z 0 H -
func (cpu *CPU) incHL() {
	addr := cpu.getHL()
	value := cpu.bus.Read(addr)
	result := byte(value + 1)

	// Flags
	cpu.f &^= FlagN // N = 0
	if result == 0 {
		cpu.f |= FlagZ // Z
	} else {
		cpu.f &^= FlagZ
	}
	if (value&0x0F)+1 > 0x0F { // Half Carry
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	cpu.tick(4)
	cpu.bus.Write(addr, result)
}

// DEC (HL) 8 bits Z 1 H -
func (cpu *CPU) decHL() {
	addr := cpu.getHL()
	value := cpu.bus.Read(addr)
	result := byte(value - 1)

	// Flags
	cpu.f |= FlagN // N = 1
	if result == 0 {
		cpu.f |= FlagZ // Z
	} else {
		cpu.f &^= FlagZ
	}
	if (value & 0x0F) == 0x00 { // Half Borrow (H flag)
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	cpu.tick(4)
	cpu.bus.Write(addr, result)
}

// INC 16 bits - - - -
func (cpu *CPU) inc16(set func(uint16), value uint16) {
	set(value + 1)
}

// DEC 16 bits - - - -
func (cpu *CPU) dec16(set func(uint16), value uint16) {
	set(value - 1)
}

// Add 16 bits - 0 H C
func (cpu *CPU) addHL(b uint16) {
	a := cpu.getHL()
	cpu.ldHL(a + b)
	cpu.f &^= FlagN                             // N = 0
	if ((a & 0x0FFF) + (b & 0x0FFF)) > 0x0FFF { // H
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}
	if uint32(a)+uint32(b) > 0xFFFF { // C
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
	valA := *a
	result := int(valA) - int(b)

	cpu.f |= FlagN // Subtraction

	// Half-Carry: Borrow from bit 4
	if (valA & 0x0F) < (b & 0x0F) {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}

	// Carry: Borrow from full byte
	if valA < b {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}

	final := byte(result)
	*a = final

	// Zero flag
	if final == 0 {
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
	valA := *a
	sub := int(valA) - int(b) - int(carry)
	result := byte(sub)

	cpu.f |= FlagN // Subtract

	// Half-carry (borrow from bit 4)
	if (valA & 0x0F) < ((b & 0x0F) + carry) {
		cpu.f |= FlagH
	} else {
		cpu.f &^= FlagH
	}

	// Full borrow (carry flag)
	if uint16(valA) < uint16(b)+uint16(carry) {
		cpu.f |= FlagC
	} else {
		cpu.f &^= FlagC
	}

	*a = result

	// Zero flag
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
	result := byte(sum)

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

	*a = result

	if result == 0 {
		cpu.f |= FlagZ
	} else {
		cpu.f &^= FlagZ
	}
}

// Or 8 Z 0 0 0
func (cpu *CPU) or8(a *byte, b byte) {
	*a |= b
	cpu.f = 0 // N, H, C = 0 por definición

	if *a == 0 {
		cpu.f |= FlagZ
	}
}

// Xor 8 bits Z 0 0 0
func (cpu *CPU) xor8(a *byte, b byte) {
	*a ^= b
	cpu.f = 0 // Todos los flags se limpian en XOR

	if *a == 0 {
		cpu.f |= FlagZ
	}
}

// cp Compare 8 bits Z 1 H C
func (cpu *CPU) cp8(a byte, b byte) {
	cpu.f = FlagN // N = 1, los demás se limpian

	if (a & 0x0F) < (b & 0x0F) {
		cpu.f |= FlagH
	}
	if a < b {
		cpu.f |= FlagC
	}
	if a == b {
		cpu.f |= FlagZ
	}
}

// and Z 0 1 0
func (cpu *CPU) and8(a *byte, b byte) {
	*a &= b
	cpu.f = FlagH // H siempre se activa, N y C son 0

	if *a == 0 {
		cpu.f |= FlagZ
	}
}

// DAA (Decimal Adjust Accumulator) Z - 0 C
func (cpu *CPU) daa() {
	a := cpu.a
	var adjust byte = 0
	carry := false

	// Ajuste cuando el valor en A es mayor que 9 o hay acarreo en el nibble bajo
	if (cpu.f&FlagH) != 0 || ((cpu.f&FlagN) == 0 && (a&0x0F) > 9) {
		adjust |= 0x06
	}

	// Ajuste cuando el valor en A es mayor que 0x99 o hay acarreo en C
	if (cpu.f&FlagC) != 0 || ((cpu.f&FlagN) == 0 && a > 0x99) {
		adjust |= 0x60
		carry = true
	}

	// Si la bandera N está activa (resta) o no, ajusta el valor
	if (cpu.f & FlagN) != 0 {
		a -= adjust
	} else {
		a += adjust
	}

	cpu.a = a
	// Limpia los flags Z, H y C
	cpu.f &^= FlagZ | FlagH | FlagC

	// Si el resultado es 0, entonces setea la bandera Z
	if a == 0 {
		cpu.f |= FlagZ
	}

	// Si hubo acarreo, setea la bandera C
	if carry {
		cpu.f |= FlagC
	}
}

// SCF (Set Carry Flag) - 0 0 1
func (cpu *CPU) scf() {
	cpu.f &^= FlagN | FlagH // Clear N and H
	cpu.f |= FlagC          // Set C
}

// CPL - 1 1 -
func (cpu *CPU) cpl() {
	cpu.a = ^cpu.a
	cpu.f |= FlagN | FlagH // Set N and H
}
