package apu

import (
	"log"

	"github.com/deybismelendez/liteboy/bus"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100
)

type APU struct {
	audio   *audio.Context
	bus     *bus.Bus
	chan1   *SquareChannel
	chan2   *SquareChannel
	chan3   *WaveChannel
	chan4   *NoiseChannel
	player1 *audio.Player
	player2 *audio.Player
	player3 *audio.Player
	player4 *audio.Player
	reader1 *squareWaveReader
	reader2 *squareWaveReader
	reader3 *waveReader
	reader4 *noiseReader
	ticks   int
}

func NewAPU(bus *bus.Bus) *APU {
	ctx := audio.NewContext(sampleRate)

	ch1 := &SquareChannel{}
	ch2 := &SquareChannel{}
	ch3 := &WaveChannel{}
	ch4 := &NoiseChannel{}
	reader1 := &squareWaveReader{channel: ch1}
	reader2 := &squareWaveReader{channel: ch2}
	reader3 := &waveReader{channel: ch3}
	reader4 := &noiseReader{channel: ch4}

	player1, err := ctx.NewPlayer(reader1)
	if err != nil {
		log.Fatal("error al crear audio player canal 1:", err)
	}
	player2, err := ctx.NewPlayer(reader2)
	if err != nil {
		log.Fatal("error al crear audio player canal 2:", err)
	}
	player3, err := ctx.NewPlayer(reader3)
	if err != nil {
		log.Fatal("error al crear audio player canal 3:", err)
	}
	player4, err := ctx.NewPlayer(reader4)
	if err != nil {
		log.Fatal("error al crear audio player canal 4:", err)
	}
	player1.Play()
	player2.Play()
	player3.Play()
	player4.Play()

	return &APU{
		audio:   ctx,
		bus:     bus,
		chan1:   ch1,
		chan2:   ch2,
		chan3:   ch3,
		chan4:   ch4,
		player1: player1,
		player2: player2,
		player3: player3,
		player4: player4,
		reader1: reader1,
		reader2: reader2,
		reader3: reader3,
		reader4: reader4,
	}

}

func (apu *APU) Step() {
	if apu.ticks >= 24 {
		apu.ticks -= 24
		apu.updateChannel1()
		apu.updateChannel2()
		apu.updateChannel3()
		//apu.updateChannel4()
	}
	apu.ticks++
}

func (apu *APU) updateChannel1() {
	c := apu.chan1
	c.mu.Lock()
	defer c.mu.Unlock()
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

	c.updateEnvelope()

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

	c.updateLengthTimer()
}
func (apu *APU) updateChannel2() {
	c := apu.chan2
	c.mu.Lock()
	defer c.mu.Unlock()
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
		c.frequency = 131072.0 / (2048 - float64(freq))

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

	c.updateEnvelope()

	c.updateLengthTimer()
}
func (apu *APU) updateChannel3() {
	c := apu.chan3
	c.mu.Lock()
	defer c.mu.Unlock()

	nr30 := apu.bus.Read(0xFF1A)
	nr31 := apu.bus.Read(0xFF1B)
	nr32 := apu.bus.Read(0xFF1C)
	nr33 := apu.bus.Read(0xFF1D)
	nr34 := apu.bus.Read(0xFF1E)

	if nr30&0x80 == 0 {
		c.enabled = false
		return
	}

	if nr34&0x80 != 0 {
		c.enabled = true
		c.triggered = true
		c.lengthTimer = 256 - int(nr31)
		c.wavePos = 1
		c.phase = 0.0

		// Cargar wave RAM
		for i := 0; i < 16; i++ {
			c.waveRAM[i] = apu.bus.Read(0xFF30 + uint16(i))
		}

		switch (nr32 >> 5) & 0x03 {
		case 0:
			c.volumeShift = -1 // Silencio
		case 1:
			c.volumeShift = 0 // 100%
		case 2:
			c.volumeShift = 1 // 50%
		case 3:
			c.volumeShift = 2 // 25%
		}
	}

	// Frecuencia: 131072 / (2048 - freq)
	freq := uint16(nr33) | (uint16(nr34&0x07) << 8)
	if freq >= 2048 {
		freq = 2047
	}
	c.frequency = 2097152.0 / (2.0 * float64(2048-freq))
	//c.frequency = 131072.0 / float64(2048-freq)

	c.updateLengthTimer()
}

func (apu *APU) updateChannel4() {
	c := apu.chan4
	c.mu.Lock()
	defer c.mu.Unlock()
	nr41 := apu.bus.Read(0xFF20)
	nr42 := apu.bus.Read(0xFF21)
	nr43 := apu.bus.Read(0xFF22)
	nr44 := apu.bus.Read(0xFF23)

	if nr44&0x80 != 0 {
		c.enabled = true
		c.triggered = true
		c.lengthTimer = 64 - int(nr41&0x3F)
		c.initialVolume = int(nr42 >> 4)
		c.volume = float64(c.initialVolume) / 15.0
		c.envelopeDir = 1
		if nr42&0x08 == 0 {
			c.envelopeDir = -1
		}
		c.envelopeStep = int(nr42 & 0x07)
		c.envelopeTimer = c.envelopeStep
		c.clockShift = int(nr43 >> 4)
		c.widthMode = (nr43 & 0x08) != 0
		c.divisorCode = int(nr43 & 0x07)
	}

	c.updateEnvelope()
	c.updateLengthTimer()
}
