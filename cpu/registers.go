package cpu

func (cpu *CPU) setAddr16(addr uint16, value uint16) {
	cpu.bus.Write(addr, byte(value&0xFF)) // low
	cpu.bus.Write(addr+1, byte(value>>8)) // high
}

// getA8 devuelve una dirección de memoria agregando 0xFF00 al byte inmediato
func (cpu *CPU) getA8() uint16 {
	a8 := 0xFF00 + uint16(cpu.bus.Read(cpu.pc))
	cpu.pc++
	return a8
}

// a16
func (cpu *CPU) getA16() uint16 {
	lo := cpu.bus.Read(cpu.pc)
	cpu.pc++
	hi := cpu.bus.Read(cpu.pc)
	cpu.pc++
	return uint16(lo) | uint16(hi)<<8
}

// e8
func (cpu *CPU) getE8() int8 {
	e8 := int8(cpu.bus.Read(cpu.pc))
	cpu.pc++
	return e8
}

// n8
func (cpu *CPU) getN8() byte {
	n8 := cpu.bus.Read(cpu.pc)
	cpu.pc++
	return n8
}

// n16
func (cpu *CPU) getN16() uint16 {
	lo := cpu.bus.Read(cpu.pc)
	cpu.pc++
	hi := cpu.bus.Read(cpu.pc)
	cpu.pc++
	return uint16(lo) | uint16(hi)<<8
}

// BC combinado
func (cpu *CPU) getBC() uint16 {
	return uint16(cpu.b)<<8 | uint16(cpu.c)
}

// DE combinado
func (cpu *CPU) getDE() uint16 {
	return uint16(cpu.d)<<8 | uint16(cpu.e)
}

// HL combinado
func (cpu *CPU) getHL() uint16 {
	return uint16(cpu.h)<<8 | uint16(cpu.l)
}
