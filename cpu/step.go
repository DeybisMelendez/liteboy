package cpu

// Step ejecuta una instrucción del procesador y devuelve los ciclos utilizados
func (cpu *CPU) Step() int {
	interruptsPending := (cpu.bus.Read(0xFF0F) & cpu.bus.Read(0xFFFF)) != 0

	if cpu.halted {
		if interruptsPending {
			cpu.halted = false
		} else {
			cpu.updateTimers(4)
			// Si no hay interrupciones, CPU sigue halted, hace "nada"
			return 4
		}

	}

	// Verificar si debe manejar interrupciones
	if cpu.ime && interruptsPending {
		cpu.handleInterrupt()
		cpu.updateTimers(20)
		return 20
	}

	// Interrupciones se habilitan después de la instrucción siguiente al EI
	if cpu.enableIME {
		cpu.ime = true
		cpu.enableIME = false
	}

	// Fetch
	opcode := cpu.bus.Read(cpu.pc)
	cpu.pc++

	// Ejecutar instrucción
	cycles := cpu.execute(opcode)

	// Actualizar timers
	cpu.updateTimers(cycles)

	return cycles
}
