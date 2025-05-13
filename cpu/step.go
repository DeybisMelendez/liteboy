package cpu

// Step ejecuta una instrucción del procesador y devuelve los ciclos utilizados
func (cpu *CPU) Step() int {
	interruptsPending := (cpu.bus.Read(0xFF0F) & cpu.bus.Read(0xFFFF)) != 0

	if cpu.halted {
		if interruptsPending {
			cpu.halted = false
			if !cpu.ime {
				// HALT bug: el siguiente opcode debe ejecutarse de nuevo
				cpu.updateTimers(4)
				return 1
			}
		} else {
			cpu.updateTimers(4)
			// Si no hay interrupciones, CPU sigue halted, hace "nada"
			return 1
		}
	}

	// Interrupciones se habilitan después de la instrucción siguiente al EI
	if cpu.enableIME {
		cpu.ime = true
		cpu.enableIME = false
	}

	// Verificar si debe manejar interrupciones
	if cpu.ime && interruptsPending {
		cpu.handleInterrupt()
		return 5
	}

	// Fetch
	opcode := cpu.bus.Read(cpu.pc)
	cpu.pc++

	// Ejecutar instrucción
	cycles := cpu.execute(opcode)

	// Actualizar timers (usa ciclos en T-cycles)
	cpu.updateTimers(cycles * 4)

	return cycles
}
