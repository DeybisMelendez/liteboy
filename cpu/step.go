package cpu

import "fmt"

func (cpu *CPU) Step() int {
	opcode := cpu.bus.Read(cpu.regs.pc) //cpu.memory[cpu.regs.pc]
	fmt.Println("---")
	fmt.Printf("Opcode: %02X\n", opcode)
	fmt.Printf("Registers: pc: %04X sp: %04X\n", cpu.regs.pc, cpu.regs.sp)
	fmt.Printf("Registers: a: %04X b: %04X c: %04X d: %04X\n", cpu.regs.a, cpu.regs.b, cpu.regs.c, cpu.regs.d)
	fmt.Printf("Registers: e: %04X f: %04X h: %04X l: %04X\n", cpu.regs.e, cpu.regs.f, cpu.regs.h, cpu.regs.l)
	cycles := cpu.cycles
	if cpu.halted {
		// TODO: agregar CheckInterrupts() para cambiar halted a false
		fmt.Println("CPU Halted")
		cpu.cycles++
		return cycles - cpu.cycles
	}
	switch opcode {
	case 0x00: // NOP
		cpu.next(1, 1)

	case 0x01: // LD BC, n16
		cpu.ld16(cpu.setBC, cpu.getN16())
		cpu.next(3, 3)

	case 0x02: // LD (BC),A
		cpu.ld8(cpu.bus.GetAddress(cpu.getBC()), cpu.regs.a)
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

	case 0x21: // LD HL, n16
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
		cpu.ld8(cpu.bus.GetAddress(addr), cpu.regs.a) //&cpu.memory[addr], cpu.regs.a)
		cpu.setHL(addr - 1)
		cpu.next(2, 2)

	case 0x33: // INC SP
		cpu.regs.sp++
		cpu.next(2, 2)

	case 0x34: // INC (HL)
		cpu.inc8(cpu.bus.GetAddress(cpu.getHL())) //&cpu.memory[cpu.getHL()])
		cpu.next(3, 3)

	case 0x35: // DEC (HL)
		cpu.dec8(cpu.bus.GetAddress(cpu.getHL())) //&cpu.memory[cpu.getHL()])
		cpu.next(3, 3)

	case 0x36: // LD (HL), n8
		cpu.ld8(cpu.bus.GetAddress(cpu.getHL()), cpu.getN8()) //&cpu.memory[cpu.getHL()], cpu.getN8())
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
	case 0x40: // LD B,B
		cpu.ld8(&cpu.regs.b, cpu.regs.b)
		cpu.next(1, 1)
	case 0x41: // LD B,C
		cpu.ld8(&cpu.regs.b, cpu.regs.c)
		cpu.next(1, 1)
	case 0x42: // LD B,D
		cpu.ld8(&cpu.regs.b, cpu.regs.d)
		cpu.next(1, 1)
	case 0x43: // LD B,E
		cpu.ld8(&cpu.regs.b, cpu.regs.e)
		cpu.next(1, 1)
	case 0x44: // LD B,H
		cpu.ld8(&cpu.regs.b, cpu.regs.h)
		cpu.next(1, 1)
	case 0x45: // LD B,L
		cpu.ld8(&cpu.regs.b, cpu.regs.l)
		cpu.next(1, 1)
	case 0x46: // LD B,(HL)
		cpu.ld8(&cpu.regs.b, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x47: // LD B,A
		cpu.ld8(&cpu.regs.b, cpu.regs.a)
		cpu.next(1, 1)

	case 0x48: // LD C,B
		cpu.ld8(&cpu.regs.c, cpu.regs.b)
		cpu.next(1, 1)
	case 0x49: // LD C,C
		cpu.ld8(&cpu.regs.c, cpu.regs.c)
		cpu.next(1, 1)
	case 0x4A: // LD C,D
		cpu.ld8(&cpu.regs.c, cpu.regs.d)
		cpu.next(1, 1)
	case 0x4B: // LD C,E
		cpu.ld8(&cpu.regs.c, cpu.regs.e)
		cpu.next(1, 1)
	case 0x4C: // LD C,H
		cpu.ld8(&cpu.regs.c, cpu.regs.h)
		cpu.next(1, 1)
	case 0x4D: // LD C,L
		cpu.ld8(&cpu.regs.c, cpu.regs.l)
		cpu.next(1, 1)
	case 0x4E: // LD C,(HL)
		cpu.ld8(&cpu.regs.c, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x4F: // LD C,A
		cpu.ld8(&cpu.regs.c, cpu.regs.a)
		cpu.next(1, 1)
	case 0x50: // LD D,B
		cpu.ld8(&cpu.regs.d, cpu.regs.b)
		cpu.next(1, 1)
	case 0x51: // LD D,C
		cpu.ld8(&cpu.regs.d, cpu.regs.c)
		cpu.next(1, 1)
	case 0x52: // LD D,D
		cpu.ld8(&cpu.regs.d, cpu.regs.d)
		cpu.next(1, 1)
	case 0x53: // LD D,E
		cpu.ld8(&cpu.regs.d, cpu.regs.e)
		cpu.next(1, 1)
	case 0x54: // LD D,H
		cpu.ld8(&cpu.regs.d, cpu.regs.h)
		cpu.next(1, 1)
	case 0x55: // LD D,L
		cpu.ld8(&cpu.regs.d, cpu.regs.l)
		cpu.next(1, 1)
	case 0x56: // LD D,(HL)
		cpu.ld8(&cpu.regs.d, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x57: // LD D,A
		cpu.ld8(&cpu.regs.d, cpu.regs.a)
		cpu.next(1, 1)

	case 0x58: // LD E,B
		cpu.ld8(&cpu.regs.e, cpu.regs.b)
		cpu.next(1, 1)
	case 0x59: // LD E,C
		cpu.ld8(&cpu.regs.e, cpu.regs.c)
		cpu.next(1, 1)
	case 0x5A: // LD E,D
		cpu.ld8(&cpu.regs.e, cpu.regs.d)
		cpu.next(1, 1)
	case 0x5B: // LD E,E
		cpu.ld8(&cpu.regs.e, cpu.regs.e)
		cpu.next(1, 1)
	case 0x5C: // LD E,H
		cpu.ld8(&cpu.regs.e, cpu.regs.h)
		cpu.next(1, 1)
	case 0x5D: // LD E,L
		cpu.ld8(&cpu.regs.e, cpu.regs.l)
		cpu.next(1, 1)
	case 0x5E: // LD E,(HL)
		cpu.ld8(&cpu.regs.e, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x5F: // LD E,A
		cpu.ld8(&cpu.regs.e, cpu.regs.a)
		cpu.next(1, 1)
	case 0x60: // LD H,B
		cpu.ld8(&cpu.regs.h, cpu.regs.b)
		cpu.next(1, 1)
	case 0x61: // LD H,C
		cpu.ld8(&cpu.regs.h, cpu.regs.c)
		cpu.next(1, 1)
	case 0x62: // LD H,D
		cpu.ld8(&cpu.regs.h, cpu.regs.d)
		cpu.next(1, 1)
	case 0x63: // LD H,E
		cpu.ld8(&cpu.regs.h, cpu.regs.e)
		cpu.next(1, 1)
	case 0x64: // LD H,H
		cpu.ld8(&cpu.regs.h, cpu.regs.h)
		cpu.next(1, 1)
	case 0x65: // LD H,L
		cpu.ld8(&cpu.regs.h, cpu.regs.l)
		cpu.next(1, 1)
	case 0x66: // LD H,(HL)
		cpu.ld8(&cpu.regs.h, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x67: // LD H,A
		cpu.ld8(&cpu.regs.h, cpu.regs.a)
		cpu.next(1, 1)

	case 0x68: // LD L,B
		cpu.ld8(&cpu.regs.l, cpu.regs.b)
		cpu.next(1, 1)
	case 0x69: // LD L,C
		cpu.ld8(&cpu.regs.l, cpu.regs.c)
		cpu.next(1, 1)
	case 0x6A: // LD L,D
		cpu.ld8(&cpu.regs.l, cpu.regs.d)
		cpu.next(1, 1)
	case 0x6B: // LD L,E
		cpu.ld8(&cpu.regs.l, cpu.regs.e)
		cpu.next(1, 1)
	case 0x6C: // LD L,H
		cpu.ld8(&cpu.regs.l, cpu.regs.h)
		cpu.next(1, 1)
	case 0x6D: // LD L,L
		cpu.ld8(&cpu.regs.l, cpu.regs.l)
		cpu.next(1, 1)
	case 0x6E: // LD L,(HL)
		cpu.ld8(&cpu.regs.l, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x6F: // LD L,A
		cpu.ld8(&cpu.regs.l, cpu.regs.a)
		cpu.next(1, 1)

	case 0x70: // LD (HL),B
		cpu.bus.Write(cpu.getHL(), cpu.regs.b) //cpu.memory[cpu.getHL()] = cpu.regs.b
		cpu.next(1, 2)
	case 0x71: // LD (HL),C
		cpu.bus.Write(cpu.getHL(), cpu.regs.c) //cpu.memory[cpu.getHL()] = cpu.regs.c
		cpu.next(1, 2)
	case 0x72: // LD (HL),D
		cpu.bus.Write(cpu.getHL(), cpu.regs.d) //cpu.memory[cpu.getHL()] = cpu.regs.d
		cpu.next(1, 2)
	case 0x73: // LD (HL),E
		cpu.bus.Write(cpu.getHL(), cpu.regs.e) //cpu.memory[cpu.getHL()] = cpu.regs.e
		cpu.next(1, 2)
	case 0x74: // LD (HL),H
		cpu.bus.Write(cpu.getHL(), cpu.regs.h) //cpu.memory[cpu.getHL()] = cpu.regs.h
		cpu.next(1, 2)
	case 0x75: // LD (HL),L
		cpu.bus.Write(cpu.getHL(), cpu.regs.l) //cpu.memory[cpu.getHL()] = cpu.regs.l
		cpu.next(1, 2)
	case 0x76: // HALT
		cpu.halt()
		cpu.next(1, 1)
	case 0x77: // LD (HL),A
		cpu.bus.Write(cpu.getHL(), cpu.regs.a) //cpu.memory[cpu.getHL()] = cpu.regs.a
		cpu.next(1, 2)

	case 0x78: // LD A,B
		cpu.ld8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0x79: // LD A,C
		cpu.ld8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0x7A: // LD A,D
		cpu.ld8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0x7B: // LD A,E
		cpu.ld8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0x7C: // LD A,H
		cpu.ld8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0x7D: // LD A,L
		cpu.ld8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0x7E: // LD A,(HL)
		cpu.ld8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x7F: // LD A,A
		cpu.ld8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)
	case 0x80: // ADD A,B
		cpu.add8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0x81: // ADD A,C
		cpu.add8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0x82: // ADD A,D
		cpu.add8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0x83: // ADD A,E
		cpu.add8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0x84: // ADD A,H
		cpu.add8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0x85: // ADD A,L
		cpu.add8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0x86: // ADD A,(HL)
		cpu.add8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x87: // ADD A,A
		cpu.add8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)

	case 0x88: // ADC A,B
		cpu.adc8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0x89: // ADC A,C
		cpu.adc8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0x8A: // ADC A,D
		cpu.adc8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0x8B: // ADC A,E
		cpu.adc8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0x8C: // ADC A,H
		cpu.adc8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0x8D: // ADC A,L
		cpu.adc8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0x8E: // ADC A,(HL)
		cpu.adc8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x8F: // ADC A,A
		cpu.adc8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)

	case 0x90: // SUB A, B
		cpu.sub8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0x91: // SUB A, C
		cpu.sub8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0x92: // SUB A, D
		cpu.sub8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0x93: // SUB A, E
		cpu.sub8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0x94: // SUB A, H
		cpu.sub8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0x95: // SUB A, L
		cpu.sub8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0x96: // SUB A, (HL)
		cpu.sub8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x97: // SUB A, A
		cpu.sub8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)

	case 0x98: // SBC A,B
		cpu.sbc8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0x99: // SBC A,C
		cpu.sbc8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0x9A: // SBC A,D
		cpu.sbc8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0x9B: // SBC A,E
		cpu.sbc8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0x9C: // SBC A,H
		cpu.sbc8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0x9D: // SBC A,L
		cpu.sbc8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0x9E: // SBC A,(HL)
		cpu.sbc8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0x9F: // SBC A,A
		cpu.sbc8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)

	case 0xA0: // AND B
		cpu.and8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0xA1: // AND C
		cpu.and8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0xA2: // AND D
		cpu.and8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0xA3: // AND E
		cpu.and8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0xA4: // AND H
		cpu.and8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0xA5: // AND L
		cpu.and8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0xA6: // AND (HL)
		cpu.and8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0xA7: // AND A
		cpu.and8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)

	case 0xA8: // XOR B
		cpu.xor8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0xA9: // XOR C
		cpu.xor8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0xAA: // XOR D
		cpu.xor8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0xAB: // XOR E
		cpu.xor8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0xAC: // XOR H
		cpu.xor8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0xAD: // XOR L
		cpu.xor8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0xAE: // XOR (HL)
		cpu.xor8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0xAF: // XOR A
		cpu.xor8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)

	case 0xB0: // OR B
		cpu.or8(&cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0xB1: // OR C
		cpu.or8(&cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0xB2: // OR D
		cpu.or8(&cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0xB3: // OR E
		cpu.or8(&cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0xB4: // OR H
		cpu.or8(&cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0xB5: // OR L
		cpu.or8(&cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0xB6: // OR (HL)
		cpu.or8(&cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0xB7: // OR A
		cpu.or8(&cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)

	case 0xB8: // CP B
		cpu.cp8(cpu.regs.a, cpu.regs.b)
		cpu.next(1, 1)
	case 0xB9: // CP C
		cpu.cp8(cpu.regs.a, cpu.regs.c)
		cpu.next(1, 1)
	case 0xBA: // CP D
		cpu.cp8(cpu.regs.a, cpu.regs.d)
		cpu.next(1, 1)
	case 0xBB: // CP E
		cpu.cp8(cpu.regs.a, cpu.regs.e)
		cpu.next(1, 1)
	case 0xBC: // CP H
		cpu.cp8(cpu.regs.a, cpu.regs.h)
		cpu.next(1, 1)
	case 0xBD: // CP L
		cpu.cp8(cpu.regs.a, cpu.regs.l)
		cpu.next(1, 1)
	case 0xBE: // CP (HL)
		cpu.cp8(cpu.regs.a, *cpu.getAddr(cpu.getHL()))
		cpu.next(1, 2)
	case 0xBF: // CP A
		cpu.cp8(cpu.regs.a, cpu.regs.a)
		cpu.next(1, 1)
	case 0xC0: // RET NZ
		if cpu.regs.f&FlagZ == 0 {
			cpu.ret()
			cpu.cycles += 20
		} else {
			cpu.next(1, 8)
		}

	case 0xC1: // POP BC
		cpu.pop16(cpu.setBC)
		cpu.next(1, 12)

	case 0xC2: // JP NZ, a16
		if cpu.regs.f&FlagZ == 0 {
			cpu.regs.pc = cpu.getA16()
			cpu.cycles += 4
		} else {
			cpu.next(3, 3)
		}

	case 0xC3: // JP a16
		cpu.regs.pc = cpu.getA16()
		cpu.cycles += 4

	case 0xC4: // CALL NZ, a16
		if cpu.regs.f&FlagZ == 0 {
			cpu.call16(cpu.getA16())
			cpu.next(6, 4)
		} else {
			cpu.next(3, 3)
		}

	case 0xC5: // PUSH BC
		cpu.push16(cpu.getBC())
		cpu.next(1, 4)

	case 0xC6: // ADD A, n8
		cpu.add8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0xC7: // RST 00H
		cpu.rst16(0x00)

	case 0xC8: // RET Z
		if cpu.regs.f&FlagZ != 0 {
			cpu.ret()
			cpu.next(1, 5)
		} else {
			cpu.next(1, 2)
		}

	case 0xC9: // RET
		cpu.ret()
		cpu.next(1, 4)

	case 0xCA: // JP Z, a16
		if cpu.regs.f&FlagZ != 0 {
			cpu.regs.pc = cpu.getA16()
			cpu.cycles += 4
		} else {
			cpu.next(3, 3)
		}

	case 0xCB: // PREFIX CB
		cpu.next(1, 4)
		// TODO: Implementar Prefix $CB
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		panic("Detenido")

	case 0xCC: // CALL Z, a16
		if cpu.regs.f&FlagZ != 0 {
			cpu.call16(cpu.getA16())
			cpu.next(3, 6)
		} else {
			cpu.next(3, 3)
		}

	case 0xCD: // CALL a16
		cpu.call16(cpu.getA16())
		cpu.next(3, 6)

	case 0xCE: // ADC A, n8
		cpu.adc8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 8)

	case 0xCF: // RST 08H
		cpu.rst16(0x08)
		cpu.next(1, 4)
	case 0xD0: // RET NC
		if cpu.regs.f&FlagC == 0 {
			cpu.ret()
			cpu.next(1, 5)
		} else {
			cpu.next(1, 2)
		}

	case 0xD1: // POP DE
		cpu.pop16(cpu.setDE)
		cpu.next(1, 3)

	case 0xD2: // JP NC, a16
		if cpu.regs.f&FlagC == 0 {
			cpu.regs.pc = cpu.getA16()
			cpu.cycles += 4
		} else {
			cpu.next(3, 3)
		}

	case 0xD4: // CALL NC, a16
		if cpu.regs.f&FlagC == 0 {
			cpu.call16(cpu.getA16())
			cpu.next(3, 6)
		} else {
			cpu.next(3, 3)
		}

	case 0xD5: // PUSH DE
		cpu.push16(cpu.getDE())
		cpu.next(1, 4)

	case 0xD6: // SUB n8
		cpu.sub8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0xD7: // RST 10H
		cpu.rst16(0x10)
		cpu.next(1, 4)

	case 0xD8: // RET C
		if cpu.regs.f&FlagC != 0 {
			cpu.ret()
			cpu.next(1, 5)
		} else {
			cpu.next(1, 2)
		}

	case 0xD9: // RETI
		cpu.ret()
		// TODO: Habilitar interrputs en RETI
		cpu.next(1, 4)

	case 0xDA: // JP C, a16
		if cpu.regs.f&FlagC != 0 {
			cpu.regs.pc = cpu.getA16()
			cpu.cycles += 4
		} else {
			cpu.next(3, 3)
		}

	case 0xDC: // CALL C, a16
		if cpu.regs.f&FlagC != 0 {
			cpu.call16(cpu.getA16())
			cpu.next(1, 6)
		} else {
			cpu.next(3, 3)
		}

	case 0xDE: // SBC A, n8
		cpu.sbc8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0xDF: // RST 18H
		cpu.rst16(0x18)
		cpu.next(1, 4)
	case 0xE0: // LDH (n), A
		addr := 0xFF00 + uint16(cpu.getN8())
		*cpu.getAddr(addr) = cpu.regs.a
		cpu.next(2, 3)

	case 0xE1: // POP HL
		cpu.pop16(cpu.setHL)
		cpu.next(1, 3)

	case 0xE2: // LDH (C), A
		addr := 0xFF00 + uint16(cpu.regs.c)
		*cpu.getAddr(addr) = cpu.regs.a
		cpu.next(1, 2)

	case 0xE5: // PUSH HL
		cpu.push16(cpu.getHL())
		cpu.next(1, 16)

	case 0xE6: // AND A,n8
		cpu.add8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0xE8: // ADD SP, e8
		e8 := int16(cpu.getE8())
		sp := int16(cpu.regs.sp)
		result := sp + e8
		cpu.regs.f = 0
		if ((sp ^ e8 ^ result) & 0x10) != 0 {
			cpu.regs.f |= FlagH
		}
		if ((sp ^ e8 ^ result) & 0x100) != 0 {
			cpu.regs.f |= FlagC
		}
		cpu.regs.sp = uint16(result)
		cpu.next(2, 4)

	case 0xE9: // JP HL
		cpu.regs.pc = cpu.getHL()
		cpu.cycles++

	case 0xEA: // LD (a16), A
		*cpu.getAddr(cpu.getA16()) = cpu.regs.a
		cpu.next(3, 16)

	case 0xEE: // XOR A, n8
		cpu.xor8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0xF0: // LDH A, (a8)
		cpu.regs.a = *cpu.getAddr(0xFF00 + uint16(cpu.getN8()))
		cpu.next(2, 3)

	case 0xF1: // POP AF
		cpu.pop16(cpu.setAF)
		cpu.next(1, 3)

	case 0xF2: // LDH A, (C)
		addr := 0xFF00 + uint16(cpu.regs.c)
		cpu.regs.a = *cpu.getAddr(addr)
		cpu.next(1, 2)

	case 0xF3: // DI
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		// TODO: Implementar DI
		cpu.next(1, 1)

	case 0xF5: // PUSH AF
		cpu.push16(cpu.getAF())
		cpu.next(1, 4)

	case 0xF6: // OR A, n8
		cpu.or8(&cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0xF8: // LD HL, SP+e8
		e8 := int16(cpu.getE8())
		sp := int16(cpu.regs.sp)
		result := sp + e8
		cpu.setHL(uint16(result))
		cpu.regs.f = 0
		if ((sp ^ e8 ^ result) & 0x10) != 0 {
			cpu.regs.f |= FlagH
		}
		if ((sp ^ e8 ^ result) & 0x100) != 0 {
			cpu.regs.f |= FlagC
		}
		cpu.next(2, 3)

	case 0xF9: // LD SP, HL
		cpu.regs.sp = cpu.getHL()
		cpu.next(1, 2)

	case 0xFA: // LD A, (a16)
		cpu.regs.a = *cpu.getAddr(cpu.getA16())
		cpu.next(3, 4)

	case 0xFB: // EI
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		// TODO: Implementar EI
		cpu.next(1, 1)

	case 0xFE: // CP A, n8
		cpu.cp8(cpu.regs.a, cpu.getN8())
		cpu.next(2, 2)

	case 0xFF: // RST 38H
		cpu.rst16(0x0038)
		cpu.next(1, 4)

	default:
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		panic("Detenido")
	}
	return cycles - cpu.cycles
}
