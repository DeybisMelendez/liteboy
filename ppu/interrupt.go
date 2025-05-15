package ppu

const (
	InterruptVBlank byte = 0
	InterruptSTAT   byte = 1
	InterruptTimer  byte = 2
	InterruptSerial byte = 3
	InterruptJoypad byte = 4
)

func (ppu *PPU) requestInterrupt(interruptBit byte) {
	const IF = 0xFF0F // Interrupt Flag register
	current := ppu.bus.Read(IF)
	ppu.bus.Write(IF, current|(1<<interruptBit))
}
