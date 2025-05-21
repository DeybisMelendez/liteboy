package cpu

// ret return from subroutine suma 12 tcycles
func (cpu *CPU) ret() {
	lo := cpu.bus.Read(cpu.sp)
	cpu.tick()
	hi := cpu.bus.Read(cpu.sp + 1)
	cpu.tick()
	cpu.sp += 2
	cpu.pc = uint16(hi)<<8 | uint16(lo)
	cpu.tick()
}

// suma 12 tcycles
func (cpu *CPU) call16(addr uint16) {
	cpu.tick() // Internal Delay
	cpu.sp -= 1
	cpu.bus.Write(cpu.sp, byte(cpu.pc>>8)) // PC high byte
	cpu.tick()
	cpu.sp -= 1
	cpu.bus.Write(cpu.sp, byte(cpu.pc&0xFF)) // PC low byte
	cpu.tick()
	cpu.pc = addr
}
