package main

import (
	"testing"

	"github.com/deybismelendez/liteboy/apu"
	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
	"github.com/deybismelendez/liteboy/timer"
)

var passValues []byte = []byte{3, 5, 8, 13, 21, 34}
var failValues []byte = []byte{0x42, 0x42, 0x42, 0x42, 0x42, 0x42}
var mooneyeAcceptance = map[string]string{
	"add_sp_e_timing":                 "roms/mooneye/acceptance/add_sp_e_timing.gb",
	"bits/mem_oam":                    "roms/mooneye/acceptance/bits/mem_oam.gb",
	"bits/reg_f":                      "roms/mooneye/acceptance/bits/reg_f.gb",
	"boot_regs-dmgABC":                "roms/mooneye/acceptance/boot_regs-dmgABC.gb",
	"call_cc_timing":                  "roms/mooneye/acceptance/call_cc_timing.gb",
	"call_cc_timing2":                 "roms/mooneye/acceptance/call_cc_timing2.gb",
	"call_timing":                     "roms/mooneye/acceptance/call_timing.gb",
	"call_timing2":                    "roms/mooneye/acceptance/call_timing2.gb",
	"ei_sequence":                     "roms/mooneye/acceptance/ei_sequence.gb",
	"ei_timing":                       "roms/mooneye/acceptance/ei_timing.gb",
	"halt_ime0_ei":                    "roms/mooneye/acceptance/halt_ime0_ei.gb",
	"halt_ime0_nointr_timing":         "roms/mooneye/acceptance/halt_ime0_nointr_timing.gb",
	"halt_ime1_timing":                "roms/mooneye/acceptance/halt_ime1_timing.gb",
	"if_ie_registers":                 "roms/mooneye/acceptance/if_ie_registers.gb",
	"instr/daa":                       "roms/mooneye/acceptance/instr/daa.gb",
	"interrupts/ie_push":              "roms/mooneye/acceptance/interrupts/ie_push.gb",
	"intr_timing":                     "roms/mooneye/acceptance/intr_timing.gb",
	"jp_cc_timing":                    "roms/mooneye/acceptance/jp_cc_timing.gb",
	"jp_timing":                       "roms/mooneye/acceptance/jp_timing.gb",
	"ld_hl_sp_e_timing":               "roms/mooneye/acceptance/ld_hl_sp_e_timing.gb",
	"oam_dma/basic":                   "roms/mooneye/acceptance/oam_dma/basic.gb",
	"oam_dma/reg_read":                "roms/mooneye/acceptance/oam_dma/reg_read.gb",
	"oam_dma_restart":                 "roms/mooneye/acceptance/oam_dma_restart.gb",
	"oam_dma_start":                   "roms/mooneye/acceptance/oam_dma_start.gb",
	"oam_dma_timing":                  "roms/mooneye/acceptance/oam_dma_timing.gb",
	"pop_timing":                      "roms/mooneye/acceptance/pop_timing.gb",
	"ppu/intr_2_0_timing":             "roms/mooneye/acceptance/ppu/intr_2_0_timing.gb",
	"ppu/intr_2_mode0_timing_sprites": "roms/mooneye/acceptance/ppu/intr_2_mode0_timing_sprites.gb",
	"ppu/intr_2_mode0_timing":         "roms/mooneye/acceptance/ppu/intr_2_mode0_timing.gb",
	"ppu/intr_2_mode3_timing":         "roms/mooneye/acceptance/ppu/intr_2_mode3_timing.gb",
	"ppu/intr_2_oam_ok_timing":        "roms/mooneye/acceptance/ppu/intr_2_oam_ok_timing.gb",
	"ppu/stat_irq_blocking":           "roms/mooneye/acceptance/ppu/stat_irq_blocking.gb",
	"ppu/stat_lyc_onoff":              "roms/mooneye/acceptance/ppu/stat_lyc_onoff.gb",
	"push_timing":                     "roms/mooneye/acceptance/push_timing.gb",
	"rapid_di_ei":                     "roms/mooneye/acceptance/rapid_di_ei.gb",
	"reti_intr_timing":                "roms/mooneye/acceptance/reti_intr_timing.gb",
	"reti_timing":                     "roms/mooneye/acceptance/reti_timing.gb",
	"ret_cc_timing":                   "roms/mooneye/acceptance/ret_cc_timing.gb",
	"ret_timing":                      "roms/mooneye/acceptance/ret_timing.gb",
	"rst_timing":                      "roms/mooneye/acceptance/rst_timing.gb",
	"timer/div_write":                 "roms/mooneye/acceptance/timer/div_write.gb",
	"timer/rapid_toggle":              "roms/mooneye/acceptance/timer/rapid_toggle.gb",
	"timer/tim00_div_trigger":         "roms/mooneye/acceptance/timer/tim00_div_trigger.gb",
	"timer/tim00":                     "roms/mooneye/acceptance/timer/tim00.gb",
	"timer/tim01_div_trigger":         "roms/mooneye/acceptance/timer/tim01_div_trigger.gb",
	"timer/tim01":                     "roms/mooneye/acceptance/timer/tim01.gb",
	"timer/tim10_div_trigger":         "roms/mooneye/acceptance/timer/tim10_div_trigger.gb",
	"timer/tim10":                     "roms/mooneye/acceptance/timer/tim10.gb",
	"timer/tim11_div_trigger":         "roms/mooneye/acceptance/timer/tim11_div_trigger.gb",
	"timer/tim11":                     "roms/mooneye/acceptance/timer/tim11.gb",
	"timer/tima_reload":               "roms/mooneye/acceptance/timer/tima_reload.gb",
	"timer/tima_write_reloading":      "roms/mooneye/acceptance/timer/tima_write_reloading.gb",
	"timer/tma_write_reloading":       "roms/mooneye/acceptance/timer/tma_write_reloading.gb",
}

func TestMooneyeAcceptance(t *testing.T) {
	for name, path := range mooneyeAcceptance {
		t.Run(name, func(t *testing.T) {
			if ok := runMooneyeTestROM(path); !ok {
				t.Errorf("Mooneye test %s failed", name)
			}
		})
	}
}

// Verifica registros después de muchos ciclos buscando los valores esperados
func runMooneyeTestROM(path string) bool {
	cart := cartridge.NewCartridge(path)
	gameBus := bus.NewBus(cart)
	gamePPU := ppu.NewPPU(gameBus)
	gameTimer := timer.NewTimer(gameBus)
	gameAPU := apu.NewAPU(gameBus)
	gameCPU := cpu.NewCPU(gameBus, gameTimer, gamePPU, gameAPU)

	for range 1_000_000 {
		opcode := gameCPU.GetOpcode()
		//c := gameCPU.Step()
		//gamePPU.Step(c)
		gameCPU.Step()

		// Si ejecuta LD B, B (0x40), revisamos los registros
		if opcode == 0x40 {
			regs := gameCPU.GetRegisters()

			if equalBytes(regs, passValues) {
				return true
			} else if equalBytes(regs, failValues) {
				return false
			}
		}
	}
	return false // Timeout o sin detectar resultado
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
