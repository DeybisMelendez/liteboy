package cpu

import (
	"github.com/deybismelendez/liteboy/bus"
)

// Representa la CPU del Game Boy DMG
type CPU struct {
	regs   registers
	cycles int
	halted bool
	bus    *bus.Bus
}

// Crea e inicializa una nueva CPU
func NewCPU(bus *bus.Bus) *CPU {
	cpu := &CPU{}
	cpu.regs.a = 0x01
	cpu.regs.f = 0x00
	cpu.regs.b = 0xFF
	cpu.regs.c = 0x13
	cpu.regs.d = 0x00
	cpu.regs.e = 0xC1
	cpu.regs.h = 0x84
	cpu.regs.l = 0x03
	cpu.regs.pc = 0x0100
	cpu.regs.sp = 0xFFFE
	cpu.halted = false
	cpu.bus = bus
	return cpu
}

func (cpu *CPU) GetCycles() int {
	return cpu.cycles
}
