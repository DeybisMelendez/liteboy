package apu

type SquareChannel struct {
	frequency     float64
	volume        float64
	dutyRatio     float64
	enabled       bool
	lengthTimer   int
	envelopeStep  int
	envelopeTimer int
	envelopeDir   int
	initialVolume int
	sweepTime     int
	sweepCounter  int
	sweepShift    int
	sweepDir      int
	shadowFreq    uint16
	triggered     bool
	phase         float64
}
type WaveChannel struct {
	enabled     bool
	triggered   bool
	lengthTimer int
	volumeShift int
	frequency   float64
	phase       float64
	waveRAM     [32]byte // 32 muestras de 4 bits
	wavePos     int
}
