package ppu

const (
	InterruptVBlank = 0
	InterruptSTAT   = 1
	InterruptTimer  = 2
	InterruptSerial = 3
	InterruptJoypad = 4
)

func (ppu *PPU) requestInterrupt(interruptBit byte) {
	const IF = 0xFF0F // Interrupt Flag register
	current := ppu.bus.Read(IF)
	ppu.bus.Write(IF, current|(1<<interruptBit))
}
