package cartridge

import "log"

type romOnly struct {
	ROM [][0x4000]byte
}

func (r *romOnly) Read(addr uint16) byte {
	if addr < 0x8000 {
		bank := addr / 0x4000
		offset := addr % 0x4000
		if int(bank) < len(r.ROM) {
			return r.ROM[bank][offset]
		}
	}
	return 0xFF
}

func (r *romOnly) Write(addr uint16, value byte) {
	// romOnly no permite escritura
	log.Fatalf("Intento de escritura en ROM en %04X: %02X\n", addr, value)
}
