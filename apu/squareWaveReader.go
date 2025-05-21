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
