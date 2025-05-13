package cpu

func (cpu *CPU) rlca() {
	carry := cpu.a >> 7
	cpu.a = (cpu.a << 1) | carry

	cpu.f = 0 // RLCA siempre limpia Z, N, H
	if carry != 0 {
		cpu.f |= FlagC
	}
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
	// Guardamos el carry anterior
	oldCarry := byte(0)
	if cpu.f&FlagC != 0 {
		oldCarry = 1
	}

	// Calculamos el nuevo carry desde el bit 7
	newCarry := cpu.a >> 7

	// Rotamos A a la izquierda, insertando el viejo carry
	cpu.a = (cpu.a << 1) | oldCarry

	cpu.f = 0

	// Establecemos el nuevo carry si corresponde
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
func (cpu *CPU) rlc(value byte) byte {
	result := (value << 1) | (value >> 7)
	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	if value&0x80 != 0 {
		cpu.f |= 0x10 // C
	}
	return result
}

func (cpu *CPU) rrc(value byte) byte {
	result := (value >> 1) | (value << 7)
	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	if value&0x01 != 0 {
		cpu.f |= 0x10 // C
	}
	return result
}

func (cpu *CPU) rl(value byte) byte {
	carry := (value & 0x80) >> 7
	result := (value << 1) | (cpu.f>>4)&1 // old carry flag
	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	if carry == 1 {
		cpu.f |= 0x10 // C
	}
	return result
}
func (cpu *CPU) rr(value byte) byte {
	oldCarry := byte(0)
	if cpu.f&0x10 != 0 {
		oldCarry = 1
	}
	newCarry := value & 1
	result := (value >> 1) | (oldCarry << 7)

	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	if newCarry != 0 {
		cpu.f |= 0x10 // C
	}
	return result
}
func (cpu *CPU) bit(bitIndex uint8, value byte) {
	// Verifica si el bit está en 0 o 1
	if value&(1<<bitIndex) == 0 {
		cpu.f |= 0x80 // Z (si es 0)
	} else {
		cpu.f &= 0x7F // Limpiar la bandera Z
	}

	// N siempre es 0 para BIT
	// H siempre es 1 para BIT
	cpu.f |= 0x40
}

func (cpu *CPU) sla(value byte) byte {
	carry := (value >> 7) & 1
	result := value << 1

	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	if carry != 0 {
		cpu.f |= 0x10 // C
	}
	return result
}

func (cpu *CPU) sra(value byte) byte {
	carry := value & 1
	msb := value & 0x80
	result := (value >> 1) | msb

	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	if carry != 0 {
		cpu.f |= 0x10 // C
	}
	return result
}
func (cpu *CPU) swap(value byte) byte {
	result := (value >> 4) | (value << 4)
	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	return result
}

func (cpu *CPU) srl(value byte) byte {
	carry := value & 1
	result := value >> 1

	cpu.f = 0
	if result == 0 {
		cpu.f |= 0x80 // Z
	}
	if carry != 0 {
		cpu.f |= 0x10 // C
	}
	return result
}
func (cpu *CPU) res(bit uint, value byte) byte {
	return value & ^(1 << bit)
}
func (cpu *CPU) set(bit uint, value byte) byte {
	return value | (1 << bit)
}
