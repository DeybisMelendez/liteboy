package cpu

// Step ejecuta una instrucción del procesador y devuelve los t-ciclos utilizados
func (cpu *CPU) Step() int {
	cpu.tCycles = 0
	interruptsPending := (cpu.bus.Read(0xFF0F) & cpu.bus.Read(0xFFFF)) != 0

	if cpu.halted {
		if interruptsPending {
			cpu.halted = false
		} else {
			cpu.tick()
			// Si no hay interrupciones, CPU sigue halted, hace "nada"
			return cpu.tCycles
		}

	}

	// Verificar si debe manejar interrupciones
	if cpu.ime && interruptsPending {
		cpu.handleInterrupt()
		return cpu.tCycles
	}

	// Interrupciones se habilitan después de la instrucción siguiente al EI
	if cpu.enableIME {
		cpu.ime = true
		cpu.enableIME = false
	}

	// Fetch
	opcode := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.tick()
	// Decode, Execute
	cpu.execute(opcode)

	return cpu.tCycles
}

func (cpu *CPU) tick() {
	cpu.tCycles += 4
	cpu.ppu.Step(4)
	cpu.timer.Step(4)
}
