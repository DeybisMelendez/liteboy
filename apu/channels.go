package apu

import (
	"math"
	"sync"
)

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
	if c.envelopeStep > 0 && c.envelopeTimer > 0 {
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

func (c *SquareChannel) GetSample() int {
	freqRatio := c.frequency / sampleRate
	var sample int = 0

	if c.enabled && c.volume > 0 {
		pos := math.Mod(c.phase, 1.0)
		if pos < c.dutyRatio {
			sample = int(float64(c.volume) * 32767)
		} else {
			sample = -int(float64(c.volume) * 32767)
		}
		c.phase += freqRatio
		if c.phase >= 1.0 {
			c.phase -= 1.0
		}
	}
	return sample
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

func (c *WaveChannel) GetSample() int {

	freqRatio := c.frequency / sampleRate
	var sample int = 0

	if c.enabled {
		index := c.wavePos % 32
		data := c.waveRAM[index/2]
		var waveSample byte
		if index%2 == 0 {
			waveSample = (data >> 4) & 0x0F
		} else {
			waveSample = data & 0x0F
		}

		if c.volumeShift == -1 {
			waveSample = 0
		} else {
			waveSample >>= c.volumeShift
		}

		// Escalar a [-32767, 32767] sin distorsión
		waveValue := int32(waveSample)
		sample = int((waveValue * 2 * 32767 / 15) - 32767)

		c.phase += freqRatio
		if c.phase >= 1.0 {
			c.phase -= 1.0
			c.wavePos = (c.wavePos + 1) % 32
		}
	}
	return sample
}

type NoiseChannel struct {
	Channel
	mu          sync.Mutex
	triggered   bool
	clockShift  int
	widthMode   bool
	divisorCode int
	lfsr        uint16
	timer       float64
}

func (c *NoiseChannel) GetSample() int {
	var sample int = 0
	// Divisores reales según Game Boy hardware
	divisors := []int{8, 16, 32, 48, 64, 80, 96, 112}
	div := 8
	if c.divisorCode >= 0 && c.divisorCode < len(divisors) {
		div = divisors[c.divisorCode]
	}
	if c.enabled && c.volume > 0 {
		freq := 524288.0 / float64(div<<uint(c.clockShift))
		if freq < 1 {
			freq = 1
		}

		c.timer -= 1
		if c.timer <= 0 {
			c.timer += sampleRate / freq

			if c.lfsr == 0 {
				c.lfsr = 0x7FFF
			}

			// LFSR feedback calculation
			bit := (c.lfsr ^ (c.lfsr >> 1)) & 1
			c.lfsr = (c.lfsr >> 1) | (bit << 14)

			if c.widthMode {
				// 7-bit mode: bit 6 also updated
				c.lfsr &= ^uint16(1 << 6)
				c.lfsr |= (bit << 6)
			}
		}

		if c.lfsr&1 == 0 {
			sample = int(c.volume * 32767)
			//sample = int16(c.volume * 32767 / 15) // normalizado
		} else {
			sample = -int(c.volume * 32767)
			//sample = -int16(c.volume * 32767 / 15)
		}
	}
	return sample
}
