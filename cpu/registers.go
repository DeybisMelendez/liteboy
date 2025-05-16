package cpu

// Escribe 16 bits en memoria, suma 8 tycles
func (cpu *CPU) setAddr16(addr uint16, value uint16) {
	cpu.bus.Write(addr, byte(value&0xFF)) // low
	cpu.tick()
	cpu.bus.Write(addr+1, byte(value>>8)) // high
	cpu.tick()
}

/*
getA8
devuelve una dirección de memoria agregando 0xFF00 al byte inmediato, suma 4 tcycles
*/
func (cpu *CPU) getA8() uint16 {
	a8 := 0xFF00 + uint16(cpu.bus.Read(cpu.pc))
	cpu.pc++
	cpu.tick()
	return a8
}

// a16 lee 16 bits inmediato, suma 8 tcycles
func (cpu *CPU) getA16() uint16 {
	lo := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.tick()
	hi := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.tick()
	return uint16(lo) | uint16(hi)<<8
}

// e8
func (cpu *CPU) getE8() int8 {
	e8 := int8(cpu.bus.Read(cpu.pc))
	cpu.pc++
	return e8
}

// n8 lee 8 bits inmediatos, no suma tcycles
func (cpu *CPU) getN8() byte {
	n8 := cpu.bus.Read(cpu.pc)
	cpu.pc++
	return n8
}

// n16 lee 16 bits inmediato, suma 8 tcycles
func (cpu *CPU) getN16() uint16 {
	lo := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.tick()
	hi := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.tick()
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
