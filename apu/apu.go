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
	audio  *audio.Context
	bus    *bus.Bus
	chan1  *SquareChannel
	chan2  *SquareChannel
	chan3  *WaveChannel
	chan4  *NoiseChannel
	player *audio.Player
}

func NewAPU(bus *bus.Bus) *APU {
	ctx := audio.NewContext(sampleRate)

	ch1 := &SquareChannel{}
	ch2 := &SquareChannel{}
	ch3 := &WaveChannel{bus: bus}
	ch4 := &NoiseChannel{}
	reader := &Reader{ch1: ch1, ch2: ch2, ch3: ch3, ch4: ch4}

	player, err := ctx.NewPlayer(reader)
	if err != nil {
		log.Fatal("error al crear audio player canal 1:", err)
	}
	// Inicializar waveform RAM con patrón 00 FF 00 FF ...
	for i := uint16(0); i < 0x10; i++ {
		var addr uint16 = 0xFF30 + i
		if i&2 == 0 {
			bus.Write(addr, 0x00)
		} else {
			bus.Write(addr, 0xFF)
		}
	}

	player.Play()

	return &APU{
		audio:  ctx,
		bus:    bus,
		chan1:  ch1,
		chan2:  ch2,
		chan3:  ch3,
		chan4:  ch4,
		player: player,
	}

}

func (apu *APU) Step() {
	apu.bus.Client = 4
	apu.updateChannel1()
	apu.updateChannel2()
	apu.updateChannel3()
	//apu.updateChannel4()
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
			nextFreq := c.shadowFreq
			if c.sweepDir < 0 {
				nextFreq -= change
			} else {
				nextFreq += change
			}
			if nextFreq > 2047 {
				c.enabled = false
			} else {
				c.shadowFreq = nextFreq
				c.frequency = 131072.0 / float64(2048-nextFreq)
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

	// Trigger
	if nr34&0x80 != 0 {
		c.enabled = (nr30 & 0x80) != 0
		c.triggered = true
		c.lengthTimer = 256 - int(nr31)
		// Cargar wave RAM
		for i := 0; i < 16; i++ {
			c.waveRAM[i] = apu.bus.Read(0xFF30 + uint16(i))
		}
		// Volume shift (NR32 bits 5-6)
		code := (nr32 >> 5) & 0x03
		switch code {
		case 0:
			c.volumeShift = -1 // mute
		case 1:
			c.volumeShift = 0 // 100%
		case 2:
			c.volumeShift = 1 // 50%
		case 3:
			c.volumeShift = 2 // 25%
		}

		// Frequency
		freq := uint16(nr33) | (uint16(nr34&0x07) << 8)
		c.frequency = 65536.0 / (2048.0 - float64(freq))

		// Load wave RAM
		for i := uint16(0); i < 16; i++ {
			c.waveRAM[i] = apu.bus.Read(0xFF30 + i)
		}
	}

	// Length timer
	if (nr34 & 0x40) != 0 {
		c.updateLengthTimer()
	}
}

// updateChannel4 debe inicializar y disparar el canal de ruido
func (apu *APU) updateChannel4() {
	c := apu.chan4
	c.mu.Lock()
	defer c.mu.Unlock()

	nr41 := apu.bus.Read(0xFF20)
	nr42 := apu.bus.Read(0xFF21)
	nr43 := apu.bus.Read(0xFF22)
	nr44 := apu.bus.Read(0xFF23)

	// Trigger (bit 7 de NR44)
	if nr44&0x80 != 0 {
		c.enabled = true
		c.triggered = true

		// Length timer (si NR44 bit 6 = 1, se usa en step())
		c.lengthTimer = 64 - int(nr41&0x3F)

		// Envelope (NR42)
		c.initialVolume = int(nr42 >> 4)
		c.volume = float64(c.initialVolume) / 15.0
		c.envelopeDir = 1
		if nr42&0x08 == 0 {
			c.envelopeDir = -1
		}
		c.envelopeStep = int(nr42 & 0x07)
		c.envelopeTimer = c.envelopeStep

		// Parámetros de ruido (NR43)
		c.clockShift = int(nr43 >> 4)
		c.widthMode = (nr43 & 0x08) != 0
		c.divisorCode = int(nr43 & 0x07)

		// Reiniciar LFSR
		c.lfsr = 0x7FFF // 15 bits todos a 1
		if c.widthMode {
			c.lfsr = 0x7F // 7 bits todos a 1
		}

		// Reset del timer
		c.timer = 0
	}

	// Length timer automático si NR44 bit 6 = 1
	if nr44&0x40 != 0 {
		c.updateLengthTimer()
	}

	// Envelope
	c.updateEnvelope()
}
