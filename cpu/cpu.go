package cpu

import "fmt"

// Representa la CPU del Game Boy DMG
type CPU struct {
	regs registers

	memory [0x10000]byte // Memoria de 64 KB
	cycles int
}

// Crea e inicializa una nueva CPU
func NewCPU() *CPU {
	cpu := &CPU{}
	cpu.regs.a = 0x01
	cpu.regs.pc = 0x0100
	cpu.regs.sp = 0xFFFE
	return cpu
}

func (cpu *CPU) step() {
	opcode := cpu.memory[cpu.regs.pc]
	fmt.Printf("PC: %04X  Opcode: %02X\n", cpu.regs.pc, opcode)
	cpu.regs.pc++

	switch opcode {
	case 0x00: // NOP
		cpu.nop(4)

	case 0x01: // LD BC, n16
		lo := cpu.memory[cpu.regs.pc]
		hi := cpu.memory[cpu.regs.pc+1]
		value := uint16(hi)<<8 | uint16(lo)
		cpu.ld16(cpu.setBC, value, 12)

	case 0x02: // LD (BC),A
		cpu.ld8(&cpu.memory[cpu.getBC()], cpu.regs.a, 8)

	case 0x03: // INC BC
		cpu.inc16(cpu.setBC, cpu.getBC, 8)

	case 0x04: // INC B
		cpu.inc8(&cpu.regs.b, 4)

	case 0x05: // DEC B
		cpu.dec8(&cpu.regs.b, 4)

	case 0x06: // LD B,n8
		cpu.ld8(&cpu.regs.b, cpu.memory[cpu.regs.pc], 8)
		// TODO: Continuar revisando las opcode
	case 0x07: // RLCA
		cpu.regs.a = cpu.rlca(cpu.regs.a)
		cpu.cycles += 4

	case 0x08: // LD (a16),SP
		lo := cpu.memory[cpu.regs.pc]
		hi := cpu.memory[cpu.regs.pc+1]
		addr := uint16(hi)<<8 | uint16(lo)
		cpu.memory[addr] = byte(cpu.regs.sp & 0xFF)
		cpu.memory[addr+1] = byte(cpu.regs.sp >> 8)
		cpu.regs.pc += 2
		cpu.cycles += 20

	case 0x09: // ADD HL,BC
		hl := cpu.GetHL()
		result := uint32(hl) + uint32(cpu.GetBC())
		cpu.SetHL(uint16(result))
		cpu.updateAdd16Flags(hl, cpu.GetBC())
		cpu.cycles += 8

	case 0x0A: // LD A,(BC)
		cpu.regs.a = cpu.memory[cpu.GetBC()]
		cpu.cycles += 8

	case 0x0B: // DEC BC
		cpu.SetBC(cpu.GetBC() - 1)
		cpu.cycles += 8

	case 0x0C: // INC C
		cpu.regs.c = cpu.inc8(cpu.regs.c)
		cpu.cycles += 4

	case 0x0D: // DEC C
		cpu.regs.c = cpu.dec8(cpu.regs.c)
		cpu.cycles += 4

	case 0x0E: // LD C,d8
		cpu.regs.c = cpu.memory[cpu.regs.pc]
		cpu.regs.pc++
		cpu.cycles += 8

	case 0x0F: // RRCA
		cpu.regs.a = cpu.rrca(cpu.regs.a)
		cpu.cycles += 4

	default:
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		panic("Detenido")
	}
}
