package cartridge

import (
	"log"
	"os"
)

type mbc2 struct {
	ROM        [][0x4000]byte // Bancos ROM de 16 KiB
	ERAM       [512]byte      // RAM interna de MBC2 (512 x 4 bits)
	ramEnabled bool
	romBank    byte // Solo 4 bits válidos (1-15)
}

func (m *mbc2) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		// Banco fijo 0
		return m.ROM[0][addr]

	case addr >= 0x4000 && addr < 0x8000:
		// Banco conmutable (1-15)
		bank := m.romBank
		if bank == 0 {
			bank = 1
		}
		bank %= byte(len(m.ROM))
		offset := addr - 0x4000
		return m.ROM[bank][offset]

	case addr >= 0xA000 && addr < 0xA200:
		if !m.ramEnabled {
			return 0xFF
		}
		offset := addr - 0xA000
		val := m.ERAM[offset]
		return val | 0xF0 // solo 4 bits significativos

	default:
		return 0xFF
	}
}

func (m *mbc2) Write(addr uint16, value byte) {
	switch {
	case addr < 0x2000:
		// RAM Enable solo si el bit 8 de la dirección es 0
		if (addr & 0x0100) == 0 {
			m.ramEnabled = (value & 0x0F) == 0x0A
		}

	case addr >= 0x2000 && addr < 0x4000:
		// Cambiar banco ROM solo si el bit 8 de la dirección es 1
		if (addr & 0x0100) != 0 {
			m.romBank = value & 0x0F // Solo 4 bits válidos (1-15)
			if m.romBank == 0 {
				m.romBank = 1
			}
		}

	case addr >= 0xA000 && addr < 0xA200:
		if !m.ramEnabled {
			return
		}
		offset := addr - 0xA000
		m.ERAM[offset] = value & 0x0F // Solo 4 bits significativos
	default:
		log.Printf("Write ignorado en dirección %04X valor %02X", addr, value)
	}
}

// SaveToFile guarda la RAM de MBC2 en un archivo .sav
func (m *mbc2) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(m.ERAM[:])
	return err
}

// LoadFromFile carga la RAM desde un archivo .sav
func (m *mbc2) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Read(m.ERAM[:])
	return err
}
