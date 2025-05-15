package bus

import (
	"log"

	"github.com/deybismelendez/liteboy/cartridge"
)

const (
	DIVRegister  = 0xFF04
	TIMARegister = 0xFF05
	TACRegister  = 0xFF07
)

type Bus struct {
	cart       *cartridge.Cartridge
	BootROM    [0x100]byte
	bootActive bool
	ROM00      *[0x4000]byte // 0x0000 - 0x3FFF
	ROMNN      *[0x4000]byte // 0x4000 - 0x7FFF
	VRAM       [0x2000]byte  // 0x8000 - 0x9FFF
	ERAM       *[0x2000]byte // 0xA000 - 0xBFFF
	WRAM       [0x2000]byte  // 0xC000 - 0xDFFF
	OAM        [0xA0]byte    // 0xFE00 - 0xFE9F
	IO         [0x80]byte    // 0xFF00 - 0xFF7F
	HRAM       [0x7F]byte    // 0xFF80 - 0xFFFE
	IE         byte          // 0xFFFF
}

func NewBus(cart *cartridge.Cartridge) *Bus {
	bus := &Bus{
		cart:       cart,
		BootROM:    BootROM,
		bootActive: false,
		ROM00:      &cart.ROM[0],
		ROMNN:      &cart.ROM[1],
		ERAM:       &[0x2000]byte{},
	}

	// Valores por defecto de los registros, igual que hiciste
	bus.Write(0xFF00, 0xCF) // P1
	bus.Write(0xFF01, 0x00) // SB
	bus.Write(0xFF02, 0x7E) // SC
	bus.Write(0xFF04, 0xAB) // DIV
	bus.Write(0xFF05, 0x00) // TIMA
	bus.Write(0xFF06, 0x00) // TMA
	bus.Write(0xFF07, 0xF8) // TAC
	bus.Write(0xFF0F, 0xE1) // IF

	// Audio registers
	bus.Write(0xFF10, 0x80)
	bus.Write(0xFF11, 0xBF)
	bus.Write(0xFF12, 0xF3)
	bus.Write(0xFF13, 0xFF)
	bus.Write(0xFF14, 0xBF)
	bus.Write(0xFF16, 0x3F)
	bus.Write(0xFF17, 0x00)
	bus.Write(0xFF18, 0xFF)
	bus.Write(0xFF19, 0xBF)
	bus.Write(0xFF1A, 0x7F)
	bus.Write(0xFF1B, 0xFF)
	bus.Write(0xFF1C, 0x9F)
	bus.Write(0xFF1D, 0xFF)
	bus.Write(0xFF1E, 0xBF)
	bus.Write(0xFF20, 0xFF)
	bus.Write(0xFF21, 0x00)
	bus.Write(0xFF22, 0x00)
	bus.Write(0xFF23, 0xBF)
	bus.Write(0xFF24, 0x77)
	bus.Write(0xFF25, 0xF3)
	bus.Write(0xFF26, 0xF0) // 0xF0 = DMG, 0xF1 = CGB

	// PPU
	bus.Write(0xFF40, 0x91) // LCDC
	bus.Write(0xFF41, 0x85) // STAT (o 0x81 también se ve)
	bus.Write(0xFF42, 0x00) // SCY
	bus.Write(0xFF43, 0x00) // SCX
	bus.Write(0xFF44, 0x00) // LY
	bus.Write(0xFF45, 0x00) // LYC
	bus.Write(0xFF46, 0xFF) // DMA
	bus.Write(0xFF47, 0xFC) // BGP
	bus.Write(0xFF48, 0xFF) // OBP0
	bus.Write(0xFF49, 0xFF) // OBP1
	bus.Write(0xFF4A, 0x00) // WY
	bus.Write(0xFF4B, 0x00) // WX

	bus.Write(0xFFFF, 0x00) // IE

	return bus
}

func (b *Bus) Read(addr uint16) byte {
	switch {
	case addr < 0x100 && b.bootActive:
		return b.BootROM[addr]
	case addr < 0x4000:
		return b.ROM00[addr]

	case addr < 0x8000:
		return b.ROMNN[addr-0x4000]

	case addr >= 0x8000 && addr < 0xA000:
		return b.VRAM[addr-0x8000]

	case addr >= 0xA000 && addr < 0xC000:
		return b.ERAM[addr-0xA000]

	case addr >= 0xC000 && addr < 0xE000:
		return b.WRAM[addr-0xC000]

	case addr >= 0xE000 && addr < 0xFE00:
		// Echo RAM (mirror of C000–DDFF)
		return b.WRAM[addr-0xE000]

	case addr >= 0xFE00 && addr < 0xFEA0:
		return b.OAM[addr-0xFE00]

	case addr >= 0xFEA0 && addr < 0xFF00:
		log.Printf("Intento de lectura en zona no usable en %04X, se retorna 0xFF\n", addr)
		return 0xFF

	case addr >= 0xFF00 && addr < 0xFF80:
		return b.IO[addr-0xFF00]

	case addr >= 0xFF80 && addr < 0xFFFF:
		return b.HRAM[addr-0xFF80]

	case addr == 0xFFFF:
		return b.IE

	default:
		log.Printf("Intento de lectura fuera de rango en %04X\n", addr)
		return 0xFF
	}
}
func (b *Bus) Write(addr uint16, value byte) {
	switch {
	case addr < 0x8000:
		log.Printf("Intento de escritura en ROM en %04X: %02X\n", addr, value)
		return

	case addr >= 0x8000 && addr < 0xA000:
		b.VRAM[addr-0x8000] = value

	case addr >= 0xA000 && addr < 0xC000:
		b.ERAM[addr-0xA000] = value

	case addr >= 0xC000 && addr < 0xE000:
		b.WRAM[addr-0xC000] = value

	case addr >= 0xE000 && addr < 0xFE00:
		// Echo RAM (mirror of C000–DDFF)
		b.WRAM[addr-0xE000] = value

	case addr >= 0xFE00 && addr < 0xFEA0:
		b.OAM[addr-0xFE00] = value

	case addr >= 0xFEA0 && addr < 0xFF00:
		log.Printf("Intento de escritura en zona no usable en %04X: %02X\n", addr, value)

	case addr >= 0xFF00 && addr < 0xFF80:
		// Si se intenta escribir en DIV se establece en 0
		//https://gbdev.io/pandocs/Timer_and_Divider_Registers.html#ff04--div-divider-register
		if addr == DIVRegister {
			b.IO[addr-0xFF00] = 0x00
			return
		}
		if addr == 0xFF46 {
			b.IO[addr-0xFF00] = value
			b.doDMATransfer(value)
			return
		}
		if addr == 0xFF50 && value != 0 {
			b.bootActive = false // Desactiva Boot ROM
		}
		b.IO[addr-0xFF00] = value

	case addr >= 0xFF80 && addr < 0xFFFF:
		b.HRAM[addr-0xFF80] = value

	case addr == 0xFFFF:
		b.IE = value

	default:
		log.Printf("Intento de escritura fuera de rango en %04X: %02X\n", addr, value)
	}
}
func (b *Bus) doDMATransfer(value byte) {
	source := uint16(value) << 8 // Dirección de origen base
	for i := 0; i < 0xA0; i++ {
		data := b.Read(source + uint16(i))
		b.OAM[i] = data
	}
}
