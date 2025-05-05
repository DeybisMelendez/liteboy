package cpu

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
	return cpu
}

func (cpu *CPU) LoadMemory(memory *[]byte) {
	cpu.memory = *memory
}
