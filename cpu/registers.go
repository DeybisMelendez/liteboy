package cpu

type registers struct {
	a, f byte
	b, c byte
	d, e byte
	h, l byte

	pc uint16 // Program Counter
	sp uint16 // Stack Pointer
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
