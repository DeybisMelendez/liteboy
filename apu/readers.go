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

	// Cada frame estéreo son 4 bytes (2 canales x 2 bytes por muestra)
	n := 5000

	for i := 0; i < n; i += 4 {
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

		// Escribir el mismo sample para canal izquierdo y derecho (estéreo)
		binary.LittleEndian.PutUint16(p[i:], uint16(sample))   // Left
		binary.LittleEndian.PutUint16(p[i+2:], uint16(sample)) // Right
	}

	return n, nil
}

type waveReader struct {
	channel *WaveChannel
}

func (r *waveReader) Read(p []byte) (int, error) {
	c := r.channel
	freqRatio := c.frequency / sampleRate

	n := 5000

	for i := 0; i < n; i += 4 {
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

			if c.volumeShift == -1 {
				waveSample = 0
			} else {
				waveSample >>= c.volumeShift
			}

			sample = int16((int(waveSample) - 8) * 4096)

			// Avanzar la posición de onda
			c.phase += freqRatio
			if c.phase >= 1.0 {
				c.phase -= 1.0
				c.wavePos = (c.wavePos + 1) % 32
			}
		}

		// Escribir la muestra a ambos canales (L y R)
		binary.LittleEndian.PutUint16(p[i:], uint16(sample))   // Left
		binary.LittleEndian.PutUint16(p[i+2:], uint16(sample)) // Right
	}

	return n, nil
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

	n := 5000

	for i := 0; i < n; i += 4 {
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

		// Escribir la muestra a ambos canales
		binary.LittleEndian.PutUint16(p[i:], uint16(sample))   // Left
		binary.LittleEndian.PutUint16(p[i+2:], uint16(sample)) // Right
	}

	return n, nil
}
