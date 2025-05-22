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

	for i := 0; i < len(p); i += 2 {
		var sample int16 = 0

		if c.enabled && c.volume > 0 {
			pos := math.Mod(c.phase, 1.0) // ciclo entre 0 y 1
			if pos < c.dutyRatio {
				sample = int16(float64(c.volume) * 32767)
			} else {
				sample = -int16(float64(c.volume) * 32767)
			}

			// Avanzar fase según frecuencia y sampleRate
			c.phase += freqRatio
			if c.phase >= 1.0 {
				c.phase -= 1.0
			}
		}

		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}
	return len(p) / 6, nil
}

type waveReader struct {
	channel *WaveChannel
}

func (r *waveReader) Read(p []byte) (int, error) {
	c := r.channel

	freqRatio := c.frequency / sampleRate

	for i := 0; i < len(p); i += 2 {
		var sample int16 = 0

		if c.enabled {
			// Calcular índice actual según la posición en la onda
			index := c.wavePos

			data := c.waveRAM[index/2]
			var waveSample byte
			if index%2 == 0 {
				waveSample = data >> 4
			} else {
				waveSample = data & 0x0F
			}

			// Aplicar volumen (0: mute, 1: 100%, 2: 50%, 3: 25%)
			switch c.volumeShift {
			case 0:
				waveSample = 0
			case 1:
				// 100%, sin cambio
			case 2:
				waveSample >>= 1
			case 3:
				waveSample >>= 2
			}

			// Ajustar sample para centrar en 0 y escalar
			sample = int16((int(waveSample) - 8) * 4096)

			// Avanzar fase y actualizar posición en la onda cuando fase complete ciclo
			c.phase += freqRatio
			if c.phase >= 1.0 {
				c.phase -= 1.0
				c.wavePos = (c.wavePos + 1) % 32
			}
		}

		binary.LittleEndian.PutUint16(p[i:], uint16(sample))
	}

	return len(p) / 6, nil
}

type noiseReader struct {
	channel *NoiseChannel
	lfsr    uint16
	timer   float64
}

func (r *noiseReader) Read(p []byte) (int, error) {
	c := r.channel

	// Precalcular divisores para evitar recreación cada ciclo
	divisors := []int{8, 16, 32, 48, 64, 80, 96, 112}
	div := 8 // valor por defecto
	if c.divisorCode >= 0 && c.divisorCode < len(divisors) {
		div = divisors[c.divisorCode]
	}

	for i := 0; i < len(p); i += 2 {
		var sample int16 = 0

		if c.enabled && c.volume > 0 {
			// Calcular frecuencia del ruido
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

				// XOR bit 0 y bit 1
				bit := (r.lfsr ^ (r.lfsr >> 1)) & 1
				r.lfsr = (r.lfsr >> 1) | (bit << 14)

				if c.widthMode {
					// modo 7 bits en LFSR (bit 6 se rellena con bit XOR)
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

	return len(p) / 6, nil
}
