package cpu

const (
	FlagZ byte = 0x80 // Zero
	FlagN byte = 0x40 // Subtract
	FlagH byte = 0x20 // Half Carry
	FlagC byte = 0x10 // Carry
)
