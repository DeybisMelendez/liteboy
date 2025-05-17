package timer

import "github.com/deybismelendez/liteboy/bus"

const (
	DIVRegister  = 0xFF04
	TIMARegister = 0xFF05
	TMARegister  = 0xFF06
	TACRegister  = 0xFF07
)

type Timer struct {
	internalCounter uint16
	bus             *bus.Bus
}

func NewTimer(bus *bus.Bus) *Timer {
	return &Timer{
		bus: bus,
	}
}

func (timer *Timer) Step(tCycles int) {
	timer.bus.Client = 2
	// Hack para resetear DIV
	if timer.bus.ResetDIV {
		timer.internalCounter = 0
		timer.bus.ResetDIV = false
	}
	// Guardamos el estado anterior
	oldCounter := timer.internalCounter
	timer.internalCounter += uint16(tCycles)

	// --- DIV update (upper 8 bits of counter) ---
	oldDIV := byte(oldCounter >> 8)
	newDIV := byte(timer.internalCounter >> 8)
	if newDIV != oldDIV {
		timer.bus.Write(DIVRegister, newDIV)
	}

	// --- TIMA update ---
	tac := timer.bus.Read(TACRegister)
	if tac&0x04 == 0 {
		return // timer disabled
	}

	// Determine bit to watch based on TAC
	var bitIndex uint
	switch tac & 0x03 {
	case 0:
		bitIndex = 9 // 1024 cycles
	case 1:
		bitIndex = 3 // 16 cycles
	case 2:
		bitIndex = 5 // 64 cycles
	case 3:
		bitIndex = 7 // 256 cycles
	}

	oldBit := (oldCounter >> bitIndex) & 1
	newBit := (timer.internalCounter >> bitIndex) & 1

	// Detect flanco descendente: 1 -> 0
	if oldBit == 1 && newBit == 0 {
		tima := timer.bus.Read(TIMARegister)
		if tima == 0xFF {
			timer.bus.Write(TIMARegister, timer.bus.Read(TMARegister))
			ifReg := timer.bus.Read(0xFF0F)
			timer.bus.Write(0xFF0F, ifReg|0x04)
		} else {
			timer.bus.Write(TIMARegister, tima+1)
		}
	}
}
