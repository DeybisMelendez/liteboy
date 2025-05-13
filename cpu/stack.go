package cpu

func (cpu *CPU) pop16(set func(uint16)) {
	lo := cpu.bus.Read(cpu.sp)
	hi := cpu.bus.Read(cpu.sp + 1)
	cpu.sp += 2
	value := uint16(hi)<<8 | uint16(lo)
	set(value)
}
func (cpu *CPU) push16(value uint16) {
	cpu.sp -= 2
	cpu.bus.Write(cpu.sp, byte(value&0xFF)) // LSB
	cpu.bus.Write(cpu.sp+1, byte(value>>8)) // MSB
}
