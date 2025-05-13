package cpu

import "log"

func (cpu *CPU) runCB() {
	cbOpcode := cpu.getN8()

	switch cbOpcode {
	case 0x00: // RLC B
		cpu.b = cpu.rlc(cpu.b)
	case 0x01: // RLC C
		cpu.c = cpu.rlc(cpu.c)
	case 0x02: // RLC D
		cpu.d = cpu.rlc(cpu.d)
	case 0x03: // RLC E
		cpu.e = cpu.rlc(cpu.e)
	case 0x04: // RLC H
		cpu.h = cpu.rlc(cpu.h)
	case 0x05: // RLC L
		cpu.l = cpu.rlc(cpu.l)
	case 0x06: // RLC (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.rlc(value)
		cpu.bus.Write(addr, value)
	case 0x07: // RLC A
		cpu.a = cpu.rlc(cpu.a)

	case 0x08: // RRC B
		cpu.b = cpu.rrc(cpu.b)
	case 0x09: // RRC C
		cpu.c = cpu.rrc(cpu.c)
	case 0x0A: // RRC D
		cpu.d = cpu.rrc(cpu.d)
	case 0x0B: // RRC E
		cpu.e = cpu.rrc(cpu.e)
	case 0x0C: // RRC H
		cpu.h = cpu.rrc(cpu.h)
	case 0x0D: // RRC L
		cpu.l = cpu.rrc(cpu.l)
	case 0x0E: // RRC (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.rrc(value)
		cpu.bus.Write(addr, value)
	case 0x0F: // RRC A
		cpu.a = cpu.rrc(cpu.a)
	case 0x10: // RL B
		cpu.b = cpu.rl(cpu.b)
	case 0x11: // RL C
		cpu.c = cpu.rl(cpu.c)
	case 0x12: // RL D
		cpu.d = cpu.rl(cpu.d)
	case 0x13: // RL E
		cpu.e = cpu.rl(cpu.e)
	case 0x14: // RL H
		cpu.h = cpu.rl(cpu.h)
	case 0x15: // RL L
		cpu.l = cpu.rl(cpu.l)
	case 0x16: // RL (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.rl(value)
		cpu.bus.Write(addr, value)
	case 0x17: // RL A
		cpu.a = cpu.rl(cpu.a)

	case 0x18: // RR B
		cpu.b = cpu.rr(cpu.b)
	case 0x19: // RR C
		cpu.c = cpu.rr(cpu.c)
	case 0x1A: // RR D
		cpu.d = cpu.rr(cpu.d)
	case 0x1B: // RR E
		cpu.e = cpu.rr(cpu.e)
	case 0x1C: // RR H
		cpu.h = cpu.rr(cpu.h)
	case 0x1D: // RR L
		cpu.l = cpu.rr(cpu.l)
	case 0x1E: // RR (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.rr(value)
		cpu.bus.Write(addr, value)
	case 0x1F: // RR A
		cpu.a = cpu.rr(cpu.a)
	case 0x20: // SLA B
		cpu.b = cpu.sla(cpu.b)
	case 0x21: // SLA C
		cpu.c = cpu.sla(cpu.c)
	case 0x22: // SLA D
		cpu.d = cpu.sla(cpu.d)
	case 0x23: // SLA E
		cpu.e = cpu.sla(cpu.e)
	case 0x24: // SLA H
		cpu.h = cpu.sla(cpu.h)
	case 0x25: // SLA L
		cpu.l = cpu.sla(cpu.l)
	case 0x26: // SLA (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.sla(value)
		cpu.bus.Write(addr, value)
	case 0x27: // SLA A
		cpu.a = cpu.sla(cpu.a)

	case 0x28: // SRA B
		cpu.b = cpu.sra(cpu.b)
	case 0x29: // SRA C
		cpu.c = cpu.sra(cpu.c)
	case 0x2A: // SRA D
		cpu.d = cpu.sra(cpu.d)
	case 0x2B: // SRA E
		cpu.e = cpu.sra(cpu.e)
	case 0x2C: // SRA H
		cpu.h = cpu.sra(cpu.h)
	case 0x2D: // SRA L
		cpu.l = cpu.sra(cpu.l)
	case 0x2E: // SRA (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.sra(value)
		cpu.bus.Write(addr, value)
	case 0x2F: // SRA A
		cpu.a = cpu.sra(cpu.a)
	case 0x30: // SWAP B
		cpu.b = cpu.swap(cpu.b)
	case 0x31: // SWAP C
		cpu.c = cpu.swap(cpu.c)
	case 0x32: // SWAP D
		cpu.d = cpu.swap(cpu.d)
	case 0x33: // SWAP E
		cpu.e = cpu.swap(cpu.e)
	case 0x34: // SWAP H
		cpu.h = cpu.swap(cpu.h)
	case 0x35: // SWAP L
		cpu.l = cpu.swap(cpu.l)
	case 0x36: // SWAP (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.swap(value)
		cpu.bus.Write(addr, value)
	case 0x37: // SWAP A
		cpu.a = cpu.swap(cpu.a)

	case 0x38: // SRL B
		cpu.b = cpu.srl(cpu.b)
	case 0x39: // SRL C
		cpu.c = cpu.srl(cpu.c)
	case 0x3A: // SRL D
		cpu.d = cpu.srl(cpu.d)
	case 0x3B: // SRL E
		cpu.e = cpu.srl(cpu.e)
	case 0x3C: // SRL H
		cpu.h = cpu.srl(cpu.h)
	case 0x3D: // SRL L
		cpu.l = cpu.srl(cpu.l)
	case 0x3E: // SRL (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		value = cpu.srl(value)
		cpu.bus.Write(addr, value)
	case 0x3F: // SRL A
		cpu.a = cpu.srl(cpu.a)
	case 0x40: // BIT 0, B
		cpu.bit(0, cpu.b)
	case 0x41: // BIT 0, C
		cpu.bit(0, cpu.c)
	case 0x42: // BIT 0, D
		cpu.bit(0, cpu.d)
	case 0x43: // BIT 0, E
		cpu.bit(0, cpu.e)
	case 0x44: // BIT 0, H
		cpu.bit(0, cpu.h)
	case 0x45: // BIT 0, L
		cpu.bit(0, cpu.l)
	case 0x46: // BIT 0, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(0, value)
	case 0x47: // BIT 0, A
		cpu.bit(0, cpu.a)

	case 0x48: // BIT 1, B
		cpu.bit(1, cpu.b)
	case 0x49: // BIT 1, C
		cpu.bit(1, cpu.c)
	case 0x4A: // BIT 1, D
		cpu.bit(1, cpu.d)
	case 0x4B: // BIT 1, E
		cpu.bit(1, cpu.e)
	case 0x4C: // BIT 1, H
		cpu.bit(1, cpu.h)
	case 0x4D: // BIT 1, L
		cpu.bit(1, cpu.l)
	case 0x4E: // BIT 1, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(1, value)
	case 0x4F: // BIT 1, A
		cpu.bit(1, cpu.a)
	case 0x50: // BIT 2, B
		cpu.bit(2, cpu.b)
	case 0x51: // BIT 2, C
		cpu.bit(2, cpu.c)
	case 0x52: // BIT 2, D
		cpu.bit(2, cpu.d)
	case 0x53: // BIT 2, E
		cpu.bit(2, cpu.e)
	case 0x54: // BIT 2, H
		cpu.bit(2, cpu.h)
	case 0x55: // BIT 2, L
		cpu.bit(2, cpu.l)
	case 0x56: // BIT 2, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(2, value)
	case 0x57: // BIT 2, A
		cpu.bit(2, cpu.a)

	case 0x58: // BIT 3, B
		cpu.bit(3, cpu.b)
	case 0x59: // BIT 3, C
		cpu.bit(3, cpu.c)
	case 0x5A: // BIT 3, D
		cpu.bit(3, cpu.d)
	case 0x5B: // BIT 3, E
		cpu.bit(3, cpu.e)
	case 0x5C: // BIT 3, H
		cpu.bit(3, cpu.h)
	case 0x5D: // BIT 3, L
		cpu.bit(3, cpu.l)
	case 0x5E: // BIT 3, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(3, value)
	case 0x5F: // BIT 3, A
		cpu.bit(3, cpu.a)
	case 0x60: // BIT 4, B
		cpu.bit(4, cpu.b)
	case 0x61: // BIT 4, C
		cpu.bit(4, cpu.c)
	case 0x62: // BIT 4, D
		cpu.bit(4, cpu.d)
	case 0x63: // BIT 4, E
		cpu.bit(4, cpu.e)
	case 0x64: // BIT 4, H
		cpu.bit(4, cpu.h)
	case 0x65: // BIT 4, L
		cpu.bit(4, cpu.l)
	case 0x66: // BIT 4, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(4, value)
	case 0x67: // BIT 4, A
		cpu.bit(4, cpu.a)

	case 0x68: // BIT 5, B
		cpu.bit(5, cpu.b)
	case 0x69: // BIT 5, C
		cpu.bit(5, cpu.c)
	case 0x6A: // BIT 5, D
		cpu.bit(5, cpu.d)
	case 0x6B: // BIT 5, E
		cpu.bit(5, cpu.e)
	case 0x6C: // BIT 5, H
		cpu.bit(5, cpu.h)
	case 0x6D: // BIT 5, L
		cpu.bit(5, cpu.l)
	case 0x6E: // BIT 5, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(5, value)
	case 0x6F: // BIT 5, A
		cpu.bit(5, cpu.a)
	case 0x70: // BIT 6, B
		cpu.bit(6, cpu.b)
	case 0x71: // BIT 6, C
		cpu.bit(6, cpu.c)
	case 0x72: // BIT 6, D
		cpu.bit(6, cpu.d)
	case 0x73: // BIT 6, E
		cpu.bit(6, cpu.e)
	case 0x74: // BIT 6, H
		cpu.bit(6, cpu.h)
	case 0x75: // BIT 6, L
		cpu.bit(6, cpu.l)
	case 0x76: // BIT 6, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(6, value)
	case 0x77: // BIT 6, A
		cpu.bit(6, cpu.a)

	case 0x78: // BIT 7, B
		cpu.bit(7, cpu.b)
	case 0x79: // BIT 7, C
		cpu.bit(7, cpu.c)
	case 0x7A: // BIT 7, D
		cpu.bit(7, cpu.d)
	case 0x7B: // BIT 7, E
		cpu.bit(7, cpu.e)
	case 0x7C: // BIT 7, H
		cpu.bit(7, cpu.h)
	case 0x7D: // BIT 7, L
		cpu.bit(7, cpu.l)
	case 0x7E: // BIT 7, (HL)
		addr := cpu.getHL()
		value := cpu.bus.Read(addr)
		cpu.bit(7, value)
	case 0x7F: // BIT 7, A
		cpu.bit(7, cpu.a)
		// RES 0, r
	case 0x80:
		cpu.b = cpu.res(0, cpu.b)
	case 0x81:
		cpu.c = cpu.res(0, cpu.c)
	case 0x82:
		cpu.d = cpu.res(0, cpu.d)
	case 0x83:
		cpu.e = cpu.res(0, cpu.e)
	case 0x84:
		cpu.h = cpu.res(0, cpu.h)
	case 0x85:
		cpu.l = cpu.res(0, cpu.l)
	case 0x86:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(0, cpu.bus.Read(addr)))
	case 0x87:
		cpu.a = cpu.res(0, cpu.a)

	// RES 1, r
	case 0x88:
		cpu.b = cpu.res(1, cpu.b)
	case 0x89:
		cpu.c = cpu.res(1, cpu.c)
	case 0x8A:
		cpu.d = cpu.res(1, cpu.d)
	case 0x8B:
		cpu.e = cpu.res(1, cpu.e)
	case 0x8C:
		cpu.h = cpu.res(1, cpu.h)
	case 0x8D:
		cpu.l = cpu.res(1, cpu.l)
	case 0x8E:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(1, cpu.bus.Read(addr)))
	case 0x8F:
		cpu.a = cpu.res(1, cpu.a)

	// RES 2, r
	case 0x90:
		cpu.b = cpu.res(2, cpu.b)
	case 0x91:
		cpu.c = cpu.res(2, cpu.c)
	case 0x92:
		cpu.d = cpu.res(2, cpu.d)
	case 0x93:
		cpu.e = cpu.res(2, cpu.e)
	case 0x94:
		cpu.h = cpu.res(2, cpu.h)
	case 0x95:
		cpu.l = cpu.res(2, cpu.l)
	case 0x96:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(2, cpu.bus.Read(addr)))
	case 0x97:
		cpu.a = cpu.res(2, cpu.a)

	// RES 3, r
	case 0x98:
		cpu.b = cpu.res(3, cpu.b)
	case 0x99:
		cpu.c = cpu.res(3, cpu.c)
	case 0x9A:
		cpu.d = cpu.res(3, cpu.d)
	case 0x9B:
		cpu.e = cpu.res(3, cpu.e)
	case 0x9C:
		cpu.h = cpu.res(3, cpu.h)
	case 0x9D:
		cpu.l = cpu.res(3, cpu.l)
	case 0x9E:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(3, cpu.bus.Read(addr)))
	case 0x9F:
		cpu.a = cpu.res(3, cpu.a)

	// RES 4, r
	case 0xA0:
		cpu.b = cpu.res(4, cpu.b)
	case 0xA1:
		cpu.c = cpu.res(4, cpu.c)
	case 0xA2:
		cpu.d = cpu.res(4, cpu.d)
	case 0xA3:
		cpu.e = cpu.res(4, cpu.e)
	case 0xA4:
		cpu.h = cpu.res(4, cpu.h)
	case 0xA5:
		cpu.l = cpu.res(4, cpu.l)
	case 0xA6:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(4, cpu.bus.Read(addr)))
	case 0xA7:
		cpu.a = cpu.res(4, cpu.a)

	// RES 5, r
	case 0xA8:
		cpu.b = cpu.res(5, cpu.b)
	case 0xA9:
		cpu.c = cpu.res(5, cpu.c)
	case 0xAA:
		cpu.d = cpu.res(5, cpu.d)
	case 0xAB:
		cpu.e = cpu.res(5, cpu.e)
	case 0xAC:
		cpu.h = cpu.res(5, cpu.h)
	case 0xAD:
		cpu.l = cpu.res(5, cpu.l)
	case 0xAE:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(5, cpu.bus.Read(addr)))
	case 0xAF:
		cpu.a = cpu.res(5, cpu.a)

	// RES 6, r
	case 0xB0:
		cpu.b = cpu.res(6, cpu.b)
	case 0xB1:
		cpu.c = cpu.res(6, cpu.c)
	case 0xB2:
		cpu.d = cpu.res(6, cpu.d)
	case 0xB3:
		cpu.e = cpu.res(6, cpu.e)
	case 0xB4:
		cpu.h = cpu.res(6, cpu.h)
	case 0xB5:
		cpu.l = cpu.res(6, cpu.l)
	case 0xB6:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(6, cpu.bus.Read(addr)))
	case 0xB7:
		cpu.a = cpu.res(6, cpu.a)

	// RES 7, r
	case 0xB8:
		cpu.b = cpu.res(7, cpu.b)
	case 0xB9:
		cpu.c = cpu.res(7, cpu.c)
	case 0xBA:
		cpu.d = cpu.res(7, cpu.d)
	case 0xBB:
		cpu.e = cpu.res(7, cpu.e)
	case 0xBC:
		cpu.h = cpu.res(7, cpu.h)
	case 0xBD:
		cpu.l = cpu.res(7, cpu.l)
	case 0xBE:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.res(7, cpu.bus.Read(addr)))
	case 0xBF:
		cpu.a = cpu.res(7, cpu.a)
		// SET 0, r
	case 0xC0:
		cpu.b = cpu.set(0, cpu.b)
	case 0xC1:
		cpu.c = cpu.set(0, cpu.c)
	case 0xC2:
		cpu.d = cpu.set(0, cpu.d)
	case 0xC3:
		cpu.e = cpu.set(0, cpu.e)
	case 0xC4:
		cpu.h = cpu.set(0, cpu.h)
	case 0xC5:
		cpu.l = cpu.set(0, cpu.l)
	case 0xC6:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(0, cpu.bus.Read(addr)))
	case 0xC7:
		cpu.a = cpu.set(0, cpu.a)

	// SET 1, r
	case 0xC8:
		cpu.b = cpu.set(1, cpu.b)
	case 0xC9:
		cpu.c = cpu.set(1, cpu.c)
	case 0xCA:
		cpu.d = cpu.set(1, cpu.d)
	case 0xCB:
		cpu.e = cpu.set(1, cpu.e)
	case 0xCC:
		cpu.h = cpu.set(1, cpu.h)
	case 0xCD:
		cpu.l = cpu.set(1, cpu.l)
	case 0xCE:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(1, cpu.bus.Read(addr)))
	case 0xCF:
		cpu.a = cpu.set(1, cpu.a)

	// SET 2, r
	case 0xD0:
		cpu.b = cpu.set(2, cpu.b)
	case 0xD1:
		cpu.c = cpu.set(2, cpu.c)
	case 0xD2:
		cpu.d = cpu.set(2, cpu.d)
	case 0xD3:
		cpu.e = cpu.set(2, cpu.e)
	case 0xD4:
		cpu.h = cpu.set(2, cpu.h)
	case 0xD5:
		cpu.l = cpu.set(2, cpu.l)
	case 0xD6:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(2, cpu.bus.Read(addr)))
	case 0xD7:
		cpu.a = cpu.set(2, cpu.a)

	// SET 3, r
	case 0xD8:
		cpu.b = cpu.set(3, cpu.b)
	case 0xD9:
		cpu.c = cpu.set(3, cpu.c)
	case 0xDA:
		cpu.d = cpu.set(3, cpu.d)
	case 0xDB:
		cpu.e = cpu.set(3, cpu.e)
	case 0xDC:
		cpu.h = cpu.set(3, cpu.h)
	case 0xDD:
		cpu.l = cpu.set(3, cpu.l)
	case 0xDE:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(3, cpu.bus.Read(addr)))
	case 0xDF:
		cpu.a = cpu.set(3, cpu.a)

	// SET 4, r
	case 0xE0:
		cpu.b = cpu.set(4, cpu.b)
	case 0xE1:
		cpu.c = cpu.set(4, cpu.c)
	case 0xE2:
		cpu.d = cpu.set(4, cpu.d)
	case 0xE3:
		cpu.e = cpu.set(4, cpu.e)
	case 0xE4:
		cpu.h = cpu.set(4, cpu.h)
	case 0xE5:
		cpu.l = cpu.set(4, cpu.l)
	case 0xE6:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(4, cpu.bus.Read(addr)))
	case 0xE7:
		cpu.a = cpu.set(4, cpu.a)

	// SET 5, r
	case 0xE8:
		cpu.b = cpu.set(5, cpu.b)
	case 0xE9:
		cpu.c = cpu.set(5, cpu.c)
	case 0xEA:
		cpu.d = cpu.set(5, cpu.d)
	case 0xEB:
		cpu.e = cpu.set(5, cpu.e)
	case 0xEC:
		cpu.h = cpu.set(5, cpu.h)
	case 0xED:
		cpu.l = cpu.set(5, cpu.l)
	case 0xEE:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(5, cpu.bus.Read(addr)))
	case 0xEF:
		cpu.a = cpu.set(5, cpu.a)

	// SET 6, r
	case 0xF0:
		cpu.b = cpu.set(6, cpu.b)
	case 0xF1:
		cpu.c = cpu.set(6, cpu.c)
	case 0xF2:
		cpu.d = cpu.set(6, cpu.d)
	case 0xF3:
		cpu.e = cpu.set(6, cpu.e)
	case 0xF4:
		cpu.h = cpu.set(6, cpu.h)
	case 0xF5:
		cpu.l = cpu.set(6, cpu.l)
	case 0xF6:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(6, cpu.bus.Read(addr)))
	case 0xF7:
		cpu.a = cpu.set(6, cpu.a)

	// SET 7, r
	case 0xF8:
		cpu.b = cpu.set(7, cpu.b)
	case 0xF9:
		cpu.c = cpu.set(7, cpu.c)
	case 0xFA:
		cpu.d = cpu.set(7, cpu.d)
	case 0xFB:
		cpu.e = cpu.set(7, cpu.e)
	case 0xFC:
		cpu.h = cpu.set(7, cpu.h)
	case 0xFD:
		cpu.l = cpu.set(7, cpu.l)
	case 0xFE:
		addr := cpu.getHL()
		cpu.bus.Write(addr, cpu.set(7, cpu.bus.Read(addr)))
	case 0xFF:
		cpu.a = cpu.set(7, cpu.a)

	default:
		log.Fatalf("Instrucción CB no implementada: 0xCB 0x%02X en PC=%04X", cbOpcode, cpu.pc-1)
	}
}
