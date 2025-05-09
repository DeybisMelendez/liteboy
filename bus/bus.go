package bus

import (
	"fmt"
	"log"

	"github.com/deybismelendez/liteboy/cartridge"
)

type Bus struct {
	cart  *cartridge.Cartridge
	ROM00 *[0x4000]byte // 0x0000 - 0x3FFF
	ROMNN *[0x4000]byte // 0x4000 - 0x7FFF
	VRAM  [0x2000]byte  // 0x8000 - 0x9FFF
	ERAM  *[0x2000]byte // 0xA000 - 0xBFFF
	WRAM  [0x2000]byte  // 0xC000 - 0xDFFF
	OAM   [0xA0]byte    // 0xFE00 - 0xFE9F
	IO    [0x80]byte    // 0xFF00 - 0xFF7F
	HRAM  [0x7F]byte    // 0xFF80 - 0xFFFE
	IE    byte          // 0xFFFF
}

func NewBus(cart *cartridge.Cartridge) *Bus {
	return &Bus{
		cart:  cart,
		ROM00: &cart.ROM[0],
		ROMNN: &cart.ROM[1],
		ERAM:  &[0x2000]byte{},
		// TODO: ERAM
	}
}

func (b *Bus) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		return b.ROM00[addr]

	case addr < 0x8000:
		return b.ROMNN[addr-0x4000]

	case addr < 0xA000:
		return b.VRAM[addr-0x8000]

	case addr < 0xC000:
		return b.VRAM[addr-0xA000]

	case addr < 0xE000:
		return b.WRAM[addr-0xC000]

	case addr < 0xFEA0:
		fmt.Println(addr)
		return b.OAM[addr-0xFE00]

	case addr < 0xFF80:
		return b.IO[addr-0xFF00]

	case addr < 0xFFFF:
		return b.HRAM[addr-0xFF80]

	case addr == 0xFFFF:
		return b.IE

	default:
		panic(fmt.Sprintf("Intento de lectura fuera de rango: %04X", addr))
	}
}

func (b *Bus) Write(addr uint16, value byte) {
	switch {
	case addr < 0x8000:
		// ROM: no se puede escribir.
		// TODO: En el futuro aquí iría MBC.
		log.Printf("Intento de escritura en ROM en %04X: %02X", addr, value)

	case addr < 0xA000:
		b.VRAM[addr-0x8000] = value

	case addr < 0xC000:
		b.VRAM[addr-0xA000] = value

	case addr < 0xE000:
		b.WRAM[addr-0xC000] = value

	case addr < 0xFEA0:
		b.OAM[addr-0xFE00] = value

	case addr < 0xFF80:
		b.IO[addr-0xFF00] = value

	case addr < 0xFFFF:
		b.HRAM[addr-0xFF80] = value

	default:
		panic(fmt.Sprintf("Intento de escritura fuera de rango: %04X", addr))
	}
}
