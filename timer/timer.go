package timer

import "github.com/deybismelendez/liteboy/bus"

const (
	DIVRegister  = 0xFF04
	TIMARegister = 0xFF05
	TMARegister  = 0xFF06
	TACRegister  = 0xFF07
	IFRegister   = 0xFF0F
)

type Timer struct {
	internalCounter uint16
	bus             *bus.Bus
}

func NewTimer(bus *bus.Bus) *Timer {
	return &Timer{bus: bus}
}

func (t *Timer) Step(tCycles int) {
	t.bus.Client = 2
	if t.bus.TACWrite {
		t.OnTACWrite(t.bus.TACOld, t.bus.Read(TACRegister))
		t.bus.TACWrite = false
	}
	if t.bus.TimerReloading {
		t.bus.Write(TIMARegister, t.bus.Read(TMARegister))
		ifReg := t.bus.Read(IFRegister)
		t.bus.Write(IFRegister, ifReg|0x04)
		t.bus.TimerReloading = false
	}

	// Manejo de reinicio de DIV
	if t.bus.ResetDIV {
		t.checkFallingEdge()
		t.internalCounter = 0
		t.bus.ResetDIV = false
	}

	oldCounter := t.internalCounter
	t.internalCounter += uint16(tCycles)

	t.updateDIV(oldCounter)
	t.updateTIMA(oldCounter)
}

func (t *Timer) updateDIV(oldCounter uint16) {
	oldDIV := byte(oldCounter >> 8)
	newDIV := byte(t.internalCounter >> 8)
	if oldDIV != newDIV {
		t.bus.Write(DIVRegister, newDIV)
	}
}

func (t *Timer) updateTIMA(oldCounter uint16) {
	tac := t.bus.Read(TACRegister)
	if tac&0x04 == 0 {
		return // Timer deshabilitado
	}

	bitIndex := getTimerBitIndex(tac)
	oldBit := (oldCounter >> bitIndex) & 1
	newBit := (t.internalCounter >> bitIndex) & 1

	if oldBit == 1 && newBit == 0 {
		t.incrementTIMA()
	}
}

func (t *Timer) checkFallingEdge() {
	tac := t.bus.Read(TACRegister)
	if tac&0x04 == 0 {
		return // Timer deshabilitado
	}

	bitIndex := getTimerBitIndex(tac)
	oldBit := (t.internalCounter >> bitIndex) & 1

	if oldBit == 1 {
		t.incrementTIMA()
	}
}

func (t *Timer) incrementTIMA() {
	tima := t.bus.Read(TIMARegister)
	if tima == 0xFF {
		t.bus.Write(TIMARegister, 0x00)
		t.bus.TimerReloading = true
	} else {
		t.bus.Write(TIMARegister, tima+1)
	}
}

func getTimerBitIndex(tac byte) uint {
	switch tac & 0x03 {
	case 0:
		return 9 // 1024 ciclos
	case 1:
		return 3 // 16 ciclos
	case 2:
		return 5 // 64 ciclos
	case 3:
		return 7 // 256 ciclos
	}
	return 0
}
func (t *Timer) OnTACWrite(oldTAC, newTAC byte) {
	oldEnabled := oldTAC & 0x04
	newEnabled := newTAC & 0x04

	oldBitIndex := getTimerBitIndex(oldTAC)
	newBitIndex := getTimerBitIndex(newTAC)

	oldBit := (t.internalCounter >> oldBitIndex) & 1
	newBit := (t.internalCounter >> newBitIndex) & 1

	if oldEnabled != 0 && newEnabled != 0 && oldBitIndex == newBitIndex {
		// Timer sigue encendido y misma frecuencia: no hacer nada especial
		return
	}

	// Detectar flanco de bajada como en el hardware
	if oldEnabled != 0 && oldBit == 1 && (newEnabled == 0 || newBit == 0) {
		t.incrementTIMA()
	}
}
