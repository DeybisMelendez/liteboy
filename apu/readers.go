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
	for i := 0; i < len(p); i += 2 {
		var sample int16 = 0
		if c.enabled {
			pos := math.Mod(c.phase, 1.0)
			if pos < c.dutyRatio {
				sample = int16(c.volume * 32767)
			} else {
				sample = -int16(c.volume * 32767)
			}
			c.phase += c.frequency / sampleRate
			if c.phase >= 1.0 {
				c.phase -= 1.0
			}
		}
		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}
	return len(p), nil
}

type waveReader struct {
	channel *WaveChannel
}

func (r *waveReader) Read(p []byte) (int, error) {
	c := r.channel
	for i := 0; i < len(p); i += 2 {
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

			// Volume shift (0: mute, 1: 100%, 2: 50%, 3: 25%)
			switch c.volumeShift {
			case 0:
				waveSample = 0
			case 1:
				// 100%, no shift
			case 2:
				waveSample >>= 1
			case 3:
				waveSample >>= 2
			}

			sample = int16((int(waveSample) - 8) * 4096) // center at 0
			c.phase += c.frequency / sampleRate
			if c.phase >= 1.0 {
				c.phase -= 1.0
				c.wavePos = (c.wavePos + 1) % 32
			}
		}

		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}

	return len(p), nil
}

type noiseReader struct {
	channel *NoiseChannel
	lfsr    uint16
}

func (r *noiseReader) Read(p []byte) (int, error) {
	c := r.channel
	for i := 0; i < len(p); i += 2 {
		var sample int16 = 0
		if c.enabled {
			if r.lfsr == 0 {
				r.lfsr = 0x7FFF // Inicializar si es cero
			}

			// LFSR
			bit := (r.lfsr ^ (r.lfsr >> 1)) & 1
			r.lfsr = (r.lfsr >> 1) | (bit << 14)

			if c.widthMode {
				// LFSR de 7 bits
				r.lfsr &= ^uint16(1 << 6)
				r.lfsr |= (bit << 6)
			}

			if r.lfsr&1 == 0 {
				sample = int16(c.volume * 32767)
			} else {
				sample = -int16(c.volume * 32767)
			}
		}

		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}
	return len(p), nil
}
