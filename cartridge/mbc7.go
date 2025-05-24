package cartridge

import (
	"log"
)

type mbc7 struct {
	ROM         [][0x4000]byte // Bancos de ROM
	eeprom      [256]byte      // EEPROM de 256 bytes
	romBank     uint16         // Banco ROM actual
	ramEnabled  bool           // Habilita acceso EEPROM
	eepromAddr  uint8          // Dirección EEPROM seleccionada
	eepromCmd   byte           // Comando EEPROM actual
	eepromWrite bool           // Bandera de escritura activa
	tiltX       byte           // Valor del sensor de inclinación X
	tiltY       byte           // Valor del sensor de inclinación Y
}

func (m *mbc7) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		return m.ROM[0][addr]

	case addr >= 0x4000 && addr < 0x8000:
		bank := int(m.romBank) % len(m.ROM)
		return m.ROM[bank][addr-0x4000]

	case addr >= 0xA000 && addr < 0xA200:
		if !m.ramEnabled {
			return 0xFF
		}
		return m.readEEPROM(addr)
	}
	return 0xFF
}

func (m *mbc7) Write(addr uint16, value byte) {
	switch {
	case addr < 0x2000:
		m.ramEnabled = (value & 0x0F) == 0x0A

	case addr >= 0x2000 && addr < 0x4000:
		m.romBank = uint16(value)

	case addr >= 0xA000 && addr < 0xA200:
		if m.ramEnabled {
			m.writeEEPROM(addr, value)
		}
	}
}

// EEPROM serial emulation (simplificada)
func (m *mbc7) readEEPROM(addr uint16) byte {
	offset := addr - 0xA000
	switch offset {
	case 0x0000: // Simula lectura de datos
		if m.eepromWrite {
			return 0xFF
		}
		return m.eeprom[m.eepromAddr]
	case 0x0010: // Sensor Tilt-X
		return m.tiltX
	case 0x0011: // Sensor Tilt-Y
		return m.tiltY
	}
	return 0xFF
}

func (m *mbc7) writeEEPROM(addr uint16, value byte) {
	offset := addr - 0xA000
	switch offset {
	case 0x0000: // Enviar comando de control EEPROM
		m.eepromCmd = value
		m.eepromWrite = (value & 0x80) != 0
	case 0x0001: // Dirección EEPROM
		m.eepromAddr = value & 0xFF
	case 0x0002: // Escribir datos (si escritura habilitada)
		if m.eepromWrite {
			m.eeprom[m.eepromAddr] = value
		}
	case 0x0010: // Simula sensor Tilt-X (para testeo manual)
		m.tiltX = value
	case 0x0011: // Simula sensor Tilt-Y
		m.tiltY = value
	default:
		log.Printf("MBC7 Write %04X: %02X\n", addr, value)
	}
}
