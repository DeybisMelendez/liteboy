package cartridge

import (
	"os"
)

type mbc5 struct {
	ROM        [][0x4000]byte   // Bancos de 16KiB de ROM
	ERAM       [16][0x2000]byte // Hasta 16 bancos de 8KiB RAM
	ramEnabled bool
	romBank    uint16 // Banco de ROM (9 bits → 0x000–0x1FF)
	ramBank    byte   // Banco de RAM (0x00–0x0F)
}

func (m *mbc5) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		return m.ROM[0][addr]

	case addr >= 0x4000 && addr < 0x8000:
		bank := m.romBank % uint16(len(m.ROM))
		return m.ROM[bank][addr-0x4000]

	case addr >= 0xA000 && addr < 0xC000:
		if !m.ramEnabled {
			return 0xFF
		}
		offset := addr - 0xA000
		return m.ERAM[m.ramBank][offset]
	}
	return 0xFF
}

func (m *mbc5) Write(addr uint16, value byte) {
	switch {
	case addr < 0x2000:
		m.ramEnabled = (value & 0x0F) == 0x0A

	case addr >= 0x2000 && addr < 0x3000:
		m.romBank = (m.romBank & 0x100) | uint16(value) // LSB 8 bits

	case addr >= 0x3000 && addr < 0x4000:
		m.romBank = (m.romBank & 0xFF) | (uint16(value&0x01) << 8) // 9º bit

	case addr >= 0x4000 && addr < 0x6000:
		m.ramBank = value & 0x0F // solo 4 bits permitidos

	case addr >= 0xA000 && addr < 0xC000:
		if !m.ramEnabled {
			return
		}
		offset := addr - 0xA000
		m.ERAM[m.ramBank][offset] = value
	}
}

// SaveToFile guarda el contenido de la RAM en un archivo
func (m *mbc5) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, bank := range m.ERAM {
		if _, err := file.Write(bank[:]); err != nil {
			return err
		}
	}
	return nil
}

// LoadFromFile carga la RAM desde un archivo .sav
func (m *mbc5) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := range m.ERAM {
		if _, err := file.Read(m.ERAM[i][:]); err != nil {
			return err
		}
	}
	return nil
}
