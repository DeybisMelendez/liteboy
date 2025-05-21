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
