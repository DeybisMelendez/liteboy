package cpu

import "log"

// Representa la CPU del Game Boy DMG
type CPU struct {
	regs   registers
	memory []byte
	cycles int
	halted bool
}

// Crea e inicializa una nueva CPU
func NewCPU() *CPU {
	cpu := &CPU{}
	cpu.regs.a = 0x01
	cpu.regs.pc = 0x0100
	cpu.regs.sp = 0xFFFE
	cpu.halted = false
	cpu.memory = make([]byte, 0x10000) // 64 KB de espacio direccionable
	return cpu
}

func (cpu *CPU) LoadMemory(rom *[]byte) {
	if len(*rom) < 0x8000 {
		log.Fatalf("Error: la ROM es demasiado pequeña, se esperaban al menos 32KB, pero se recibieron %d bytes", len(*rom))
	}
	// Copiamos los primeros 32KB de la ROM en la memoria del CPU
	copy(cpu.memory[:0x8000], (*rom)[:0x8000])
}
