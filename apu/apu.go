package apu

import (
	"log"

	"github.com/deybismelendez/liteboy/bus"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate      = 44100
	frameRate       = 60
	samplesPerFrame = sampleRate / frameRate
)

type APU struct {
	audio   *audio.Context
	bus     *bus.Bus
	chan1   *SquareChannel
	chan2   *SquareChannel
	player1 *audio.Player
	player2 *audio.Player
	reader1 *squareWaveReader
	reader2 *squareWaveReader
}

func NewAPU(bus *bus.Bus) *APU {
	ctx := audio.NewContext(sampleRate)

	ch1 := &SquareChannel{}
	ch2 := &SquareChannel{}
	reader1 := &squareWaveReader{channel: ch1}
	reader2 := &squareWaveReader{channel: ch2}

	player1, err := audio.NewPlayer(ctx, reader1)
	if err != nil {
		log.Fatal("error al crear audio player canal 1:", err)
	}
	player2, err := audio.NewPlayer(ctx, reader2)
	if err != nil {
		log.Fatal("error al crear audio player canal 2:", err)
	}

	player1.Play()
	player2.Play()

	return &APU{
		audio:   ctx,
		bus:     bus,
		chan1:   ch1,
		chan2:   ch2,
		player1: player1,
		player2: player2,
		reader1: reader1,
		reader2: reader2,
	}
}

func (apu *APU) Step() {
	apu.updateChannel1()
	apu.updateChannel2()
}

func (apu *APU) updateChannel1() {
	c := apu.chan1

	nr10 := apu.bus.Read(0xFF10)
	nr11 := apu.bus.Read(0xFF11)
	nr12 := apu.bus.Read(0xFF12)
	nr13 := apu.bus.Read(0xFF13)
	nr14 := apu.bus.Read(0xFF14)

	// Trigger
	if nr14&0x80 != 0 {
		c.enabled = true
		c.triggered = true
		c.lengthTimer = 64 - int(nr11&0x3F)
		c.initialVolume = int(nr12 >> 4)
		c.volume = float64(c.initialVolume) / 15.0
		c.envelopeDir = 1
		if nr12&0x08 == 0 {
			c.envelopeDir = -1
		}
		c.envelopeStep = int(nr12 & 0x07)
		c.envelopeTimer = c.envelopeStep
		c.sweepTime = int((nr10 >> 4) & 0x07)
		c.sweepDir = 1
		if nr10&0x08 != 0 {
			c.sweepDir = -1
		}
		c.sweepShift = int(nr10 & 0x07)
		c.shadowFreq = uint16(nr13) | (uint16(nr14&0x07) << 8)

		c.frequency = 131072.0 / float64(2048-c.shadowFreq)

		// Duty
		switch (nr11 >> 6) & 0x03 {
		case 0:
			c.dutyRatio = 0.125
		case 1:
			c.dutyRatio = 0.25
		case 2:
			c.dutyRatio = 0.5
		case 3:
			c.dutyRatio = 0.75
		}
	}

	// Envelope
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

	// Sweep
	if c.sweepTime > 0 {
		c.sweepCounter++
		if c.sweepCounter >= c.sweepTime {
			c.sweepCounter = 0
			change := c.shadowFreq >> uint16(c.sweepShift)
			if c.sweepDir < 0 {
				c.shadowFreq -= change
			} else {
				c.shadowFreq += change
			}
			if c.shadowFreq > 2047 {
				c.enabled = false
			} else {
				c.frequency = 131072.0 / float64(2048-c.shadowFreq)
			}
		}
	}

	// Length timer
	if c.lengthTimer > 0 {
		c.lengthTimer--
		if c.lengthTimer == 0 {
			c.enabled = false
		}
	}
}
func (apu *APU) updateChannel2() {
	c := apu.chan2

	nr21 := apu.bus.Read(0xFF16)
	nr22 := apu.bus.Read(0xFF17)
	nr23 := apu.bus.Read(0xFF18)
	nr24 := apu.bus.Read(0xFF19)

	// Trigger
	if nr24&0x80 != 0 {
		c.enabled = true
		c.triggered = true
		c.lengthTimer = 64 - int(nr21&0x3F)
		c.initialVolume = int(nr22 >> 4)
		c.volume = float64(c.initialVolume) / 15.0
		c.envelopeDir = 1
		if nr22&0x08 == 0 {
			c.envelopeDir = -1
		}
		c.envelopeStep = int(nr22 & 0x07)
		c.envelopeTimer = c.envelopeStep

		freq := uint16(nr23) | (uint16(nr24&0x07) << 8)
		c.frequency = 131072.0 / float64(2048-freq)

		// Duty
		switch (nr21 >> 6) & 0x03 {
		case 0:
			c.dutyRatio = 0.125
		case 1:
			c.dutyRatio = 0.25
		case 2:
			c.dutyRatio = 0.5
		case 3:
			c.dutyRatio = 0.75
		}
	}

	// Envelope
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

	// Length timer
	if c.lengthTimer > 0 {
		c.lengthTimer--
		if c.lengthTimer == 0 {
			c.enabled = false
		}
	}
}
