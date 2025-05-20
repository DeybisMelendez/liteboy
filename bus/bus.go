package bus

import (
	"log"

	"github.com/deybismelendez/liteboy/cartridge"
)

const (
	DIVRegister  = 0xFF04
	TIMARegister = 0xFF05
	TACRegister  = 0xFF07

	ClientCPU     = 0
	ClientPPU     = 1
	ClientTimer   = 2
	ClientDMA     = 3
	ClientLiteBoy = 255
)

type Bus struct {
	cart       *cartridge.Cartridge
	BootROM    [0x100]byte
	bootActive bool
	//ROM00      *[0x4000]byte // 0x0000 - 0x3FFF
	//ROMNN      *[0x4000]byte // 0x4000 - 0x7FFF
	VRAM [0x2000]byte  // 0x8000 - 0x9FFF
	ERAM *[0x2000]byte // 0xA000 - 0xBFFF
	WRAM [0x2000]byte  // 0xC000 - 0xDFFF
	OAM  [0xA0]byte    // 0xFE00 - 0xFE9F
	IO   [0x80]byte    // 0xFF00 - 0xFF7F
	HRAM [0x7F]byte    // 0xFF80 - 0xFFFE
	IE   byte          // 0xFFFF
	// DIV Register
	// Hack para resetear timer.internalCounter cuando CPU escribe en el registro
	ResetDIV bool
	// DMA
	DMAIsActive      bool
	enableDMA        bool
	dmaSource        uint16
	dmaIndex         uint16
	dmaCyclesLeft    byte
	dmaDelay         byte    // ciclos de retardo inicial (2)
	pendingDMASource *uint16 // nuevo origen DMA si hay reinicio
	Client           byte
}

func (b *Bus) Read(addr uint16) byte {
	if !b.isAccessible(addr) {
		log.Printf("Acceso denegado en lectura para el cliente %d en %04X\n", b.Client, addr)
		return 0xFF
	}
	switch {
	case addr < 0x100 && b.bootActive:
		return b.BootROM[addr]
	case addr < 0x8000 || (addr >= 0xA000 && addr < 0xC000):
		return b.cart.Memory.Read(addr)

	case addr >= 0x8000 && addr < 0xA000:
		return b.VRAM[addr-0x8000]

	case addr >= 0xC000 && addr < 0xE000:
		return b.WRAM[addr-0xC000]

	case addr >= 0xE000 && addr < 0xFE00:
		// Echo RAM (mirror of C000–DDFF)
		return b.WRAM[addr-0xE000]

	case addr >= 0xFE00 && addr < 0xFEA0:
		return b.OAM[addr-0xFE00]

	case addr >= 0xFEA0 && addr < 0xFF00:
		log.Printf("Intento de lectura en zona no usable en %04X, se retorna 0xFF por cliente %d\n", addr, b.Client)
		return 0xFF

	case addr >= 0xFF00 && addr < 0xFF80:
		// El registro IF los bits 5, 6 y 7 siempre deben leerse con 1 y no con 0
		if addr == 0xFF0F {
			return b.IO[addr-0xFF00] | 0xE0
		}
		return b.IO[addr-0xFF00]

	case addr >= 0xFF80 && addr < 0xFFFF:
		return b.HRAM[addr-0xFF80]

	case addr == 0xFFFF:
		return b.IE

	default:
		log.Printf("Intento de lectura fuera de rango en %04X por cliente %d\n", addr, b.Client)
		return 0xFF
	}
}

func (b *Bus) Write(addr uint16, value byte) {
	if !b.isAccessible(addr) {
		log.Printf("Acceso denegado en escritura para el cliente %d en %04X\n", b.Client, addr)
		return
	}
	switch {
	case addr < 0x8000 || (addr >= 0xA000 && addr < 0xC000):
		b.cart.Memory.Write(addr, value)
		return

	case addr >= 0x8000 && addr < 0xA000:
		b.VRAM[addr-0x8000] = value

	case addr >= 0xC000 && addr < 0xE000:
		b.WRAM[addr-0xC000] = value

	case addr >= 0xE000 && addr < 0xFE00:
		// Echo RAM (mirror of C000–DDFF)
		b.WRAM[addr-0xE000] = value

	case addr >= 0xFE00 && addr < 0xFEA0:
		b.OAM[addr-0xFE00] = value

	case addr >= 0xFEA0 && addr < 0xFF00:
		log.Printf("Intento de escritura en zona no usable en %04X: %02X por cliente %d\n", addr, value, b.Client)

	case addr >= 0xFF00 && addr < 0xFF80:
		// Si se intenta escribir en DIV se establece en 0
		//https://gbdev.io/pandocs/Timer_and_Divider_Registers.html#ff04--div-divider-register
		if b.Client == ClientCPU && addr == DIVRegister {
			b.IO[addr-0xFF00] = 0x00
			b.ResetDIV = true
			return
		}
		// Activa el DMA
		if addr == 0xFF46 {
			b.IO[addr-0xFF00] = value
			b.doDMATransfer(value)
			return
		}
		// Desactiva el Boot ROM
		if addr == 0xFF50 && value != 0 {
			b.bootActive = false // Desactiva Boot ROM
		}
		b.IO[addr-0xFF00] = value

	case addr >= 0xFF80 && addr < 0xFFFF:
		b.HRAM[addr-0xFF80] = value

	case addr == 0xFFFF:
		b.IE = value

	default:
		log.Printf("Intento de escritura fuera de rango en %04X: %02X por cliente %d\n", addr, value, b.Client)
	}
}
func (b *Bus) isAccessible(addr uint16) bool {
	switch b.Client {
	case ClientCPU:
		if b.DMAIsActive {
			if (addr >= 0x8000 && addr < 0xA000) || (addr >= 0xFE00 && addr < 0xFEA0) {
				log.Printf("Bloqueando acceso a %04X por cliente %d (DMA activo: %v)\n", addr, b.Client, b.DMAIsActive)
				return false
			}
		}
		return true
	case ClientPPU:
		// TODO: PPU accede solo a VRAM y OAM
		return true
	case ClientTimer:
		// TODO: Timer accede a registros de temporizador
		return true
	case ClientDMA:
		// DMA puede leer desde la fuente DMA (dmaSource..dmaSource+0x9F)
		// y escribir a OAM (0xFE00..0xFE9F)
		if addr == b.dmaSource+b.dmaIndex {
			return true
		}
		if addr >= 0xFE00 && addr < 0xFEA0 {
			return true
		}
		return false
	case ClientLiteBoy:
		return true
	default:
		return false
	}
}
