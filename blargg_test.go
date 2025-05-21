package main

import (
	"strings"
	"testing"

	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
	"github.com/deybismelendez/liteboy/timer"
)

var cpu_instrs = map[string]string{
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
}
var instr_timing = map[string]string{
	"instr_timing": "roms/blargg/instr_timing/instr_timing.gb",
}

var mem_timing = map[string]string{
	"01-read_timing":   "roms/blargg/mem_timing/individual/01-read_timing.gb",
	"02-write_timing":  "roms/blargg/mem_timing/individual/02-write_timing.gb",
	"03-modify_timing": "roms/blargg/mem_timing/individual/03-modify_timing.gb",
}
var mem_timing_2 = map[string]string{
	"01-read_timing":   "roms/blargg/mem_timing-2/rom_singles/01-read_timing.gb",
	"02-write_timing":  "roms/blargg/mem_timing-2/rom_singles/02-write_timing.gb",
	"03-modify_timing": "roms/blargg/mem_timing-2/rom_singles/03-modify_timing.gb",
}

func TestBlargg_cpu_instrs(t *testing.T) {
	for name, path := range cpu_instrs {
		t.Run(name, func(t *testing.T) {
			if ok := runTestROM(path); !ok {
				t.Errorf("Test %s failed", name)
			}
		})
	}
}

func TestBlargg_instr_timing(t *testing.T) {
	for name, path := range instr_timing {
		t.Run(name, func(t *testing.T) {
			if ok := runTestROM(path); !ok {
				t.Errorf("Test %s failed", name)
			}
		})
	}
}
func TestBlargg_mem_timing(t *testing.T) {
	for name, path := range mem_timing {
		t.Run(name, func(t *testing.T) {
			if ok := runTestROM(path); !ok {
				t.Errorf("Test %s failed", name)
			}
		})
	}
}
func TestBlargg_mem_timing_2(t *testing.T) {
	for name, path := range mem_timing_2 {
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
	gamePPU := ppu.NewPPU(gameBus)
	gameTimer := timer.NewTimer(gameBus)
	gameCPU := cpu.NewCPU(gameBus, gameTimer, gamePPU)

	for range 20 {
		for range 400_000 {
			gameCPU.Step()
		}
		// Inspecciona el texto en pantalla (desde VRAM)
		text := extractScreenText(gameBus)
		if strings.Contains(text, "Passed") {
			return true
		}
	}
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

/*var oam_bug = map[string]string{
	"1-lcd_sync":        "roms/blargg/oam_bug/rom_singles/1-lcd_sync.gb",
	"2-causes":          "roms/blargg/oam_bug/rom_singles/2-causes.gb",
	"3-non_causes":      "roms/blargg/oam_bug/rom_singles/3-non_causes.gb",
	"4-scanline_timing": "roms/blargg/oam_bug/rom_singles/4-scanline_timing.gb",
	"5-timing_bug":      "roms/blargg/oam_bug/rom_singles/5-timing_bug.gb",
	"6-timing_no_bug":   "roms/blargg/oam_bug/rom_singles/6-timing_no_bug.gb",
	"7-timing_effect":   "roms/blargg/oam_bug/rom_singles/7-timing_effect.gb",
	"8-instr_effect":    "roms/blargg/oam_bug/rom_singles/8-instr_effect.gb",
}*/

/*func TestBlargg_oam_bug(t *testing.T) {
	for name, path := range oam_bug {
		t.Run(name, func(t *testing.T) {
			if ok := runTestROM(path); !ok {
				t.Errorf("Test %s failed", name)
			}
		})
	}
}*/
