package cpu

type registers struct {
	a, f byte
	b, c byte
	d, e byte
	h, l byte

	pc uint16 // Program Counter
	sp uint16 // Stack Pointer
}

// Memory Address
func (cpu *CPU) getAddr(value uint16) *byte {
	return &cpu.memory[value]
}

/*func (cpu *CPU) setAddr8(addr uint16) func(byte) {
	return func(value byte) {
		cpu.memory[addr] = value
	}
}*/

func (cpu *CPU) setAddr16(addr uint16) func(uint16) {
	return func(value uint16) {
		cpu.memory[addr] = byte(value & 0xFF)
		cpu.memory[addr+1] = byte(value >> 8)
	}
}

func (cpu *CPU) getA16() uint16 {
	lo := cpu.memory[cpu.regs.pc+1]
	hi := cpu.memory[cpu.regs.pc+2]
	return uint16(hi)<<8 | uint16(lo)
}

// e8
func (cpu *CPU) getE8() int8 {
	return int8(cpu.memory[cpu.regs.pc+1])
}

// n8
func (cpu *CPU) getN8() byte {
	return cpu.memory[cpu.regs.pc+1]
}

// n16
func (cpu *CPU) getN16() uint16 {
	lo := cpu.memory[cpu.regs.pc+1]
	hi := cpu.memory[cpu.regs.pc+2]
	return uint16(hi)<<8 | uint16(lo)
}

// AF combinado
func (cpu *CPU) getAF() uint16 {
	return uint16(cpu.regs.a)<<8 | uint16(cpu.regs.f&0xF0) // Solo los 4 bits superiores de F son válidos
}

func (cpu *CPU) setAF(value uint16) {
	cpu.regs.a = byte(value >> 8)
	cpu.regs.f = byte(value & 0xF0) // Solo se permiten los bits de flags (Z, N, H, C)
}

// BC combinado
func (cpu *CPU) getBC() uint16 {
	return uint16(cpu.regs.b)<<8 | uint16(cpu.regs.c)
}

func (cpu *CPU) setBC(value uint16) {
	cpu.regs.b = byte(value >> 8)
	cpu.regs.c = byte(value & 0xFF)
}

// DE combinado
func (cpu *CPU) getDE() uint16 {
	return uint16(cpu.regs.d)<<8 | uint16(cpu.regs.e)
}

func (cpu *CPU) setDE(value uint16) {
	cpu.regs.d = byte(value >> 8)
	cpu.regs.e = byte(value & 0xFF)
}

// HL combinado
func (cpu *CPU) getHL() uint16 {
	return uint16(cpu.regs.h)<<8 | uint16(cpu.regs.l)
}

func (cpu *CPU) setHL(value uint16) {
	cpu.regs.h = byte(value >> 8)
	cpu.regs.l = byte(value & 0xFF)
}
