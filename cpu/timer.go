package cpu

const (
	DIVRegister  = 0xFF04
	TIMARegister = 0xFF05
	TMARegister  = 0xFF06
	TACRegister  = 0xFF07
)

func (cpu *CPU) updateTimers(tCycles int) {
	cpu.tCycles += tCycles
	// --- DIV siempre avanza a 16384 Hz (cada 256 ciclos de CPU) ---
	cpu.divCounter += uint16(tCycles)
	if cpu.divCounter >= 256 {
		cpu.divCounter -= 256
		div := cpu.bus.Read(DIVRegister)
		cpu.bus.Write(DIVRegister, div+1) // incrementar DIV
	}

	// --- TIMA controlado por TAC ---
	tac := cpu.bus.Read(TACRegister)

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
	cpu.timerCounter += tCycles
	for cpu.timerCounter >= threshold {
		cpu.timerCounter -= threshold

		tima := cpu.bus.Read(TIMARegister)
		if tima == 0xFF {
			// Desbordamiento: TIMA = TMA, IF |= 0x04
			cpu.bus.Write(TIMARegister, cpu.bus.Read(TMARegister))
			IFRegister := cpu.bus.Read(0xFF0F)
			cpu.bus.Write(0xFF0F, IFRegister|0x04)
		} else {
			// Solo incrementa
			cpu.bus.Write(TIMARegister, tima+1)
		}
	}
}
