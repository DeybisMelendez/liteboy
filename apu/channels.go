package apu

import (
	"math"
	"sync"

	"github.com/deybismelendez/liteboy/bus"
)

type Channel struct {
	enabled       bool
	frequency     float64
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
	phase       float64
	waveRAM     [16]byte // 32 muestras de 4 bits
	bus         *bus.Bus
}

func (c *WaveChannel) GetSample() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.enabled || c.volumeShift < 0 {
		return 0
	}

	// Avance de fase correcto escalando por 32 muestras
	delta := c.frequency * 32.0 / sampleRate
	c.phase += delta
	if c.phase >= 32 {
		c.phase = 0
	}

	idx := int(c.phase) // 0..31
	byteIdx := idx / 2  // 0..15
	isHigh := (idx % 2) == 0
	raw := c.waveRAM[byteIdx]
	var sampleValue byte
	if isHigh {
		sampleValue = (raw >> 4) & 0x0F
	} else {
		sampleValue = raw & 0x0F
	}

	// Aplicar volumen
	adjusted := sampleValue >> uint(c.volumeShift)
	// Normalizar a [-1..1] y luego a int16
	normalized := (float64(adjusted)/7.5 - 1.0)
	return int(normalized * 32767)
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

// GetSample genera una muestra de ruido PCM de 16 bits [-32767, +32767].
func (c *NoiseChannel) GetSample() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Si el canal no está encendido o el volumen es cero, no hay salida.
	if !c.enabled || c.volume == 0 {
		return 0
	}

	// Calcular cuántas muestras deben transcurrir entre actualizaciones del LFSR:
	// periodCPU = divisor * 2^(clockShift+1) ciclos de CPU (4.194304 MHz)
	// samplesPerStep = periodCPU / 4194304 * sampleRate
	div := c.divisorCode
	if div == 0 {
		div = 8
	} else {
		div *= 16
	}
	periodCPU := float64(div) * math.Pow(2, float64(c.clockShift+1))
	samplesPerStep := periodCPU * sampleRate / 4194304.0

	// Avanzar el timer y actualizar LFSR cuando toque
	c.timer -= 1.0
	if c.timer <= 0 {
		c.timer += samplesPerStep

		// Retroalimentación del LFSR: bit0 xor bit1
		bit0 := c.lfsr & 1
		bit1 := (c.lfsr >> 1) & 1
		feedback := bit0 ^ bit1

		// Desplazar y poner el nuevo bit en posición 14
		c.lfsr = (c.lfsr >> 1) | (feedback << 14)

		// Si está en modo 7-bit, también actualizar bit6 y enmascarar resto
		if c.widthMode {
			c.lfsr = (c.lfsr &^ (1 << 6)) | (feedback << 6)
			c.lfsr &= 0x7F // solo 7 bits efectivos
		}
	}

	// Extraer el bit 0 como salida de nivel
	var sample int
	if (c.lfsr & 1) == 0 {
		sample = int(c.volume * 32767)
	} else {
		sample = -int(c.volume * 32767)
	}
	return sample
}
