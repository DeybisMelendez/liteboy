package apu

import (
	"math"
	"sync"

	"github.com/deybismelendez/liteboy/bus"
)

var noiseDivisors = [8]float64{8, 16, 32, 48, 64, 80, 96, 112}

type Channel struct {
	enabled       bool
	frequency     float64
	lengthTimer   int
	envelopeStep  int
	envelopeTimer int
	envelopeDir   int
	initialVolume int
	volume        float64
	currentVolume int
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
	if c.envelopeStep == 0 {
		// Si envelopeStep es 0, el envelope no hace nada
		return
	}
	if c.envelopeTimer > 0 {
		c.envelopeTimer--
	}
	if c.envelopeTimer == 0 {
		c.envelopeTimer = c.envelopeStep
		newVolume := c.currentVolume + c.envelopeDir
		if newVolume >= 0 && newVolume <= 15 {
			c.currentVolume = newVolume
			c.volume = float64(c.currentVolume) / 15.0
		} else {
			// Si se pasa del rango, el envelope deja de cambiar (comportamiento real)
			// En hardware GB el envelope deja de cambiar al llegar a 0 o 15
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
	lfsr        uint16
	phase       float64
	divisorCode int
	shift       int
	widthMode   int
}

// GetSample genera una muestra de ruido PCM de 16 bits [-32767, +32767].
func (c *NoiseChannel) GetSample() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.enabled || c.volume == 0 {
		return 0
	}

	// Tabla de divisores según la documentación
	divisorTable := [8]float64{0.5, 1, 2, 3, 4, 5, 6, 7}
	divisor := divisorTable[c.divisorCode]
	frequency := 262144.0 / (divisor * math.Pow(2, float64(c.shift)))

	phaseIncrement := frequency / float64(sampleRate)

	// Actualizar fase y clock del LFSR
	c.phase += phaseIncrement
	numClocks := int(c.phase)
	if numClocks > 0 {
		c.phase -= float64(numClocks)
		for i := 0; i < numClocks; i++ {
			feedback := (c.lfsr & 1) ^ ((c.lfsr >> 1) & 1)
			c.lfsr = (c.lfsr >> 1) | (feedback << 14)
			if c.widthMode == 1 {
				// Modo 7 bits: establecer bit 6
				c.lfsr = (c.lfsr & 0xFFBF) | (feedback << 6)
			}
		}
	}

	// Generar muestra: bit 0 invertido, escalado por volumen
	var sample int
	if (c.lfsr & 1) == 0 {
		sample = int(c.volume * 32767)
	} else {
		sample = -int(c.volume * 32767)
	}
	return sample
}
