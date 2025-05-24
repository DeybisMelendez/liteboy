package cartridge

import (
	"log"
	"os"
)

type mbc3 struct {
	ROM        [][0x4000]byte  // ROM dividida en bancos de 16KiB
	ERAM       [4][0x2000]byte // hasta 4 bancos de 8KiB RAM
	rtcRegs    [5]byte         // RTC registers: 0→seconds, 1→minutes, 2→hours, 3→day low, 4→day high/control
	rtcLatch   byte            // latched RTC register index
	ramEnabled bool
	romBank    byte // banco de ROM actual (7 bits)
	ramBank    byte // banco RAM o registro RTC seleccionado
	latchClock byte // estado del latch clock
}

func (m *mbc3) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		return m.ROM[0][addr]

	case addr >= 0x4000 && addr < 0x8000:
		bank := m.romBank
		if bank == 0 {
			bank = 1
		}
		bank %= byte(len(m.ROM))
		offset := addr - 0x4000
		return m.ROM[bank][offset]

	case addr >= 0xA000 && addr < 0xC000:
		if !m.ramEnabled {
			return 0xFF
		}

		if m.ramBank <= 0x03 {
			offset := addr - 0xA000
			return m.ERAM[m.ramBank][offset]
		} else if m.ramBank >= 0x08 && m.ramBank <= 0x0C {
			return m.rtcRegs[m.ramBank-0x08]
		}
	}
	return 0xFF
}

func (m *mbc3) Write(addr uint16, value byte) {
	switch {
	case addr < 0x2000:
		m.ramEnabled = (value & 0x0F) == 0x0A

	case addr >= 0x2000 && addr < 0x4000:
		value &= 0x7F
		if value == 0 {
			value = 1
		}
		m.romBank = value

	case addr >= 0x4000 && addr < 0x6000:
		m.ramBank = value

	case addr >= 0x6000 && addr < 0x8000:
		// Latch clock data (0x00 → 0x01)
		if m.latchClock == 0x00 && value == 0x01 {
			// Aquí se debería hacer latch del RTC real, si fuera implementado
			log.Println("RTC latched (simulado)")
		}
		m.latchClock = value

	case addr >= 0xA000 && addr < 0xC000:
		if !m.ramEnabled {
			return
		}
		if m.ramBank <= 0x03 {
			offset := addr - 0xA000
			m.ERAM[m.ramBank][offset] = value
		} else if m.ramBank >= 0x08 && m.ramBank <= 0x0C {
			m.rtcRegs[m.ramBank-0x08] = value
		}
	}
}

// SaveToFile guarda RAM y registros RTC a archivo .sav
func (m *mbc3) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Guarda 4 bancos de RAM (4x8KiB)
	for _, bank := range m.ERAM {
		if _, err := file.Write(bank[:]); err != nil {
			return err
		}
	}
	// Guarda RTC
	if _, err := file.Write(m.rtcRegs[:]); err != nil {
		return err
	}
	return nil
}

// LoadFromFile carga RAM y registros RTC desde archivo .sav
func (m *mbc3) LoadFromFile(filename string) error {
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
	if _, err := file.Read(m.rtcRegs[:]); err != nil {
		return err
	}
	return nil
}
