package cartridge

import (
	"fmt"
	"log"
	"os"
)

type Cartridge struct {
	path             string
	entry            []byte
	logo             []byte
	title            string
	manuFacturerCode string
	cgbFlag          byte
	sgbFlag          byte
	cartridgeType    string
	newLic           string
	romSize          int
	ramSize          int
	destination      string
	oldLic           string
	version          byte
	checksum         byte
	globalChecksum   uint16
	rom              []byte
}

func NewCartridge(path string) *Cartridge {
	cart := &Cartridge{}
	cart.path = path
	cart.load(path)
	return cart
}

func (cartridge *Cartridge) load(path string) {
	rom, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error al leer la ROM: %v", err)
	}
	cartridge.rom = rom
	cartridge.entry = rom[0x0100:0x0103]
	cartridge.logo = rom[0x0104:0x133]
	cartridge.title = string(rom[0x0134:0x0144])
	cartridge.manuFacturerCode = string(rom[0x013F:0x0142])
	cartridge.cgbFlag = rom[0x143]
	cartridge.newLic = newLicCodes[fmt.Sprintf("%02X", string(rom[0x0144:0x0145]))]
	cartridge.sgbFlag = rom[0x0146]
	cartridge.cartridgeType = cartridgeTypes[rom[0x0147]]
	cartridge.romSize = romSizes[rom[0x0148]]
	cartridge.ramSize = ramSizes[rom[0x0149]]
	cartridge.destination = destinationCodes[rom[0x014A]]
	cartridge.oldLic = oldLicCodes[rom[0x014B]]
	cartridge.version = rom[0x014C]
	cartridge.checksum = rom[0x014D]
	cartridge.globalChecksum = uint16(rom[0x14E])<<8 | uint16(rom[0x14F])
	fmt.Println("--- Cartucho cargado ---")
	fmt.Println("Titulo:", cartridge.title)
}

func (cartridge *Cartridge) PrintHeaderInfo() {
	fmt.Println("--- Información del cartucho: ---")
	fmt.Println("Titulo:", cartridge.title)
	fmt.Println("Código de Manufactura:", cartridge.manuFacturerCode)
	fmt.Println("Nueva Licencia:", cartridge.newLic)
	fmt.Println("Tipo de cartucho:", cartridge.cartridgeType)
	fmt.Println("Tamaño de la rom:", cartridge.romSize, "bytes")
	fmt.Println("Tamaño de la ram:", cartridge.ramSize, "bytes")
	fmt.Println("Destino:", cartridge.destination)
	fmt.Println("Licencia:", cartridge.oldLic)
	fmt.Println("version:", cartridge.version)
	fmt.Println("Checksum:", cartridge.checksum, cartridge.validateChecksum())
	fmt.Println("Global Checksum:", cartridge.globalChecksum)
}

func (cartridge *Cartridge) validateChecksum() string {
	var checksum byte = 0
	for addr := 0x0134; addr <= 0x014C; addr++ {
		checksum = checksum - cartridge.rom[addr] - 1
	}
	if cartridge.checksum == checksum {
		return "Valido"
	} else {
		return "No valido"
	}
}
