package apu

import (
	"encoding/binary"
	"math"
)

type squareWaveReader struct {
	channel *SquareChannel
}

func (r *squareWaveReader) Read(p []byte) (int, error) {
	c := r.channel
	freqRatio := c.frequency / sampleRate

	for i := 0; i < 5000; i += 2 {
		var sample int16 = 0

		if c.enabled && c.volume > 0 {
			pos := math.Mod(c.phase, 1.0)
			if pos < c.dutyRatio {
				sample = int16(float64(c.volume) * 32767)
			} else {
				sample = -int16(float64(c.volume) * 32767)
			}
			c.phase += freqRatio
			if c.phase >= 1.0 {
				c.phase -= 1.0
			}
		}

		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}

	return 5000, nil
}

type waveReader struct {
	channel *WaveChannel
}

func (r *waveReader) Read(p []byte) (int, error) {
	c := r.channel
	freqRatio := c.frequency / sampleRate

	for i := 0; i < 5000; i += 2 {
		var sample int16 = 0

		if c.enabled {
			index := c.wavePos % 32
			data := c.waveRAM[index/2]
			var waveSample byte
			if index%2 == 0 {
				waveSample = data >> 4
			} else {
				waveSample = data & 0x0F
			}

			// Volume adjustment
			if c.volumeShift == -1 {
				waveSample = 0
			} else {
				waveSample >>= c.volumeShift
			}
			sample = int16((int(waveSample) - 8) * 4096)

			// Advance wave position based on frequency
			c.phase += freqRatio
			if c.phase >= 1.0 {
				c.phase -= 1.0
				c.wavePos = (c.wavePos + 1) % 32
			}
		}

		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}

	return 5000, nil
}

type noiseReader struct {
	channel *NoiseChannel
	lfsr    uint16
	timer   float64
}

func (r *noiseReader) Read(p []byte) (int, error) {
	c := r.channel

	divisors := []int{8, 16, 32, 48, 64, 80, 96, 112}
	div := 8
	if c.divisorCode >= 0 && c.divisorCode < len(divisors) {
		div = divisors[c.divisorCode]
	}

	for i := 0; i < 5000; i += 2 {
		var sample int16 = 0

		if c.enabled && c.volume > 0 {
			freq := 524288.0 / float64(div<<c.clockShift)
			if freq <= 0 {
				freq = 1
			}

			r.timer -= 1
			if r.timer <= 0 {
				r.timer += sampleRate / freq

				if r.lfsr == 0 {
					r.lfsr = 0x7FFF
				}

				bit := (r.lfsr ^ (r.lfsr >> 1)) & 1
				r.lfsr = (r.lfsr >> 1) | (bit << 14)

				if c.widthMode {
					r.lfsr &= ^uint16(1 << 6)
					r.lfsr |= uint16(bit << 6)
				}
			}

			if r.lfsr&1 == 0 {
				sample = int16(float64(c.volume) * 32767)
			} else {
				sample = -int16(float64(c.volume) * 32767)
			}
		}

		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}

	return 5000, nil
}
