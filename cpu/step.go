package cpu

import (
	"fmt"
	"log"
)

func (cpu *CPU) Step() int {
	opcode := cpu.bus.Read(cpu.pc)
	cpu.Trace(opcode)
	if cpu.halted {
		// TODO: agregar CheckInterrupts() para cambiar halted a false
		fmt.Println("CPU Halted")
		cpu.cycles++
		return 1
	}
	switch opcode {
	case 0x00: // NOP
		return cpu.update(1, 1)

	case 0x01: // LD BC, n16
		cpu.setBC(cpu.getN16())
		return cpu.update(3, 3)

	case 0x02: // LD (BC),A
		cpu.bus.Write(cpu.getBC(), cpu.a)
		return cpu.update(1, 2)

	case 0x03: // INC BC
		cpu.inc16(cpu.setBC, cpu.getBC())
		return cpu.update(1, 2)

	case 0x04: // INC B
		cpu.inc8(&cpu.b)
		return cpu.update(1, 1)

	case 0x05: // DEC B
		cpu.dec8(&cpu.b)
		return cpu.update(1, 1)

	case 0x06: // LD B,n8
		cpu.b = cpu.getN8()
		return cpu.update(2, 2)
	case 0x07: // RLCA
		cpu.rlca()
		return cpu.update(1, 1)

	case 0x08: // LD (a16),SP
		cpu.setAddr16(cpu.getA16(), cpu.sp)
		return cpu.update(3, 5)

	case 0x09: // ADD HL,BC
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.getBC())
		return cpu.update(1, 2)

	case 0x0A: // LD A,(BC)
		cpu.a = cpu.bus.Read(cpu.getBC())
		return cpu.update(1, 2)

	case 0x0B: // DEC BC
		cpu.dec16(cpu.setBC, cpu.getBC())
		return cpu.update(1, 2)

	case 0x0C: // INC C
		cpu.inc8(&cpu.c)
		return cpu.update(1, 1)

	case 0x0D: // DEC C
		cpu.dec8(&cpu.c)
		return cpu.update(1, 1)

	case 0x0E: // LD C,n8
		cpu.c = cpu.getN8()
		return cpu.update(2, 2)

	case 0x0F: // RRCA
		cpu.rrca()
		return cpu.update(1, 1)
	case 0x10: // STOP
		// STOP 0 instruction (detiene el reloj del sistema)
		// El siguiente byte debe ser 0x00, pero normalmente se ignora
		// TODO: Agregar la lógica del modo STOP reloj/divider
		cpu.Stopped = true
		//panic("STOP todavía no está implementado")
		return cpu.update(2, 1)

	case 0x11: // LD DE, n16
		cpu.ld16(cpu.setDE, cpu.getN16())
		return cpu.update(3, 3)

	case 0x12: // LD (DE), A
		cpu.bus.Write(cpu.getDE(), cpu.a)
		return cpu.update(1, 2)

	case 0x13: // INC DE
		cpu.inc16(cpu.setDE, cpu.getDE())
		return cpu.update(1, 2)

	case 0x14: // INC D
		cpu.inc8(&cpu.d)
		return cpu.update(1, 1)

	case 0x15: // DEC D
		cpu.dec8(&cpu.d)
		return cpu.update(1, 1)

	case 0x16: // LD D,n8
		cpu.d = cpu.getN8()
		return cpu.update(2, 2)

	case 0x17: // RLA
		cpu.rla()
		return cpu.update(1, 1)

	case 0x18: // JR e8
		offset := cpu.getE8()
		cpu.pc = uint16(int(cpu.pc) + 2 + int(offset))
		cpu.cycles += 3

	case 0x19: // ADD HL,DE
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.getDE())
		return cpu.update(1, 2)

	case 0x1A: // LD A,(DE)
		cpu.a = cpu.bus.Read(cpu.getDE())
		return cpu.update(1, 2)

	case 0x1B: // DEC DE
		cpu.dec16(cpu.setDE, cpu.getDE())
		return cpu.update(1, 2)

	case 0x1C: // INC E
		cpu.inc8(&cpu.e)
		return cpu.update(1, 1)

	case 0x1D: // DEC E
		cpu.dec8(&cpu.e)
		return cpu.update(1, 1)

	case 0x1E: // LD E,n8
		cpu.e = cpu.getN8()
		return cpu.update(2, 2)

	case 0x1F: // RRA
		cpu.rra()
		return cpu.update(1, 1)
	case 0x20: // JR NZ,e8
		offset := cpu.getE8()
		if cpu.f&FlagZ == 0 {
			oldPC := cpu.pc
			cpu.pc += 2                                    // avanzar el PC antes del salto
			cpu.pc = uint16(int32(cpu.pc) + int32(offset)) // salto relativo con signo
			log.Printf("JR NZ,e8 saltando desde %04X a %04X con offset %d", oldPC, cpu.pc, offset)
			cpu.cycles += 3
			return 3
		}
		return cpu.update(2, 2) // no salta, solo avanza

	case 0x21: // LD HL, n16
		cpu.ld16(cpu.setHL, cpu.getN16())
		return cpu.update(3, 3)

	case 0x22: // LD (HL+),A
		hl := cpu.getHL()
		cpu.bus.Write(cpu.getHL(), cpu.a)
		cpu.setHL(hl + 1)
		return cpu.update(1, 2)

	case 0x23: // INC HL
		cpu.inc16(cpu.setHL, cpu.getHL())
		return cpu.update(1, 2)

	case 0x24: // INC H
		cpu.inc8(&cpu.h)
		return cpu.update(1, 1)

	case 0x25: // DEC H
		cpu.dec8(&cpu.h)
		return cpu.update(1, 1)

	case 0x26: // LD H,n8
		cpu.h = cpu.getN8()
		return cpu.update(2, 2)

	case 0x27: // DAA
		cpu.daa()
		return cpu.update(1, 1)

	case 0x28: // JR Z,e8
		offset := cpu.getE8()
		if cpu.f&FlagZ != 0 {
			cpu.pc = uint16(int(cpu.pc) + 2 + int(offset))
			cpu.cycles += 3
		} else {
			return cpu.update(2, 2)
		}

	case 0x29: // ADD HL,HL
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.getHL())
		return cpu.update(1, 2)

	case 0x2A: // LD A,(HL+)
		hl := cpu.getHL()
		cpu.a = cpu.bus.Read(cpu.getHL())
		cpu.setHL(hl + 1)
		return cpu.update(1, 2)

	case 0x2B: // DEC HL
		cpu.dec16(cpu.setHL, cpu.getHL())
		return cpu.update(1, 2)

	case 0x2C: // INC L
		cpu.inc8(&cpu.l)
		return cpu.update(1, 1)

	case 0x2D: // DEC L
		cpu.dec8(&cpu.l)
		return cpu.update(1, 1)

	case 0x2E: // LD L,n8
		cpu.l = cpu.getN8()
		return cpu.update(2, 2)

	case 0x2F: // CPL (Complement A)
		cpu.a = ^cpu.a
		cpu.f |= FlagN | FlagH
		return cpu.update(1, 1)
	case 0x30: // JR NC, e8
		offset := cpu.getE8()
		if cpu.f&FlagC == 0 {
			cpu.pc = uint16(int(cpu.pc) + 2 + int(offset))
			cpu.cycles += 1
		} else {
			return cpu.update(2, 1)
		}

	case 0x31: // LD SP, n16
		cpu.sp = cpu.getN16()
		return cpu.update(3, 3)

	case 0x32: // LD (HL-), A
		addr := cpu.getHL()
		log.Printf("LD (HL-), A: HL=%04X, A=%02X (escribiendo en %04X)", addr, cpu.a, addr)
		cpu.bus.Write(addr, cpu.a)
		cpu.setHL(addr - 1)
		return cpu.update(1, 2)

	case 0x33: // INC SP
		cpu.sp++
		return cpu.update(2, 2)

	case 0x34: // INC (HL)
		cpu.inc8Address(cpu.getHL())
		return cpu.update(3, 3)

	case 0x35: // DEC (HL)
		cpu.dec8Address(cpu.getHL())
		return cpu.update(3, 3)

	case 0x36: // LD (HL), n8
		cpu.bus.Write(cpu.getHL(), cpu.getN8())
		return cpu.update(3, 3)

	case 0x37: // SCF (Set Carry Flag)
		cpu.scf()
		return cpu.update(1, 1)

	case 0x38: // JR C, e8
		offset := cpu.getE8()
		if cpu.f&FlagC != 0 {
			cpu.pc = uint16(int(cpu.pc) + 2 + int(offset))
			cpu.cycles += 3
		} else {
			return cpu.update(2, 2)
		}

	case 0x39: // ADD HL, SP
		cpu.add16(cpu.setHL, cpu.getHL(), cpu.sp)
		return cpu.update(2, 2)

	case 0x3A: // LD A, (HL-)
		hl := cpu.getHL()
		cpu.a = cpu.bus.Read(hl)
		cpu.dec16(cpu.setHL, hl)
		return cpu.update(2, 2)

	case 0x3B: // DEC SP
		cpu.sp--
		return cpu.update(2, 2)

	case 0x3C: // INC A
		cpu.inc8(&cpu.a)
		return cpu.update(1, 1)

	case 0x3D: // DEC A
		cpu.dec8(&cpu.a)
		return cpu.update(1, 1)

	case 0x3E: // LD A, n8
		cpu.a = cpu.getN8()
		return cpu.update(2, 2)

	case 0x3F: // CCF
		cpu.ccf()
		return cpu.update(1, 1)
	case 0x40: // LD B,B
		// redundancia
		// cpu.b = cpu.b
		return cpu.update(1, 1)
	case 0x41: // LD B,C
		cpu.b = cpu.c
		return cpu.update(1, 1)
	case 0x42: // LD B,D
		cpu.b = cpu.d
		return cpu.update(1, 1)
	case 0x43: // LD B,E
		cpu.b = cpu.e
		return cpu.update(1, 1)
	case 0x44: // LD B,H
		cpu.b = cpu.h
		return cpu.update(1, 1)
	case 0x45: // LD B,L
		cpu.b = cpu.l
		return cpu.update(1, 1)
	case 0x46: // LD B,(HL)
		cpu.b = cpu.bus.Read(cpu.getHL())
		return cpu.update(1, 2)
	case 0x47: // LD B,A
		cpu.b = cpu.a
		return cpu.update(1, 1)

	case 0x48: // LD C,B
		cpu.c = cpu.b
		return cpu.update(1, 1)
	case 0x49: // LD C,C
		// redundancia
		// cpu.c = cpu.c
		return cpu.update(1, 1)
	case 0x4A: // LD C,D
		cpu.c = cpu.d
		return cpu.update(1, 1)
	case 0x4B: // LD C,E
		cpu.c = cpu.e
		return cpu.update(1, 1)
	case 0x4C: // LD C,H
		cpu.c = cpu.h
		return cpu.update(1, 1)
	case 0x4D: // LD C,L
		cpu.c = cpu.l
		return cpu.update(1, 1)
	case 0x4E: // LD C,(HL)
		cpu.c = cpu.bus.Read(cpu.getHL())
		return cpu.update(1, 2)
	case 0x4F: // LD C,A
		cpu.c = cpu.a
		return cpu.update(1, 1)
	case 0x50: // LD D,B
		cpu.d = cpu.b
		return cpu.update(1, 1)
	case 0x51: // LD D,C
		cpu.d = cpu.c
		return cpu.update(1, 1)
	case 0x52: // LD D,D
		//redundancia
		//cpu.d = cpu.d
		return cpu.update(1, 1)
	case 0x53: // LD D,E
		cpu.d = cpu.e
		return cpu.update(1, 1)
	case 0x54: // LD D,H
		cpu.d = cpu.h
		return cpu.update(1, 1)
	case 0x55: // LD D,L
		cpu.d = cpu.l
		return cpu.update(1, 1)
	case 0x56: // LD D,(HL)
		cpu.d = cpu.bus.Read(cpu.getHL())
		return cpu.update(1, 2)
	case 0x57: // LD D,A
		cpu.d = cpu.a
		return cpu.update(1, 1)

	case 0x58: // LD E,B
		cpu.e = cpu.a
		return cpu.update(1, 1)
	case 0x59: // LD E,C
		cpu.e = cpu.a
		return cpu.update(1, 1)
	case 0x5A: // LD E,D
		cpu.e = cpu.d
		return cpu.update(1, 1)
	case 0x5B: // LD E,E
		//redundancia
		//cpu.e = cpu.e
		return cpu.update(1, 1)
	case 0x5C: // LD E,H
		cpu.e = cpu.h
		return cpu.update(1, 1)
	case 0x5D: // LD E,L
		cpu.e = cpu.l
		return cpu.update(1, 1)
	case 0x5E: // LD E,(HL)
		cpu.e = cpu.bus.Read(cpu.getHL())
		return cpu.update(1, 2)
	case 0x5F: // LD E,A
		cpu.e = cpu.a
		return cpu.update(1, 1)
	case 0x60: // LD H,B
		cpu.h = cpu.b
		return cpu.update(1, 1)
	case 0x61: // LD H,C
		cpu.h = cpu.c
		return cpu.update(1, 1)
	case 0x62: // LD H,D
		cpu.h = cpu.d
		return cpu.update(1, 1)
	case 0x63: // LD H,E
		cpu.h = cpu.e
		return cpu.update(1, 1)
	case 0x64: // LD H,H
		//redundancia
		//cpu.h = cpu.h
		return cpu.update(1, 1)
	case 0x65: // LD H,L
		cpu.h = cpu.l
		return cpu.update(1, 1)
	case 0x66: // LD H,(HL)
		cpu.h = cpu.bus.Read(cpu.getHL())
		return cpu.update(1, 2)
	case 0x67: // LD H,A
		cpu.h = cpu.a
		return cpu.update(1, 1)

	case 0x68: // LD L,B
		cpu.l = cpu.b
		return cpu.update(1, 1)
	case 0x69: // LD L,C
		cpu.l = cpu.c
		return cpu.update(1, 1)
	case 0x6A: // LD L,D
		cpu.l = cpu.d
		return cpu.update(1, 1)
	case 0x6B: // LD L,E
		cpu.l = cpu.e
		return cpu.update(1, 1)
	case 0x6C: // LD L,H
		cpu.l = cpu.h
		return cpu.update(1, 1)
	case 0x6D: // LD L,L
		//redundancia
		//cpu.l = cpu.l
		return cpu.update(1, 1)
	case 0x6E: // LD L,(HL)
		cpu.l = cpu.bus.Read(cpu.getHL())
		return cpu.update(1, 2)
	case 0x6F: // LD L,A
		cpu.l = cpu.a
		return cpu.update(1, 1)

	case 0x70: // LD (HL),B
		cpu.bus.Write(cpu.getHL(), cpu.b)
		return cpu.update(1, 2)
	case 0x71: // LD (HL),C
		cpu.bus.Write(cpu.getHL(), cpu.c)
		return cpu.update(1, 2)
	case 0x72: // LD (HL),D
		cpu.bus.Write(cpu.getHL(), cpu.d)
		return cpu.update(1, 2)
	case 0x73: // LD (HL),E
		cpu.bus.Write(cpu.getHL(), cpu.e)
		return cpu.update(1, 2)
	case 0x74: // LD (HL),H
		cpu.bus.Write(cpu.getHL(), cpu.h)
		return cpu.update(1, 2)
	case 0x75: // LD (HL),L
		cpu.bus.Write(cpu.getHL(), cpu.l)
		return cpu.update(1, 2)
	case 0x76: // HALT
		cpu.halt()
		return cpu.update(1, 1)
	case 0x77: // LD (HL),A
		cpu.bus.Write(cpu.getHL(), cpu.a)
		return cpu.update(1, 2)

	case 0x78: // LD A,B
		cpu.a = cpu.b
		return cpu.update(1, 1)
	case 0x79: // LD A,C
		cpu.a = cpu.c
		return cpu.update(1, 1)
	case 0x7A: // LD A,D
		cpu.a = cpu.d
		return cpu.update(1, 1)
	case 0x7B: // LD A,E
		cpu.a = cpu.e
		return cpu.update(1, 1)
	case 0x7C: // LD A,H
		cpu.a = cpu.h
		return cpu.update(1, 1)
	case 0x7D: // LD A,L
		cpu.a = cpu.l
		return cpu.update(1, 1)
	case 0x7E: // LD A,(HL)
		cpu.a = cpu.bus.Read(cpu.getHL())
		return cpu.update(1, 2)
	case 0x7F: // LD A,A
		//redundancia
		//cpu.a = cpu.a
		return cpu.update(1, 1)
	case 0x80: // ADD A,B
		cpu.a = cpu.b
		return cpu.update(1, 1)
	case 0x81: // ADD A,C
		cpu.add8(&cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0x82: // ADD A,D
		cpu.add8(&cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0x83: // ADD A,E
		cpu.add8(&cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0x84: // ADD A,H
		cpu.add8(&cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0x85: // ADD A,L
		cpu.add8(&cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0x86: // ADD A,(HL)
		cpu.add8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0x87: // ADD A,A
		cpu.add8(&cpu.a, cpu.a)
		return cpu.update(1, 1)

	case 0x88: // ADC A,B
		cpu.adc8(&cpu.a, cpu.b)
		return cpu.update(1, 1)
	case 0x89: // ADC A,C
		cpu.adc8(&cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0x8A: // ADC A,D
		cpu.adc8(&cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0x8B: // ADC A,E
		cpu.adc8(&cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0x8C: // ADC A,H
		cpu.adc8(&cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0x8D: // ADC A,L
		cpu.adc8(&cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0x8E: // ADC A,(HL)
		cpu.adc8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0x8F: // ADC A,A
		cpu.adc8(&cpu.a, cpu.a)
		return cpu.update(1, 1)

	case 0x90: // SUB A, B
		cpu.sub8(&cpu.a, cpu.b)
		return cpu.update(1, 1)
	case 0x91: // SUB A, C
		cpu.sub8(&cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0x92: // SUB A, D
		cpu.sub8(&cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0x93: // SUB A, E
		cpu.sub8(&cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0x94: // SUB A, H
		cpu.sub8(&cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0x95: // SUB A, L
		cpu.sub8(&cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0x96: // SUB A, (HL)
		cpu.sub8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0x97: // SUB A, A
		cpu.sub8(&cpu.a, cpu.a)
		return cpu.update(1, 1)

	case 0x98: // SBC A,B
		cpu.sbc8(&cpu.a, cpu.b)
		return cpu.update(1, 1)
	case 0x99: // SBC A,C
		cpu.sbc8(&cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0x9A: // SBC A,D
		cpu.sbc8(&cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0x9B: // SBC A,E
		cpu.sbc8(&cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0x9C: // SBC A,H
		cpu.sbc8(&cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0x9D: // SBC A,L
		cpu.sbc8(&cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0x9E: // SBC A,(HL)
		cpu.sbc8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0x9F: // SBC A,A
		cpu.sbc8(&cpu.a, cpu.a)
		return cpu.update(1, 1)

	case 0xA0: // AND B
		cpu.and8(&cpu.a, cpu.b)
		return cpu.update(1, 1)
	case 0xA1: // AND C
		cpu.and8(&cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0xA2: // AND D
		cpu.and8(&cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0xA3: // AND E
		cpu.and8(&cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0xA4: // AND H
		cpu.and8(&cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0xA5: // AND L
		cpu.and8(&cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0xA6: // AND (HL)
		cpu.and8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0xA7: // AND A
		cpu.and8(&cpu.a, cpu.a)
		return cpu.update(1, 1)

	case 0xA8: // XOR B
		cpu.xor8(&cpu.a, cpu.b)
		return cpu.update(1, 1)
	case 0xA9: // XOR C
		cpu.xor8(&cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0xAA: // XOR D
		cpu.xor8(&cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0xAB: // XOR E
		cpu.xor8(&cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0xAC: // XOR H
		cpu.xor8(&cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0xAD: // XOR L
		cpu.xor8(&cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0xAE: // XOR (HL)
		cpu.xor8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0xAF: // XOR A, A
		cpu.a ^= cpu.a
		cpu.f = 0
		cpu.f |= FlagZ
		return cpu.update(1, 1)

	case 0xB0: // OR B
		cpu.or8(&cpu.a, cpu.b)
		return cpu.update(1, 1)
	case 0xB1: // OR C
		cpu.or8(&cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0xB2: // OR D
		cpu.or8(&cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0xB3: // OR E
		cpu.or8(&cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0xB4: // OR H
		cpu.or8(&cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0xB5: // OR L
		cpu.or8(&cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0xB6: // OR (HL)
		cpu.or8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0xB7: // OR A
		cpu.or8(&cpu.a, cpu.a)
		return cpu.update(1, 1)

	case 0xB8: // CP B
		cpu.cp8(cpu.a, cpu.b)
		return cpu.update(1, 1)
	case 0xB9: // CP C
		cpu.cp8(cpu.a, cpu.c)
		return cpu.update(1, 1)
	case 0xBA: // CP D
		cpu.cp8(cpu.a, cpu.d)
		return cpu.update(1, 1)
	case 0xBB: // CP E
		cpu.cp8(cpu.a, cpu.e)
		return cpu.update(1, 1)
	case 0xBC: // CP H
		cpu.cp8(cpu.a, cpu.h)
		return cpu.update(1, 1)
	case 0xBD: // CP L
		cpu.cp8(cpu.a, cpu.l)
		return cpu.update(1, 1)
	case 0xBE: // CP (HL)
		cpu.cp8(cpu.a, cpu.bus.Read(cpu.getHL()))
		return cpu.update(1, 2)
	case 0xBF: // CP A
		cpu.cp8(cpu.a, cpu.a)
		return cpu.update(1, 1)
	case 0xC0: // RET NZ
		if cpu.f&FlagZ == 0 {
			cpu.ret()
			cpu.cycles += 20
		} else {
			return cpu.update(1, 8)
		}

	case 0xC1: // POP BC
		cpu.pop16(cpu.setBC)
		return cpu.update(1, 12)

	case 0xC2: // JP NZ, a16
		if cpu.f&FlagZ == 0 {
			cpu.pc = cpu.getA16()
			cpu.cycles += 4
			return 4
		} else {
			return cpu.update(3, 3)
		}

	case 0xC3: // JP a16
		cpu.pc = cpu.getA16()
		cpu.cycles += 4
		return 4

	case 0xC4: // CALL NZ, a16
		if cpu.f&FlagZ == 0 {
			cpu.call16(cpu.getA16())
			cpu.cycles += 4
			return 4
			//return cpu.update(6, 4)
		} else {
			return cpu.update(3, 3)
		}

	case 0xC5: // PUSH BC
		cpu.push16(cpu.getBC())
		return cpu.update(1, 4)

	case 0xC6: // ADD A, n8
		cpu.add8(&cpu.a, cpu.getN8())
		return cpu.update(2, 2)

	case 0xC7: // RST 00H
		cpu.rst16(0x00)

	case 0xC8: // RET Z
		if cpu.f&FlagZ != 0 {
			cpu.ret()
			return cpu.update(1, 5)
		} else {
			return cpu.update(1, 2)
		}

	case 0xC9: // RET
		cpu.ret()
		return cpu.update(1, 4)

	case 0xCA: // JP Z, a16
		if cpu.f&FlagZ != 0 {
			cpu.pc = cpu.getA16()
			cpu.cycles += 4
		} else {
			return cpu.update(3, 3)
		}

	case 0xCB: // PREFIX CB

		// TODO: Implementar Prefix $CB
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		panic("Detenido")
		//return cpu.update(1, 4)

	case 0xCC: // CALL Z, a16
		if cpu.f&FlagZ != 0 {
			cpu.call16(cpu.getA16())
			cpu.cycles += 6
			return 6
			//return cpu.update(3, 6)
		} else {
			return cpu.update(3, 3)
		}

	case 0xCD: // CALL a16
		cpu.call16(cpu.getA16())
		cpu.cycles += 6
		return 6
		//return cpu.update(3, 6)

	case 0xCE: // ADC A, n8
		cpu.adc8(&cpu.a, cpu.getN8())
		return cpu.update(2, 8)

	case 0xCF: // RST 08H
		cpu.rst16(0x08)
		return cpu.update(1, 4)
	case 0xD0: // RET NC
		if cpu.f&FlagC == 0 {
			cpu.ret()
			return cpu.update(1, 5)
		} else {
			return cpu.update(1, 2)
		}

	case 0xD1: // POP DE
		cpu.pop16(cpu.setDE)
		return cpu.update(1, 3)

	case 0xD2: // JP NC, a16
		if cpu.f&FlagC == 0 {
			cpu.pc = cpu.getA16()
			cpu.cycles += 4
		} else {
			return cpu.update(3, 3)
		}

	case 0xD4: // CALL NC, a16
		if cpu.f&FlagC == 0 {
			cpu.call16(cpu.getA16())
			cpu.cycles += 6
			return 6
			//return cpu.update(3, 6)
		} else {
			return cpu.update(3, 3)
		}

	case 0xD5: // PUSH DE
		cpu.push16(cpu.getDE())
		return cpu.update(1, 4)

	case 0xD6: // SUB n8
		cpu.sub8(&cpu.a, cpu.getN8())
		return cpu.update(2, 2)

	case 0xD7: // RST 10H
		cpu.rst16(0x10)
		return cpu.update(1, 4)

	case 0xD8: // RET C
		if cpu.f&FlagC != 0 {
			cpu.ret()
			return cpu.update(1, 5)
		} else {
			return cpu.update(1, 2)
		}

	case 0xD9: // RETI
		cpu.reti()
		//return cpu.update(1, 4)

	case 0xDA: // JP C, a16
		if cpu.f&FlagC != 0 {
			cpu.pc = cpu.getA16()
			cpu.cycles += 4
		} else {
			return cpu.update(3, 3)
		}

	case 0xDC: // CALL C, a16
		if cpu.f&FlagC != 0 {
			cpu.call16(cpu.getA16())
			cpu.cycles += 6
			return 6
			//return cpu.update(1, 6)
		} else {
			return cpu.update(3, 3)
		}

	case 0xDE: // SBC A, n8
		cpu.sbc8(&cpu.a, cpu.getN8())
		return cpu.update(2, 2)

	case 0xDF: // RST 18H
		cpu.rst16(0x18)
		return cpu.update(1, 4)
	case 0xE0: // LDH (a8), A
		cpu.bus.Write(cpu.getA8(), cpu.a)
		return cpu.update(2, 3)

	case 0xE1: // POP HL
		cpu.pop16(cpu.setHL)
		return cpu.update(1, 3)

	case 0xE2: // LDH (C), A
		cpu.bus.Write(0xFF00+uint16(cpu.c), cpu.a)
		return cpu.update(1, 2)

	case 0xE5: // PUSH HL
		cpu.push16(cpu.getHL())
		return cpu.update(1, 16)

	case 0xE6: // AND A,n8
		cpu.add8(&cpu.a, cpu.getN8())
		return cpu.update(2, 2)

	case 0xE8: // ADD SP, e8
		e8 := int16(cpu.getE8())
		sp := int16(cpu.sp)
		result := sp + e8
		cpu.f = 0
		if ((sp ^ e8 ^ result) & 0x10) != 0 {
			cpu.f |= FlagH
		}
		if ((sp ^ e8 ^ result) & 0x100) != 0 {
			cpu.f |= FlagC
		}
		cpu.sp = uint16(result)
		return cpu.update(2, 4)

	case 0xE9: // JP HL
		cpu.pc = cpu.getHL()
		cpu.cycles++

	case 0xEA: // LD (a16), A
		cpu.bus.Write(cpu.getA16(), cpu.a)
		return cpu.update(3, 16)

	case 0xEE: // XOR A, n8
		cpu.xor8(&cpu.a, cpu.getN8())
		return cpu.update(2, 2)

	case 0xF0: // LDH A, (a8)
		cpu.a = cpu.bus.Read(cpu.getA8())
		return cpu.update(2, 3)

	case 0xF1: // POP AF
		cpu.pop16(cpu.setAF)
		return cpu.update(1, 3)

	case 0xF2: // LDH A, (C)
		cpu.a = cpu.bus.Read(0xFF00 + uint16(cpu.c))
		return cpu.update(1, 2)

	case 0xF3: // DI
		cpu.di()
		return cpu.update(1, 1)

	case 0xF5: // PUSH AF
		cpu.push16(cpu.getAF())
		return cpu.update(1, 4)

	case 0xF6: // OR A, n8
		cpu.or8(&cpu.a, cpu.getN8())
		return cpu.update(2, 2)

	case 0xF8: // LD HL, SP+e8
		cpu.ld_HL_SP_e8()
		return cpu.update(2, 3)

	case 0xF9: // LD SP, HL
		cpu.sp = cpu.getHL()
		return cpu.update(1, 2)

	case 0xFA: // LD A, (a16)
		cpu.a = cpu.bus.Read(cpu.getA16())
		return cpu.update(3, 4)

	case 0xFB: // EI
		// The flag is only set after the instruction following EI.
		cpu.ei()
		cpu.pc++
		return 1
		//return cpu.update(1, 1)

	case 0xFE: // CP A, n8
		cpu.cp8(cpu.a, cpu.getN8())
		return cpu.update(2, 2)

	case 0xFF: // RST 38H
		cpu.rst16(0x0038)
		cpu.cycles += 4
		return 4
		//return cpu.update(1, 4)

	default:
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		panic("Detenido")
	}
	return 0
}
