package timer

import "github.com/deybismelendez/liteboy/bus"

const (
	DIVRegister  = 0xFF04
	TIMARegister = 0xFF05
	TMARegister  = 0xFF06
	TACRegister  = 0xFF07
)

type Timer struct {
	divCounter   uint16
	timerCounter int
	bus          *bus.Bus
}

func NewTimer(bus *bus.Bus) *Timer {
	return &Timer{
		bus: bus,
	}
}

func (timer *Timer) Step(tCycles int) {
	// --- DIV siempre avanza a 16384 Hz (cada 256 ciclos de CPU) ---
	timer.divCounter += uint16(tCycles)
	if timer.divCounter >= 256 {
		timer.divCounter -= 256
		div := timer.bus.Read(DIVRegister)
		timer.bus.Write(DIVRegister, div+1) // incrementar DIV
	}

	// --- TIMA controlado por TAC ---
	tac := timer.bus.Read(TACRegister)

	timerEnabled := tac&0x04 != 0

	if !timerEnabled {
		return
	}
	// Obtenemos el número de ciclos por incremento según TAC bits 1-0
	var threshold int
	switch tac & 0x03 {
	case 0:
		threshold = 1024 // 4096 Hz
	case 1:
		threshold = 16 // 262144 Hz
	case 2:
		threshold = 64 // 65536 Hz
	case 3:
		threshold = 256 // 16384 Hz
	}

	// Sumamos ciclos al timer interno
	timer.timerCounter += tCycles
	for timer.timerCounter >= threshold {
		timer.timerCounter -= threshold

		tima := timer.bus.Read(TIMARegister)
		if tima == 0xFF {
			// Desbordamiento: TIMA = TMA, IF |= 0x04
			timer.bus.Write(TIMARegister, timer.bus.Read(TMARegister))
			IFRegister := timer.bus.Read(0xFF0F)
			timer.bus.Write(0xFF0F, IFRegister|0x04)
		} else {
			// Solo incrementa
			timer.bus.Write(TIMARegister, tima+1)
		}
	}
}
