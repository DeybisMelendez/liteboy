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

func (cpu *CPU) updateTimers(cycles int) {
	// --- DIV siempre avanza a 16384 Hz (cada 256 ciclos de CPU) ---
	cpu.divCounter += uint16(cycles)
	if cpu.divCounter >= 256 {
		cpu.divCounter -= 256
		div := cpu.bus.Read(0xFF04)
		cpu.bus.Write(0xFF04, div+1) // incrementar DIV
	}

	// --- TIMA controlado por TAC ---
	tac := cpu.bus.Read(0xFF07)

	timerEnabled := tac&0x04 != 0
	if !timerEnabled {
		return
	}

	// Obtenemos el número de ciclos por incremento según TAC bits 1-0
	var threshold int
	switch tac & 0x03 {
	case 0:
		threshold = 1024 // 4096 Hz
	case 1:
		threshold = 16 // 262144 Hz
	case 2:
		threshold = 64 // 65536 Hz
	case 3:
		threshold = 256 // 16384 Hz
	}

	// Sumamos ciclos al timer interno
	cpu.timerCounter += cycles
	for cpu.timerCounter >= threshold {
		cpu.timerCounter -= threshold

		tima := cpu.bus.Read(0xFF05)
		if tima == 0xFF {
			// Desbordamiento: TIMA = TMA, IF |= 0x04
			cpu.bus.Write(0xFF05, cpu.bus.Read(0xFF06))
			ifr := cpu.bus.Read(0xFF0F)
			cpu.bus.Write(0xFF0F, ifr|0x04)
		} else {
			// Solo incrementa
			cpu.bus.Write(0xFF05, tima+1)
		}
	}
}
