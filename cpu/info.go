package cpu

const (
	FlagZ byte = 1 << 7 // Zero
	FlagN byte = 1 << 6 // Subtract
	FlagH byte = 1 << 5 // Half Carry
	FlagC byte = 1 << 4 // Carry
)
