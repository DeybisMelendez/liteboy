package cpu

// ret return from subroutine
func (cpu *CPU) ret() {
	lo := cpu.bus.Read(cpu.sp)
	hi := cpu.bus.Read(cpu.sp + 1)
	cpu.sp += 2
	cpu.pc = uint16(hi)<<8 | uint16(lo)
}

func (cpu *CPU) call16(addr uint16) {
	cpu.sp -= 2
	cpu.bus.Write(cpu.sp, byte(cpu.pc&0xFF))
	cpu.bus.Write(cpu.sp+1, byte(cpu.pc>>8))
	cpu.pc = addr
}

// reset, jump to fixed address
func (cpu *CPU) rst16(addr uint16) {
	cpu.sp -= 2
	cpu.bus.Write(cpu.sp, byte(cpu.pc&0xFF))
	cpu.bus.Write(cpu.sp+1, byte(cpu.pc>>8))
	cpu.pc = addr
}
