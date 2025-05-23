package cpu

import "fmt"

func (cpu *CPU) execute(opcode byte) {
	// Decode & Execute
	switch opcode {
	case 0x00: // NOP
		return

	case 0x01: // LD BC, n16
		cpu.ldBC(cpu.getN16())
		return

	case 0x02: // LD (BC),A
		cpu.Write(cpu.getBC(), cpu.a)
		return

	case 0x03: // INC BC
		cpu.inc16(cpu.ldBC, cpu.getBC())
		return

	case 0x04: // INC B
		cpu.incR(&cpu.b)
		return

	case 0x05: // DEC B
		cpu.decR(&cpu.b)
		return

	case 0x06: // LD B, n8
		cpu.b = cpu.getN8()
		cpu.tick()
		return
	case 0x07: // RLCA
		cpu.rlca()
		return

	case 0x08: // LD (a16),SP
		cpu.setAddr16(cpu.getA16(), cpu.sp)
		return

	case 0x09: // ADD HL,BC
		cpu.addHL(cpu.getBC())
		return

	case 0x0A: // LD A,(BC)
		cpu.a = cpu.bus.Read(cpu.getBC())
		cpu.tick()
		return

	case 0x0B: // DEC BC
		cpu.dec16(cpu.ldBC, cpu.getBC())
		return

	case 0x0C: // INC C
		cpu.incR(&cpu.c)
		return

	case 0x0D: // DEC C
		cpu.decR(&cpu.c)
		return

	case 0x0E: // LD C, n8
		cpu.c = cpu.getN8()
		cpu.tick()
		return

	case 0x0F: // RRCA
		cpu.rrca()
		return
	case 0x10: // STOP
		// STOP 0 instruction (detiene el reloj del sistema)
		// El siguiente byte debe ser 0x00, pero normalmente se ignora
		// TODO: Averiguar si stop debe saltar un byte en PC?
		cpu.bus.Write(DIVRegister, 0x00)
		cpu.Stopped = true
		return

	case 0x11: // LD DE, n16
		cpu.ldDE(cpu.getN16())
		return

	case 0x12: // LD (DE), A
		cpu.Write(cpu.getDE(), cpu.a)
		return

	case 0x13: // INC DE
		cpu.inc16(cpu.ldDE, cpu.getDE())
		return

	case 0x14: // INC D
		cpu.incR(&cpu.d)
		return

	case 0x15: // DEC D
		cpu.decR(&cpu.d)
		return

	case 0x16: // LD D, n8
		cpu.d = cpu.getN8()
		cpu.tick()
		return

	case 0x17: // RLA
		cpu.rla()
		return

	case 0x18: // JR e8
		offset := cpu.getE8()
		cpu.tick()
		cpu.pc = uint16(int32(cpu.pc) + int32(offset))
		cpu.tick()

	case 0x19: // ADD HL,DE
		cpu.addHL(cpu.getDE())
		return

	case 0x1A: // LD A,(DE)
		cpu.a = cpu.bus.Read(cpu.getDE())
		cpu.tick()
		return

	case 0x1B: // DEC DE
		cpu.dec16(cpu.ldDE, cpu.getDE())
		return

	case 0x1C: // INC E
		cpu.incR(&cpu.e)
		return

	case 0x1D: // DEC E
		cpu.decR(&cpu.e)
		return

	case 0x1E: // LD E, n8
		cpu.e = cpu.getN8()
		cpu.tick()
		return

	case 0x1F: // RRA
		cpu.rra()
		return

	case 0x20: // JR NZ,e8
		offset := cpu.getE8()
		cpu.tick()
		if cpu.f&FlagZ == 0 {
			cpu.pc = uint16(int16(cpu.pc) + int16(offset))
			cpu.tick()
			return
		}
		return

	case 0x21: // LD HL, n16
		cpu.ldHL(cpu.getN16())
		return

	case 0x22: // LD (HL+),A
		hl := cpu.getHL()
		cpu.bus.Write(hl, cpu.a)
		cpu.ldHL(hl + 1)
		cpu.tick()
		return

	case 0x23: // INC HL
		cpu.inc16(cpu.ldHL, cpu.getHL())
		return

	case 0x24: // INC H
		cpu.incR(&cpu.h)
		return

	case 0x25: // DEC H
		cpu.decR(&cpu.h)
		return

	case 0x26: // LD H, n8
		cpu.h = cpu.getN8()
		cpu.tick()
		return

	case 0x27: // DAA
		cpu.daa()
		return

	case 0x28: // JR Z,e8
		offset := cpu.getE8()
		cpu.tick()
		if cpu.f&FlagZ != 0 {
			cpu.pc = uint16(int(cpu.pc) + int(offset))
			cpu.tick()
			return
		}
		return

	case 0x29: // ADD HL,HL
		cpu.addHL(cpu.getHL())
		return

	case 0x2A: // LD A,(HL+)
		hl := cpu.getHL()
		cpu.a = cpu.bus.Read(cpu.getHL())
		cpu.ldHL(hl + 1)
		cpu.tick()
		return

	case 0x2B: // DEC HL
		cpu.dec16(cpu.ldHL, cpu.getHL())
		return

	case 0x2C: // INC L
		cpu.incR(&cpu.l)
		return

	case 0x2D: // DEC L
		cpu.decR(&cpu.l)
		return

	case 0x2E: // LD L, n8
		cpu.l = cpu.getN8()
		cpu.tick()
		return

	case 0x2F: // CPL (Complement A)
		cpu.cpl()
		return

	case 0x30: // JR NC, e8
		offset := cpu.getE8()
		cpu.tick()
		if cpu.f&FlagC == 0 {
			cpu.pc = uint16(int16(cpu.pc) + int16(offset))
			cpu.tick()
			return
		}
		return

	case 0x31: // LD SP, n16
		cpu.sp = cpu.getN16()
		return

	case 0x32: // LD (HL-), A
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.a)
		cpu.ldHL(addr - 1)
		cpu.tick()
		return

	case 0x33: // INC SP
		cpu.sp++
		cpu.tick()
		return

	case 0x34: // INC (HL)
		cpu.incHL()
		return

	case 0x35: // DEC (HL)
		cpu.decHL()
		return

	case 0x36: // LD (HL), n8
		cpu.tick()
		cpu.Write(cpu.getHL(), cpu.getN8())
		return

	case 0x37: // SCF (Set Carry Flag)
		cpu.scf()
		return

	case 0x38: // JR C, e8
		offset := cpu.getE8()
		cpu.tick()
		if cpu.f&FlagC != 0 {
			cpu.pc = uint16(int16(cpu.pc) + int16(offset))
			cpu.tick()
			return
		}
		return

	case 0x39: // ADD HL, SP
		cpu.addHL(cpu.sp)
		return

	case 0x3A: // LD A, (HL-)
		hl := cpu.getHL()
		cpu.a = cpu.bus.Read(hl)
		cpu.ldHL(hl - 1)
		cpu.tick()
		return

	case 0x3B: // DEC SP
		cpu.sp--
		cpu.tick()
		return

	case 0x3C: // INC A
		cpu.incR(&cpu.a)
		return

	case 0x3D: // DEC A
		cpu.decR(&cpu.a)
		return

	case 0x3E: // LD A, n8
		cpu.a = cpu.getN8()
		cpu.tick()
		return

	case 0x3F: // CCF
		cpu.ccf()
		return
	case 0x40: // LD B,B
		// redundancia
		// cpu.b = cpu.b
		return
	case 0x41: // LD B,C
		cpu.b = cpu.c
		return
	case 0x42: // LD B,D
		cpu.b = cpu.d
		return
	case 0x43: // LD B,E
		cpu.b = cpu.e
		return
	case 0x44: // LD B,H
		cpu.b = cpu.h
		return
	case 0x45: // LD B,L
		cpu.b = cpu.l
		return
	case 0x46: // LD B,(HL)
		cpu.b = cpu.bus.Read(cpu.getHL())
		cpu.tick()
		return
	case 0x47: // LD B,A
		cpu.b = cpu.a
		return

	case 0x48: // LD C,B
		cpu.c = cpu.b
		return
	case 0x49: // LD C,C
		// redundancia
		// cpu.c = cpu.c
		return
	case 0x4A: // LD C,D
		cpu.c = cpu.d
		return
	case 0x4B: // LD C,E
		cpu.c = cpu.e
		return
	case 0x4C: // LD C,H
		cpu.c = cpu.h
		return
	case 0x4D: // LD C,L
		cpu.c = cpu.l
		return
	case 0x4E: // LD C,(HL)
		cpu.c = cpu.bus.Read(cpu.getHL())
		cpu.tick()
		return
	case 0x4F: // LD C,A
		cpu.c = cpu.a
		return
	case 0x50: // LD D,B
		cpu.d = cpu.b
		return
	case 0x51: // LD D,C
		cpu.d = cpu.c
		return
	case 0x52: // LD D,D
		//redundancia
		//cpu.d = cpu.d
		return
	case 0x53: // LD D,E
		cpu.d = cpu.e
		return
	case 0x54: // LD D,H
		cpu.d = cpu.h
		return
	case 0x55: // LD D,L
		cpu.d = cpu.l
		return
	case 0x56: // LD D,(HL)
		cpu.d = cpu.bus.Read(cpu.getHL())
		cpu.tick()
		return
	case 0x57: // LD D,A
		cpu.d = cpu.a
		return

	case 0x58: // LD E,B
		cpu.e = cpu.b
		return
	case 0x59: // LD E,C
		cpu.e = cpu.c
		return
	case 0x5A: // LD E,D
		cpu.e = cpu.d
		return
	case 0x5B: // LD E,E
		//redundancia
		//cpu.e = cpu.e
		return
	case 0x5C: // LD E,H
		cpu.e = cpu.h
		return
	case 0x5D: // LD E,L
		cpu.e = cpu.l
		return
	case 0x5E: // LD E,(HL)
		cpu.e = cpu.bus.Read(cpu.getHL())
		cpu.tick()
		return
	case 0x5F: // LD E,A
		cpu.e = cpu.a
		return
	case 0x60: // LD H,B
		cpu.h = cpu.b
		return
	case 0x61: // LD H,C
		cpu.h = cpu.c
		return
	case 0x62: // LD H,D
		cpu.h = cpu.d
		return
	case 0x63: // LD H,E
		cpu.h = cpu.e
		return
	case 0x64: // LD H,H
		//redundancia
		//cpu.h = cpu.h
		return
	case 0x65: // LD H,L
		cpu.h = cpu.l
		return
	case 0x66: // LD H,(HL)
		cpu.h = cpu.bus.Read(cpu.getHL())
		cpu.tick()
		return
	case 0x67: // LD H,A
		cpu.h = cpu.a
		return

	case 0x68: // LD L,B
		cpu.l = cpu.b
		return
	case 0x69: // LD L,C
		cpu.l = cpu.c
		return
	case 0x6A: // LD L,D
		cpu.l = cpu.d
		return
	case 0x6B: // LD L,E
		cpu.l = cpu.e
		return
	case 0x6C: // LD L,H
		cpu.l = cpu.h
		return
	case 0x6D: // LD L,L
		//redundancia
		//cpu.l = cpu.l
		return
	case 0x6E: // LD L,(HL)
		cpu.l = cpu.bus.Read(cpu.getHL())
		cpu.tick()
		return
	case 0x6F: // LD L,A
		cpu.l = cpu.a
		return

	case 0x70: // LD (HL),B
		cpu.Write(cpu.getHL(), cpu.b)
		return
	case 0x71: // LD (HL),C
		cpu.Write(cpu.getHL(), cpu.c)
		return
	case 0x72: // LD (HL),D
		cpu.Write(cpu.getHL(), cpu.d)
		return
	case 0x73: // LD (HL),E
		cpu.Write(cpu.getHL(), cpu.e)
		return
	case 0x74: // LD (HL),H
		cpu.Write(cpu.getHL(), cpu.h)
		return
	case 0x75: // LD (HL),L
		cpu.Write(cpu.getHL(), cpu.l)
		return
	case 0x76: // HALT
		cpu.halt()
		return
	case 0x77: // LD (HL),A
		cpu.Write(cpu.getHL(), cpu.a)
		return

	case 0x78: // LD A,B
		cpu.a = cpu.b
		return
	case 0x79: // LD A,C
		cpu.a = cpu.c
		return
	case 0x7A: // LD A,D
		cpu.a = cpu.d
		return
	case 0x7B: // LD A,E
		cpu.a = cpu.e
		return
	case 0x7C: // LD A,H
		cpu.a = cpu.h
		return
	case 0x7D: // LD A,L
		cpu.a = cpu.l
		return
	case 0x7E: // LD A,(HL)
		cpu.a = cpu.bus.Read(cpu.getHL())
		cpu.tick()
		return
	case 0x7F: // LD A,A
		//redundancia
		//cpu.a = cpu.a
		return
	case 0x80: // ADD A,B
		cpu.add8(&cpu.a, cpu.b)
		return
	case 0x81: // ADD A,C
		cpu.add8(&cpu.a, cpu.c)
		return
	case 0x82: // ADD A,D
		cpu.add8(&cpu.a, cpu.d)
		return
	case 0x83: // ADD A,E
		cpu.add8(&cpu.a, cpu.e)
		return
	case 0x84: // ADD A,H
		cpu.add8(&cpu.a, cpu.h)
		return
	case 0x85: // ADD A,L
		cpu.add8(&cpu.a, cpu.l)
		return
	case 0x86: // ADD A,(HL)
		cpu.add8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0x87: // ADD A,A
		cpu.add8(&cpu.a, cpu.a)
		return

	case 0x88: // ADC A,B
		cpu.adc8(&cpu.a, cpu.b)
		return
	case 0x89: // ADC A,C
		cpu.adc8(&cpu.a, cpu.c)
		return
	case 0x8A: // ADC A,D
		cpu.adc8(&cpu.a, cpu.d)
		return
	case 0x8B: // ADC A,E
		cpu.adc8(&cpu.a, cpu.e)
		return
	case 0x8C: // ADC A,H
		cpu.adc8(&cpu.a, cpu.h)
		return
	case 0x8D: // ADC A,L
		cpu.adc8(&cpu.a, cpu.l)
		return
	case 0x8E: // ADC A,(HL)
		cpu.adc8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0x8F: // ADC A,A
		cpu.adc8(&cpu.a, cpu.a)
		return

	case 0x90: // SUB A, B
		cpu.sub8(&cpu.a, cpu.b)
		return
	case 0x91: // SUB A, C
		cpu.sub8(&cpu.a, cpu.c)
		return
	case 0x92: // SUB A, D
		cpu.sub8(&cpu.a, cpu.d)
		return
	case 0x93: // SUB A, E
		cpu.sub8(&cpu.a, cpu.e)
		return
	case 0x94: // SUB A, H
		cpu.sub8(&cpu.a, cpu.h)
		return
	case 0x95: // SUB A, L
		cpu.sub8(&cpu.a, cpu.l)
		return
	case 0x96: // SUB A, (HL)
		cpu.sub8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0x97: // SUB A, A
		cpu.sub8(&cpu.a, cpu.a)
		return

	case 0x98: // SBC A,B
		cpu.sbc8(&cpu.a, cpu.b)
		return
	case 0x99: // SBC A,C
		cpu.sbc8(&cpu.a, cpu.c)
		return
	case 0x9A: // SBC A,D
		cpu.sbc8(&cpu.a, cpu.d)
		return
	case 0x9B: // SBC A,E
		cpu.sbc8(&cpu.a, cpu.e)
		return
	case 0x9C: // SBC A,H
		cpu.sbc8(&cpu.a, cpu.h)
		return
	case 0x9D: // SBC A,L
		cpu.sbc8(&cpu.a, cpu.l)
		return
	case 0x9E: // SBC A,(HL)
		cpu.sbc8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0x9F: // SBC A,A
		cpu.sbc8(&cpu.a, cpu.a)
		return

	case 0xA0: // AND B
		cpu.and8(&cpu.a, cpu.b)
		return
	case 0xA1: // AND C
		cpu.and8(&cpu.a, cpu.c)
		return
	case 0xA2: // AND D
		cpu.and8(&cpu.a, cpu.d)
		return
	case 0xA3: // AND E
		cpu.and8(&cpu.a, cpu.e)
		return
	case 0xA4: // AND H
		cpu.and8(&cpu.a, cpu.h)
		return
	case 0xA5: // AND L
		cpu.and8(&cpu.a, cpu.l)
		return
	case 0xA6: // AND (HL)
		cpu.and8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0xA7: // AND A
		cpu.and8(&cpu.a, cpu.a)
		return

	case 0xA8: // XOR B
		cpu.xor8(&cpu.a, cpu.b)
		return
	case 0xA9: // XOR C
		cpu.xor8(&cpu.a, cpu.c)
		return
	case 0xAA: // XOR D
		cpu.xor8(&cpu.a, cpu.d)
		return
	case 0xAB: // XOR E
		cpu.xor8(&cpu.a, cpu.e)
		return
	case 0xAC: // XOR H
		cpu.xor8(&cpu.a, cpu.h)
		return
	case 0xAD: // XOR L
		cpu.xor8(&cpu.a, cpu.l)
		return
	case 0xAE: // XOR (HL)
		cpu.xor8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0xAF: // XOR A, A
		cpu.xor8(&cpu.a, cpu.a)
		return

	case 0xB0: // OR B
		cpu.or8(&cpu.a, cpu.b)
		return
	case 0xB1: // OR C
		cpu.or8(&cpu.a, cpu.c)
		return
	case 0xB2: // OR D
		cpu.or8(&cpu.a, cpu.d)
		return
	case 0xB3: // OR E
		cpu.or8(&cpu.a, cpu.e)
		return
	case 0xB4: // OR H
		cpu.or8(&cpu.a, cpu.h)
		return
	case 0xB5: // OR L
		cpu.or8(&cpu.a, cpu.l)
		return
	case 0xB6: // OR (HL)
		cpu.or8(&cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0xB7: // OR A
		cpu.or8(&cpu.a, cpu.a)
		return

	case 0xB8: // CP B
		cpu.cp8(cpu.a, cpu.b)
		return
	case 0xB9: // CP C
		cpu.cp8(cpu.a, cpu.c)
		return
	case 0xBA: // CP D
		cpu.cp8(cpu.a, cpu.d)
		return
	case 0xBB: // CP E
		cpu.cp8(cpu.a, cpu.e)
		return
	case 0xBC: // CP H
		cpu.cp8(cpu.a, cpu.h)
		return
	case 0xBD: // CP L
		cpu.cp8(cpu.a, cpu.l)
		return
	case 0xBE: // CP (HL)
		cpu.cp8(cpu.a, cpu.bus.Read(cpu.getHL()))
		cpu.tick()
		return
	case 0xBF: // CP A
		cpu.cp8(cpu.a, cpu.a)
		return
	case 0xC0: // RET NZ
		cpu.tick()
		if cpu.f&FlagZ == 0 {
			cpu.ret()
			return
		}
		return

	case 0xC1: // POP BC
		cpu.pop16(cpu.ldBC)
		return

	case 0xC2: // JP NZ, a16
		addr := cpu.getN16()
		if cpu.f&FlagZ == 0 {
			cpu.pc = addr
			cpu.tick()
			return
		}
		return

	case 0xC3: // JP a16
		addr := cpu.getA16()
		cpu.pc = addr
		cpu.tick()
		return

	case 0xC4: // CALL NZ, a16
		addr := cpu.getA16()
		if cpu.f&FlagZ == 0 {
			cpu.call16(addr)
			return
		}
		return

	case 0xC5: // PUSH BC
		cpu.push16(cpu.getBC())
		return

	case 0xC6: // ADD A, n8
		cpu.add8(&cpu.a, cpu.getN8())
		cpu.tick()
		return

	case 0xC7: // RST 00H
		cpu.call16(0x00)
		return

	case 0xC8: // RET Z
		cpu.tick()
		if cpu.f&FlagZ != 0 {
			cpu.ret()
			return
		}
		return

	case 0xC9: // RET
		cpu.ret()
		return

	case 0xCA: // JP Z, a16
		addr := cpu.getA16()
		if cpu.f&FlagZ != 0 {
			cpu.pc = addr
			cpu.tick()
			return
		}
		return

	case 0xCB: // PREFIX CB
		cpu.executeCB()
		return

	case 0xCC: // CALL Z, a16
		addr := cpu.getA16()
		if cpu.f&FlagZ != 0 {
			cpu.call16(addr)
			return
		}
		return

	case 0xCD: // CALL a16
		cpu.call16(cpu.getA16())
		return

	case 0xCE: // ADC A, n8
		cpu.adc8(&cpu.a, cpu.getN8())
		cpu.tick()
		return

	case 0xCF: // RST 08H
		cpu.call16(0x08)
		return
	case 0xD0: // RET NC
		cpu.tick()
		if cpu.f&FlagC == 0 {
			cpu.ret()
			return
		}
		return

	case 0xD1: // POP DE
		cpu.pop16(cpu.ldDE)
		return

	case 0xD2: // JP NC, a16
		addr := cpu.getA16()
		if cpu.f&FlagC == 0 {
			cpu.pc = addr
			cpu.tick()
			return
		}
		return

	case 0xD4: // CALL NC, a16
		addr := cpu.getA16()
		if cpu.f&FlagC == 0 {
			cpu.call16(addr)
			return
		}
		return

	case 0xD5: // PUSH DE
		cpu.push16(cpu.getDE())
		return

	case 0xD6: // SUB n8
		cpu.sub8(&cpu.a, cpu.getN8())
		cpu.tick()
		return

	case 0xD7: // RST 10H
		cpu.call16(0x10)
		return

	case 0xD8: // RET C
		cpu.tick()
		if cpu.f&FlagC != 0 {
			cpu.ret()
			return
		}
		return

	case 0xD9: // RETI
		cpu.reti()
		return

	case 0xDA: // JP C, a16
		addr := cpu.getA16()
		if cpu.f&FlagC != 0 {
			cpu.pc = addr
			cpu.tick()
			return
		}
		return

	case 0xDC: // CALL C, a16
		addr := cpu.getA16()
		if cpu.f&FlagC != 0 {
			cpu.call16(addr)
			return
		}
		return

	case 0xDE: // SBC A, n8
		cpu.sbc8(&cpu.a, cpu.getN8())
		cpu.tick()
		return

	case 0xDF: // RST 18H
		cpu.call16(0x18)
		return
	case 0xE0: // LDH (a8), A
		cpu.Write(cpu.getA8(), cpu.a)
		return

	case 0xE1: // POP HL
		cpu.pop16(cpu.ldHL)
		return

	case 0xE2: // LDH (C), A
		cpu.Write(0xFF00+uint16(cpu.c), cpu.a)
		return

	case 0xE5: // PUSH HL
		cpu.push16(cpu.getHL())
		return

	case 0xE6: // AND A,n8
		cpu.and8(&cpu.a, cpu.getN8())
		cpu.tick()
		return
	case 0xE7: // RST 20H
		cpu.call16(0x20)
		return

	case 0xE8: // ADD SP, e8
		e8 := int16(cpu.getE8())
		cpu.tick()
		sp := int16(cpu.sp)
		result := sp + e8
		cpu.tick()
		cpu.f = 0
		if ((sp ^ e8 ^ result) & 0x10) != 0 {
			cpu.f |= FlagH
		}
		if ((sp ^ e8 ^ result) & 0x100) != 0 {
			cpu.f |= FlagC
		}
		cpu.sp = uint16(result)
		cpu.tick()
		return

	case 0xE9: // JP HL
		cpu.pc = cpu.getHL()
		return

	case 0xEA: // LD (a16), A
		cpu.Write(cpu.getA16(), cpu.a)
		return

	case 0xEE: // XOR A, n8
		cpu.xor8(&cpu.a, cpu.getN8())
		cpu.tick()
		return
	case 0xEF: // RST 28H
		cpu.call16(0x28)
		return

	case 0xF0: // LDH A, (a8)
		cpu.a = cpu.bus.Read(cpu.getA8())
		cpu.tick()
		return

	case 0xF1: // POP AF
		cpu.popAF()
		return

	case 0xF2: // LDH A, (C)
		cpu.a = cpu.bus.Read(0xFF00 + uint16(cpu.c))
		cpu.tick()
		return

	case 0xF3: // DI
		cpu.di()
		return

	case 0xF5: // PUSH AF
		cpu.pushAF()
		return

	case 0xF6: // OR A, n8
		cpu.or8(&cpu.a, cpu.getN8())
		cpu.tick()
		return

	case 0xF7: // RST 30H
		cpu.call16(0x30)
		return

	case 0xF8: // LD HL, SP+e8
		cpu.ld_HL_SP_e8()
		return

	case 0xF9: // LD SP, HL
		cpu.sp = cpu.getHL()
		cpu.tick()
		return

	case 0xFA: // LD A, (a16)
		cpu.a = cpu.bus.Read(cpu.getA16())
		cpu.tick()
		return

	case 0xFB: // EI
		cpu.ei()
		return

	case 0xFE: // CP A, n8
		cpu.cp8(cpu.a, cpu.getN8())
		cpu.tick()
		return

	case 0xFF: // RST 38H
		cpu.call16(0x38)
		return

	default:
		fmt.Printf("Instrucción no implementada: %02X\n", opcode)
		panic("Detenido")
	}
}
