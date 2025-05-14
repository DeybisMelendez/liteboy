package cpu

func (cpu *CPU) handleInterrupt() {
	IE := cpu.bus.Read(0xFFFF)
	IF := cpu.bus.Read(0xFF0F)
	pending := IE & IF

	for i := 0; i < 5; i++ {
		if (pending & (1 << i)) != 0 {
			// Limpia bit correspondiente en IF
			cpu.bus.Write(0xFF0F, IF&^(1<<i))

			cpu.ime = false

			// Push PC a la pila
			pc := cpu.pc
			cpu.pushPC(byte((pc >> 8) & 0xFF))
			cpu.pushPC(byte(pc & 0xFF))

			// Salta al vector
			cpu.pc = uint16(0x40 + i*8)

			break
		}
	}
}
func (cpu *CPU) pushPC(value byte) {
	cpu.sp--
	cpu.bus.Write(cpu.sp, value)
}
