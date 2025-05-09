package cpu

func (cpu *CPU) setAddr16(addr uint16, value uint16) {
	cpu.bus.Write(addr, byte(value&0xFF)) // low
	cpu.bus.Write(addr+1, byte(value>>8)) // high
}

// getA8 devuelve una dirección de memoria agregando 0xFF00 al byte inmediato
func (cpu *CPU) getA8() uint16 {
	return 0xFF00 + uint16(cpu.bus.Read(cpu.pc+1))
}

// a16
func (cpu *CPU) getA16() uint16 {
	lo := cpu.bus.Read(cpu.pc + 1)
	hi := cpu.bus.Read(cpu.pc + 2)
	return uint16(lo) | uint16(hi)<<8
}

// e8
func (cpu *CPU) getE8() int8 {
	return int8(cpu.bus.Read(cpu.pc + 1))
}

// n8
func (cpu *CPU) getN8() byte {
	return cpu.bus.Read(cpu.pc + 1)
}

// n16
func (cpu *CPU) getN16() uint16 {
	lo := cpu.bus.Read(cpu.pc + 1)
	hi := cpu.bus.Read(cpu.pc + 2)
	return uint16(lo) | uint16(hi)<<8
}

// AF combinado
func (cpu *CPU) getAF() uint16 {
	return uint16(cpu.a)<<8 | uint16(cpu.f&0xF0) // Solo los 4 bits superiores de F son válidos
}

func (cpu *CPU) setAF(value uint16) {
	cpu.a = byte(value >> 8)
	cpu.f = byte(value & 0xF0) // Solo se permiten los bits de flags (Z, N, H, C)
}

// BC combinado
func (cpu *CPU) getBC() uint16 {
	return uint16(cpu.b)<<8 | uint16(cpu.c)
}

func (cpu *CPU) setBC(value uint16) {
	cpu.b = byte(value >> 8)
	cpu.c = byte(value & 0xFF)
}

// DE combinado
func (cpu *CPU) getDE() uint16 {
	return uint16(cpu.d)<<8 | uint16(cpu.e)
}

func (cpu *CPU) setDE(value uint16) {
	cpu.d = byte(value >> 8)
	cpu.e = byte(value & 0xFF)
}

// HL combinado
func (cpu *CPU) getHL() uint16 {
	return uint16(cpu.h)<<8 | uint16(cpu.l)
}

func (cpu *CPU) setHL(value uint16) {
	cpu.h = byte(value >> 8)
	cpu.l = byte(value & 0xFF)
}
