package apu

import "sync"

type Channel struct {
	enabled       bool
	lengthTimer   int
	envelopeStep  int
	envelopeTimer int
	envelopeDir   int
	initialVolume int
	volume        float64
}

func (c *Channel) updateLengthTimer() {
	if c.lengthTimer > 0 {
		c.lengthTimer--
		if c.lengthTimer == 0 {
			c.enabled = false
		}
	}
}
func (c *Channel) updateEnvelope() {
	if c.envelopeStep > 0 {
		c.envelopeTimer--
		if c.envelopeTimer <= 0 {
			c.envelopeTimer = c.envelopeStep
			newVolume := c.initialVolume + c.envelopeDir
			if newVolume >= 0 && newVolume <= 15 {
				c.initialVolume = newVolume
				c.volume = float64(c.initialVolume) / 15.0
			}
		}
	}
}

type SquareChannel struct {
	Channel
	mu           sync.Mutex
	frequency    float64
	dutyRatio    float64
	sweepTime    int
	sweepCounter int
	sweepShift   int
	sweepDir     int
	shadowFreq   uint16
	triggered    bool
	phase        float64
}
type WaveChannel struct {
	Channel
	mu          sync.Mutex
	triggered   bool
	volumeShift int
	frequency   float64
	phase       float64
	waveRAM     [32]byte // 32 muestras de 4 bits
	wavePos     int
}

type NoiseChannel struct {
	Channel
	mu          sync.Mutex
	triggered   bool
	clockShift  int
	widthMode   bool
	divisorCode int
}
