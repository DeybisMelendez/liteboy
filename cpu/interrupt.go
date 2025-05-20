package cpu

func (cpu *CPU) handleInterrupt() {
	IE := cpu.bus.Read(0xFFFF)
	IF := cpu.bus.Read(0xFF0F)
	pending := IE & IF
	for i := range 5 {
		if (pending & (1 << i)) != 0 {

			cpu.tick() // 2 M cycles de espera
			cpu.tick()
			// Push PC a la pila
			pc := cpu.pc
			cpu.tick()
			cpu.pushPC(byte((pc >> 8) & 0xFF))
			// Interrupt dispatch cancellation via IE write during PC push
			// Pasa el test de mooneye acceptance/interrupts/ie_push
			if cpu.sp == 0xFFFF {
				cpu.ime = false
				cpu.pc = 0
				continue
			}
			cpu.tick()
			cpu.pushPC(byte(pc & 0xFF))

			// Limpia bit correspondiente en IF
			cpu.bus.Write(0xFF0F, IF&^(1<<i))

			// Salta al vector
			cpu.tick()
			cpu.pc = uint16(0x40 + i*8)

			// Desactivar IME
			cpu.ime = false

			break
		}
	}
}
func (cpu *CPU) pushPC(value byte) {
	cpu.sp--
	cpu.bus.Write(cpu.sp, value)
}
