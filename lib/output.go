package lib

import (
	"bufio"
	"log"
)

const MaxShort = 0x7FFF

// It multiplies by an overall gain, converts to signed int16, and writes to the writer in little-endian format.
// When the input volts channel has no more, we write true to the done channel,
// so the main program can exit.
// Volts below -1.0 or above +1.0 get clipped hard.
func Output(volts chan Volt, w *bufio.Writer, done chan bool) {
	for {
		volt, ok := <-volts
		if !ok {
			break
		}

		y := *GAIN * float64(volt)
		// Clip hard at +/- 1 unit.
		if y > 1.0 {
			if !*CLIP {
				log.Panicf("Clipping not allowed without --clip flag: y=%g", y)
			}
			y = 1.0
		}
		if y < -1.0 {
			if !*CLIP {
				log.Panicf("Clipping not allowed without --clip flag: y=%g", y)
			}
			y = -1.0
		}
		yShort := int(MaxShort * y)
		w.WriteByte(byte(yShort))
		w.WriteByte(byte(yShort >> 8))
	}
	done <- true
}
