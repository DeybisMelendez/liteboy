package cpu

const validFlagsMask = FlagZ | FlagN | FlagH | FlagC

func (cpu *CPU) pop16(set func(uint16)) {
	lo := cpu.bus.Read(cpu.sp)
	cpu.sp++
	hi := cpu.bus.Read(cpu.sp)
	cpu.sp++
	value := uint16(hi)<<8 | uint16(lo)
	set(value)
}
func (cpu *CPU) push16(value uint16) {
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(value>>8))
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(value&0xFF))
}

func (cpu *CPU) pushAF() {
	af := uint16(cpu.a)<<8 | uint16(cpu.f&validFlagsMask)
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(af>>8)) // MSB (A)
	cpu.sp--
	cpu.bus.Write(cpu.sp, byte(af&0xFF)) // LSB (F)
}

func (cpu *CPU) popAF() {
	lo := cpu.bus.Read(cpu.sp) // LSB (F)
	cpu.sp++
	hi := cpu.bus.Read(cpu.sp) // MSB (A)
	cpu.sp++

	cpu.a = hi
	cpu.f = lo & validFlagsMask // Solo 4 bits altos válidos
}
