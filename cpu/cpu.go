package cpu

// Instructions: https://gbdev.io/gb-opcodes/optables/
// Reference: https://rgbds.gbdev.io/docs/v0.9.2/gbz80.7

import (
	"log"

	"github.com/deybismelendez/liteboy/bus"
)

// Representa la CPU del Game Boy DMG
type CPU struct {
	a, f byte
	b, c byte
	d, e byte
	h, l byte

	pc           uint16 // Program Counter
	sp           uint16 // Stack Pointer
	halted       bool
	Stopped      bool
	ime          bool
	enableIME    bool
	divCounter   uint16
	timerCounter int
	tCycles      int
	bus          *bus.Bus
}

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
	cpu.Stopped = false
	cpu.bus = bus
	cpu.ime = false
	cpu.enableIME = false
	cpu.divCounter = 0
	cpu.timerCounter = 0
	return cpu
}

const (
	FlagZ byte = 1 << 7 // Zero
	FlagN byte = 1 << 6 // Subtract
	FlagH byte = 1 << 5 // Half Carry
	FlagC byte = 1 << 4 // Carry
)

func (cpu *CPU) Trace(opcode byte) {
	log.Printf("Opcode: %02X PC=%04X SP=%04X A=%02X B=%02X C=%02X D=%02X E=%02X F=%b H=%02X L=%02X",
		opcode, cpu.pc, cpu.sp, cpu.a, cpu.b, cpu.c, cpu.d, cpu.e, cpu.f, cpu.h, cpu.l)
	//fmt.Printf("\n0x%04X\t0x%02X", cpu.pc, opcode)
}
