package apu

import "encoding/binary"

const BufferSize = 2500

type Reader struct {
	ch1         *SquareChannel
	ch2         *SquareChannel
	ch3         *WaveChannel
	ch4         *NoiseChannel
	leftVolume  float64
	rightVolume float64
}

func (r *Reader) Read(p []byte) (int, error) {
	for i := 0; i < BufferSize; i += 4 {
		// Obtener las muestras individuales
		s1 := int32(r.ch1.GetSample())
		s2 := int32(r.ch2.GetSample())
		s3 := int32(r.ch3.GetSample())
		s4 := int32(r.ch4.GetSample())

		// Mezcla simple promedio
		mixed := (s1 + s2 + s3 + s4) / 4
		// Clipping y conversión: limitar a int16
		if mixed > 32767 {
			mixed = 32767
		} else if mixed < -32768 {
			mixed = -32768
		}

		// Convertir a uint16 para LittleEndian
		// int16 a uint16 requiere máscara de bits
		sample := int(mixed)
		// Escribir muestra estéreo (izquierda y derecha igual)
		binary.LittleEndian.PutUint16(p[i:], uint16((sample*int(r.leftVolume*10000))/10000))    // Left
		binary.LittleEndian.PutUint16(p[i+2:], uint16((sample*int(r.rightVolume*10000))/10000)) // Right
	}
	return BufferSize, nil
}
