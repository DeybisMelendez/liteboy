package cartridge

type mbc1 struct {
	ROM         [][0x4000]byte  // slices de bancos ROM de 16KiB cada uno
	ERAM        [4][0x2000]byte // hasta 4 bancos de 8KiB ERAM
	ramEnabled  bool
	romBankLow5 byte // bits bajos (5 bits) del banco ROM
	ramBank     byte // 2 bits para banco RAM o bits altos ROM
	bankingMode byte // 0=ROM banking mode, 1=RAM banking mode
}

func (m *mbc1) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		if m.bankingMode == 0 {
			// modo 0: banco fijo 0
			return m.ROM[0][addr]
		} else {
			// modo 1: banco alto se aplica a 0000-3FFF
			bank := int(m.ramBank)<<5 | int(m.romBankLow5)
			// en bank 0 el banco 0 se interpreta como banco 1
			if bank&0x1F == 0 {
				bank |= 1
			}
			bank %= len(m.ROM)
			return m.ROM[bank][addr]
		}

	case addr >= 0x4000 && addr < 0x8000:
		bank := int(m.romBankLow5) | (int(m.ramBank) << 5)
		if bank&0x1F == 0 {
			bank |= 1
		}
		bank %= len(m.ROM)
		offset := addr - 0x4000
		return m.ROM[bank][offset]

	case addr >= 0xA000 && addr < 0xC000:
		if !m.ramEnabled {
			// RAM no habilitada, devuelve valor abierto
			return 0xFF
		}
		// Acceso a ERAM
		ramBank := 0
		if m.bankingMode == 1 {
			ramBank = int(m.ramBank) & 0x03
		}
		offset := addr - 0xA000
		return m.ERAM[ramBank][offset]
	}
	return 0xFF
}

func (m *mbc1) Write(addr uint16, value byte) {
	switch {
	case addr < 0x2000:
		// RAM enable (0x0A habilita, otro valor deshabilita)
		m.ramEnabled = (value & 0x0F) == 0x0A

	case addr >= 0x2000 && addr < 0x4000:
		// ROM bank lower 5 bits
		value &= 0x1F
		if value == 0 {
			value = 1
		}
		m.romBankLow5 = value

	case addr >= 0x4000 && addr < 0x6000:
		// RAM bank number or ROM bank upper bits (2 bits)
		m.ramBank = value & 0x03

	case addr >= 0x6000 && addr < 0x8000:
		// Banking mode select (0=ROM banking mode, 1=RAM banking mode)
		m.bankingMode = value & 0x01
	}
}
