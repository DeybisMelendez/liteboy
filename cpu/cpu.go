package cpu

import "fmt"

// Representa la CPU del Game Boy DMG
type CPU struct {
	regs   registers
	memory []byte
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

func (cpu *CPU) LoadMemory(memory *[]byte) {
	cpu.memory = *memory
}

func (cpu *CPU) Step() {
	opcode := cpu.memory[cpu.regs.pc]
	fmt.Printf("PC: %04X  Opcode: %02X\n", cpu.regs.pc, opcode)
	//cpu.regs.pc++

	switch opcode {
	case 0x00: // NOP
		cpu.next(1, 1)

	case 0x01: // LD BC, n16
		cpu.ld16(cpu.setBC, cpu.getN16())
		cpu.next(3, 3)

	case 0x02: // LD (BC),A
		cpu.ld8(cpu.getAddr(cpu.getBC()), cpu.regs.a)
		cpu.next(1, 2)

	case 0x03: // INC BC
		cpu.inc16(cpu.setBC, cpu.getBC())
		cpu.next(1, 2)

	case 0x04: // INC B
		cpu.inc8(&cpu.regs.b)
		cpu.next(1, 1)

	case 0x05: // DEC B
		cpu.dec8(&cpu.regs.b)
		cpu.next(1, 1)

	case 0x06: // LD B,n8
		cpu.ld8(&cpu.regs.b, cpu.getN8())
		cpu.next(2, 2)
	case 0x07: // RLCA
		cpu.rlca(&cpu.regs.a)
		cpu.next(1, 1)

	case 0x08: // LD (a16),SP
		cpu.ld16(cpu.setAddr16(cpu.getA16()), cpu.regs.sp)
		cpu.next(3, 5)

	case 0x09: // ADD HL,BC
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.getBC())
		cpu.next(1, 2)

	case 0x0A: // LD A,(BC)
		cpu.ld8(&cpu.regs.a, *cpu.getAddr(cpu.getBC()))
		cpu.next(1, 2)

	case 0x0B: // DEC BC
		cpu.dec16(cpu.setBC, cpu.getBC())
		cpu.next(1, 2)

	case 0x0C: // INC C
		cpu.inc8(&cpu.regs.c)
		cpu.next(1, 1)

	case 0x0D: // DEC C
		cpu.dec8(&cpu.regs.c)
		cpu.next(1, 1)

	case 0x0E: // LD C,n8
		cpu.ld8(&cpu.regs.c, cpu.getN8())
		cpu.next(2, 2)

	case 0x0F: // RRCA
		cpu.rrca(&cpu.regs.a)
		cpu.next(1, 1)
	case 0x10: // STOP
		// STOP 0 instruction (detiene el reloj del sistema)
		// El siguiente byte debe ser 0x00, pero normalmente se ignora
		cpu.next(2, 1)
		// TODO: Agregar la lógica del modo STOP reloj/divider
		fmt.Println("Alerta: STOP todavía no está implementado")

	case 0x11: // LD DE, n16
		cpu.ld16(cpu.setDE, cpu.getN16())
		cpu.next(3, 3)

	case 0x12: // LD (DE), A
		cpu.ld8(cpu.getAddr(cpu.getDE()), cpu.regs.a)
		cpu.next(1, 2)

	case 0x13: // INC DE
		cpu.inc16(cpu.setDE, cpu.getDE())
		cpu.next(1, 2)

	case 0x14: // INC D
		cpu.inc8(&cpu.regs.d)
		cpu.next(1, 1)

	case 0x15: // DEC D
		cpu.dec8(&cpu.regs.d)
		cpu.next(1, 1)

	case 0x16: // LD D,n8
		cpu.ld8(&cpu.regs.d, cpu.getN8())
		cpu.next(2, 2)

	case 0x17: // RLA
		cpu.rla()
		cpu.next(1, 1)

	case 0x18: // JR e8
		offset := cpu.getE8()
		cpu.regs.pc = uint16(int(cpu.regs.pc) + 2 + int(offset))
		cpu.cycles += 3

	case 0x19: // ADD HL,DE
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.getDE())
		cpu.next(1, 2)

	case 0x1A: // LD A,(DE)
		cpu.ld8(&cpu.regs.a, *cpu.getAddr(cpu.getDE()))
		cpu.next(1, 2)

	case 0x1B: // DEC DE
		cpu.dec16(cpu.setDE, cpu.getDE())
		cpu.next(1, 2)

	case 0x1C: // INC E
		cpu.inc8(&cpu.regs.e)
		cpu.next(1, 1)

	case 0x1D: // DEC E
		cpu.dec8(&cpu.regs.e)
		cpu.next(1, 1)

	case 0x1E: // LD E,n8
		cpu.ld8(&cpu.regs.e, cpu.getN8())
		cpu.next(2, 2)

	case 0x1F: // RRA
		cpu.rra()
		cpu.next(1, 1)
	case 0x20: // JR NZ,e8
		offset := cpu.getE8()
		if cpu.regs.f&FlagZ == 0 {
			cpu.regs.pc = uint16(int(cpu.regs.pc) + 2 + int(offset))
			cpu.cycles += 3
		} else {
			cpu.next(2, 2)
		}

	case 0x21: // LD HL,n16
		cpu.ld16(cpu.setHL, cpu.getN16())
		cpu.next(3, 3)

	case 0x22: // LD (HL+),A
		hl := cpu.getHL()
		cpu.ld8(cpu.getAddr(hl), cpu.regs.a)
		cpu.setHL(hl + 1)
		cpu.next(1, 2)

	case 0x23: // INC HL
		cpu.inc16(cpu.setHL, cpu.getHL())
		cpu.next(1, 2)

	case 0x24: // INC H
		cpu.inc8(&cpu.regs.h)
		cpu.next(1, 1)

	case 0x25: // DEC H
		cpu.dec8(&cpu.regs.h)
		cpu.next(1, 1)

	case 0x26: // LD H,n8
		cpu.ld8(&cpu.regs.h, cpu.getN8())
		cpu.next(2, 2)

	case 0x27: // DAA
		cpu.daa()
		cpu.next(1, 1)

	case 0x28: // JR Z,e8
		offset := cpu.getE8()
		if cpu.regs.f&FlagZ != 0 {
			cpu.regs.pc = uint16(int(cpu.regs.pc) + 2 + int(offset))
			cpu.cycles += 3
		} else {
			cpu.next(2, 2)
		}

	case 0x29: // ADD HL,HL
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.getHL())
		cpu.next(1, 2)

	case 0x2A: // LD A,(HL+)
		hl := cpu.getHL()
		cpu.ld8(&cpu.regs.a, *cpu.getAddr(hl))
		cpu.setHL(hl + 1)
		cpu.next(1, 2)

	case 0x2B: // DEC HL
		cpu.dec16(cpu.setHL, cpu.getHL())
		cpu.next(1, 2)

	case 0x2C: // INC L
		cpu.inc8(&cpu.regs.l)
		cpu.next(1, 1)

	case 0x2D: // DEC L
		cpu.dec8(&cpu.regs.l)
		cpu.next(1, 1)

	case 0x2E: // LD L,n8
		cpu.ld8(&cpu.regs.l, cpu.getN8())
		cpu.next(2, 2)

	case 0x2F: // CPL (Complement A)
		cpu.regs.a = ^cpu.regs.a
		cpu.regs.f |= FlagN | FlagH
		cpu.next(1, 1)
	case 0x30: // JR NC, e8
		offset := cpu.getE8()
		if cpu.regs.f&FlagC == 0 {
			cpu.regs.pc = uint16(int(cpu.regs.pc) + 2 + int(offset))
			cpu.cycles += 1
		} else {
			cpu.next(2, 1)
		}

	case 0x31: // LD SP, n16
		cpu.regs.sp = cpu.getN16()
		cpu.next(3, 3)

	case 0x32: // LD (HL-), A
		addr := cpu.getHL()
		cpu.ld8(&cpu.memory[addr], cpu.regs.a)
		cpu.setHL(addr - 1)
		cpu.next(2, 2)

	case 0x33: // INC SP
		cpu.regs.sp++
		cpu.next(2, 2)

	case 0x34: // INC (HL)
		cpu.inc8(&cpu.memory[cpu.getHL()])
		cpu.next(3, 3)

	case 0x35: // DEC (HL)
		cpu.dec8(&cpu.memory[cpu.getHL()])
		cpu.next(3, 3)

	case 0x36: // LD (HL), n8
		cpu.ld8(&cpu.memory[cpu.getHL()], cpu.getN8())
		cpu.next(3, 3)

	case 0x37: // SCF (Set Carry Flag)
		cpu.regs.f &^= FlagN | FlagH
		cpu.regs.f |= FlagC
		cpu.next(1, 1)

	case 0x38: // JR C, e8
		offset := cpu.getE8()
		if cpu.regs.f&FlagC != 0 {
			cpu.regs.pc = uint16(int(cpu.regs.pc) + 2 + int(offset))
			cpu.cycles += 3
		} else {
			cpu.next(2, 2)
		}

	case 0x39: // ADD HL, SP
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.regs.sp)
		cpu.next(2, 2)

	case 0x3A: // LD A, (HL-)
		hl := cpu.getHL()
		cpu.ld8(&cpu.regs.a, *cpu.getAddr(hl))
		cpu.dec16(cpu.setHL, hl)
		cpu.next(2, 2)

	case 0x3B: // DEC SP
		cpu.regs.sp--
		cpu.next(2, 2)

	case 0x3C: // INC A
		cpu.inc8(&cpu.regs.a)
		cpu.next(1, 1)

	case 0x3D: // DEC A
		cpu.dec8(&cpu.regs.a)
		cpu.next(1, 1)

	case 0x3E: // LD A, n8
		cpu.ld8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0x3F: // CCF
		cpu.ccf()
		cpu.next(1, 1)
	case 0xC3: // JP a16
		cpu.regs.pc = cpu.getA16()
		cpu.cycles += 3
	default:
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		panic("Detenido")
	}
}
