package cartridge

import (
	"fmt"
	"log"
	"os"
)

type Cartridge struct {
	Path             string
	Entry            []byte
	Logo             []byte
	Title            string
	ManufacturerCode string
	CGBFlag          byte
	SGBFlag          byte
	CartridgeType    string
	NewLicense       string
	ROMSize          int
	RAMSize          int
	Destination      string
	OldLicense       string
	Version          byte
	Checksum         byte
	GlobalChecksum   uint16
	ROM              [][0x4000]byte
}

// NewCartridge loads and parses a Game Boy ROM cartridge
func NewCartridge(path string) *Cartridge {
	rom, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}
	if len(rom) < 0x150 {
		log.Fatal("ROM demasiado corta, inválida")
	}

	cart := &Cartridge{Path: path}
	bankCount := numBanksFromHeader(rom[0x0148])
	cart.ROM = make([][0x4000]byte, bankCount)

	for i := 0; i < bankCount; i++ {
		start := i * 0x4000
		end := min(start+0x4000, len(rom))
		copy(cart.ROM[i][:], rom[start:end])
	}

	cart.Entry = rom[0x0100:0x0104]
	cart.Logo = rom[0x0104:0x0134]
	cart.Title = string(rom[0x0134:0x0143])
	cart.ManufacturerCode = string(rom[0x013F:0x0143])
	cart.CGBFlag = rom[0x0143]
	cart.NewLicense = newLicCodes[string(rom[0x0144:0x0146])]
	cart.SGBFlag = rom[0x0146]
	cart.CartridgeType = cartridgeTypes[rom[0x0147]]
	cart.ROMSize = romSizes[rom[0x0148]]
	cart.RAMSize = ramSizes[rom[0x0149]]
	cart.Destination = destinationCodes[rom[0x014A]]
	cart.OldLicense = oldLicCodes[rom[0x014B]]
	cart.Version = rom[0x014C]
	cart.Checksum = rom[0x014D]
	cart.GlobalChecksum = uint16(rom[0x014E])<<8 | uint16(rom[0x014F])

	return cart
}

func (c *Cartridge) GetROM() *[][0x4000]byte {
	return &c.ROM
}

func (c *Cartridge) PrintHeaderInfo() {
	fmt.Println("--- Información del cartucho: ---")
	fmt.Println("Título:", c.Title)
	fmt.Println("Código de manufactura:", c.ManufacturerCode)
	fmt.Println("Nueva Licencia:", c.NewLicense)
	fmt.Println("Tipo de cartucho:", c.CartridgeType)
	fmt.Println("Tamaño ROM:", c.ROMSize, "bytes")
	fmt.Println("Tamaño RAM:", c.RAMSize, "bytes")
	fmt.Println("Destino:", c.Destination)
	fmt.Println("Licencia antigua:", c.OldLicense)
	fmt.Println("Versión:", c.Version)
	fmt.Println("Checksum:", c.Checksum, "-", c.ValidateChecksum())
	fmt.Println("Global Checksum:", c.GlobalChecksum)
}

func (c *Cartridge) ValidateChecksum() string {
	var checksum byte = 0
	for addr := 0x0134; addr <= 0x014C; addr++ {
		checksum = checksum - c.ROM[0][addr] - 1
	}
	if c.Checksum == checksum {
		return "Válido"
	}
	return "No válido"
}

func numBanksFromHeader(headerByte byte) int {
	switch headerByte {
	case 0x00:
		return 2 // 32KB
	case 0x01:
		return 4
	case 0x02:
		return 8
	case 0x03:
		return 16
	case 0x04:
		return 32
	case 0x05:
		return 64
	case 0x06:
		return 128
	case 0x07:
		return 256
	case 0x08:
		return 512
	case 0x52:
		return 72
	case 0x53:
		return 80
	case 0x54:
		return 96
	default:
		return 2
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
