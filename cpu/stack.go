package cpu

const validFlagsMask = FlagZ | FlagN | FlagH | FlagC

// suma 8 tcycles //TODO: ¿o 12?
func (cpu *CPU) pop16(set func(uint16)) {
	lo := cpu.bus.Read(cpu.sp)
	cpu.sp++
	cpu.tick()
	hi := cpu.bus.Read(cpu.sp)
	cpu.sp++
	cpu.tick()
	value := uint16(hi)<<8 | uint16(lo)
	set(value)
	//cpu.tick()
}

// suma 12 tcycles
func (cpu *CPU) push16(value uint16) {
	cpu.tick() // Internal Delay
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(value>>8))
	cpu.tick()
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(value&0xFF))
	cpu.tick() // Espera extra
}

// suma 12 tcycles
func (cpu *CPU) pushAF() {
	af := uint16(cpu.a)<<8 | uint16(cpu.f&validFlagsMask)
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(af>>8)) // MSB (A)
	cpu.tick()
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(af&0xFF)) // LSB (F)
	cpu.tick()
	cpu.tick() // tick extra por restar 2 a sp
}

// suma 8 tcycles
func (cpu *CPU) popAF() {
	lo := cpu.bus.Read(cpu.sp)  // LSB (F)
	cpu.f = lo & validFlagsMask // Solo 4 bits altos válidos
	cpu.sp++
	cpu.tick()
	hi := cpu.bus.Read(cpu.sp) // MSB (A)
	cpu.a = hi
	cpu.sp++
	cpu.tick() // no hay tick extra por sumar
}
