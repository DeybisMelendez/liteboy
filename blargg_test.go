package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
)

const maxCycles = 80_000_000 // Máximo de ciclos para cada test

var testROMs = map[string]string{
	"01-special":            "roms/blargg/cpu_instrs/individual/01-special.gb",
	"02-interrupts":         "roms/blargg/cpu_instrs/individual/02-interrupts.gb",
	"03-op sp,hl":           "roms/blargg/cpu_instrs/individual/03-op sp,hl.gb",
	"04-op r,imm":           "roms/blargg/cpu_instrs/individual/04-op r,imm.gb",
	"05-op rp":              "roms/blargg/cpu_instrs/individual/05-op rp.gb",
	"06-ld r,r":             "roms/blargg/cpu_instrs/individual/06-ld r,r.gb",
	"07-jr,jp,call,ret,rst": "roms/blargg/cpu_instrs/individual/07-jr,jp,call,ret,rst.gb",
	"08-misc instrs":        "roms/blargg/cpu_instrs/individual/08-misc instrs.gb",
	"09-op r,r":             "roms/blargg/cpu_instrs/individual/09-op r,r.gb",
	"10-bit ops":            "roms/blargg/cpu_instrs/individual/10-bit ops.gb",
	"11-op a,(hl)":          "roms/blargg/cpu_instrs/individual/11-op a,(hl).gb",
	"instr_timing":          "roms/blargg/instr_timing/instr_timing.gb",
	"interrupt_time":        "roms/blargg/interrupt_time/interrupt_time.gb",
}

func TestBlarggCPUInstrs(t *testing.T) {
	for name, path := range testROMs {
		t.Run(name, func(t *testing.T) {
			if ok := runTestROM(path); !ok {
				t.Errorf("Test %s failed", name)
			}
		})
	}
}

func runTestROM(path string) bool {
	cart := cartridge.NewCartridge(path)
	gameBus := bus.NewBus(cart)
	gameCPU := cpu.NewCPU(gameBus)
	gamePPU := ppu.NewPPU(gameBus)

	cycles := 0
	for cycles < maxCycles {
		c := gameCPU.Step()
		gamePPU.Step(c)
		cycles += c
	}
	// Inspecciona el texto en pantalla (desde VRAM)
	text := extractScreenText(gameBus)
	if strings.Contains(text, "Passed") {
		return true
	}
	if strings.Contains(text, "Failed") {
		return false
	}
	fmt.Println("Una prueba falló porque le falto ciclos")
	return false // Timeout
}

// Lee los primeros caracteres del tile map en VRAM para detectar texto
func extractScreenText(b *bus.Bus) string {
	vram := b.VRAM
	base := 0x1800 // VRAM offset para $9800
	var out strings.Builder

	for i := 0; i < 32*32; i++ { // Lee ~4 filas de 32 tiles
		if i+base >= len(vram) {
			break
		}
		tile := vram[base+i]
		// Blargg imprime caracteres ASCII (espacio=0x20)
		if tile >= 0x20 && tile <= 0x7F {
			out.WriteByte(tile)
		} else {
			out.WriteByte('.') // Placeholder
		}
	}
	return out.String()
}
