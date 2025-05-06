package bus

import (
	"fmt"
	"log"
)

type Bus struct {
	ROM  *[]byte      // 0x0000 - 0x7FFF (cartucho ROM)
	VRAM [0x2000]byte // 0x8000 - 0x9FFF
	ERAM [0x2000]byte // 0xA000 - 0xBFFF
	WRAM [0x2000]byte // 0xC000 - 0xDFFF
	OAM  [0xA0]byte   // 0xFE00 - 0xFE9F
	IO   [0x80]byte   // 0xFF00 - 0xFF7F
	HRAM [0x7F]byte   // 0xFF80 - 0xFFFE (High RAM)
	IE   byte         // 0xFFFF (Interrupt Enable)
}

func NewBus(rom *[]byte) *Bus {
	if len(*rom) < 0x8000 {
		log.Fatalf("ROM demasiado pequeña: se esperaban al menos 32KB")
	}
	return &Bus{
		ROM: rom,
	}
}
func (b *Bus) GetAddress(addr uint16) *byte {
	switch {
	case addr < 0x8000:
		rom := *b.ROM
		return &rom[addr]

	case addr >= 0x8000 && addr < 0xA000:
		return &b.VRAM[addr-0x8000]

	case addr >= 0xA000 && addr < 0xC000:
		return &b.VRAM[addr-0xA000]

	case addr >= 0xC000 && addr < 0xE000:
		return &b.WRAM[addr-0xC000]

	case addr >= 0xFE00 && addr < 0xFEA0:
		return &b.OAM[addr-0xFE00]

	case addr >= 0xFF00 && addr < 0xFF80:
		return &b.IO[addr-0xFF00]

	case addr >= 0xFF80 && addr < 0xFFFF:
		return &b.HRAM[addr-0xFF80]

	case addr == 0xFFFF:
		return &b.IE

	default:
		fmt.Printf("%02X ", addr)
		panic("No implementado o fuera de rango")
		// No implementado o fuera de rango
	}
}
func (b *Bus) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		rom := *b.ROM
		return rom[addr]

	case addr >= 0x8000 && addr < 0xA000:
		return b.VRAM[addr-0x8000]

	case addr >= 0xA000 && addr < 0xC000:
		return b.VRAM[addr-0xA000]

	case addr >= 0xC000 && addr < 0xE000:
		return b.WRAM[addr-0xC000]

	case addr >= 0xFE00 && addr < 0xFEA0:
		return b.OAM[addr-0xFE00]

	case addr >= 0xFF00 && addr < 0xFF80:
		return b.IO[addr-0xFF00]

	case addr >= 0xFF80 && addr < 0xFFFF:
		return b.HRAM[addr-0xFF80]

	case addr == 0xFFFF:
		return b.IE

	default:
		fmt.Printf("%02X ", addr)
		panic("No implementado o fuera de rango")
		// No implementado o fuera de rango
	}
}

func (b *Bus) Write(addr uint16, value byte) {
	switch {
	case addr < 0x8000:
		// ROM: no se puede escribir.
		// TODO: En el futuro aquí iría MBC.
		log.Printf("Intento de escritura en ROM en %04X: %02X", addr, value)

	case addr >= 0x8000 && addr < 0xA000:
		b.VRAM[addr-0x8000] = value

	case addr >= 0xA000 && addr < 0xC000:
		b.VRAM[addr-0xA000] = value

	case addr >= 0xC000 && addr < 0xE000:
		b.WRAM[addr-0xC000] = value

	case addr >= 0xFE00 && addr < 0xFEA0:
		b.OAM[addr-0xFE00] = value

	case addr >= 0xFF00 && addr < 0xFF80:
		b.IO[addr-0xFF00] = value

	case addr >= 0xFF80 && addr < 0xFFFF:
		b.HRAM[addr-0xFF80] = value

	default:
		fmt.Printf("%02X ", addr)
		panic("No implementado o fuera de rango")
		// No implementado o fuera de rango
	}
}
