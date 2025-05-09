package cpu

// Instructions: https://gbdev.io/gb-opcodes/optables/
// Reference: https://rgbds.gbdev.io/docs/v0.9.2/gbz80.7

import (
	"github.com/deybismelendez/liteboy/bus"
)

// Representa la CPU del Game Boy DMG
type CPU struct {
	a, f byte
	b, c byte
	d, e byte
	h, l byte

	pc        uint16 // Program Counter
	sp        uint16 // Stack Pointer
	cycles    int
	halted    bool
	ime       bool
	enableIME bool
	bus       *bus.Bus
}

// Crea e inicializa una nueva CPU
func NewCPU(bus *bus.Bus) *CPU {
	cpu := &CPU{}
	cpu.a = 0x01
	cpu.f = 0xB0
	cpu.b = 0x00
	cpu.c = 0x13
	cpu.d = 0x00
	cpu.e = 0xD8
	cpu.h = 0x01
	cpu.l = 0x4D
	cpu.pc = 0x0100
	cpu.sp = 0xFFFE
	cpu.halted = false
	cpu.bus = bus
	return cpu
}

func (cpu *CPU) GetCycles() int {
	return cpu.cycles
}
