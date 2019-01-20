package lib

import (
	"math"
	"time"
)

func PlayGap(duration time.Duration, volts chan Volt) {
	numTicks := *RATE * duration.Seconds()
	for t := 0; t < int(numTicks); t++ {
		volts <- Volt(0.0)
	}
}

// Boop writes voltages in range [-1.0, +1.0] to the channel volts.
//
func PlayTone(toneBegin, toneEnd float64, ramp RampType, duration time.Duration, volts chan Volt) {
	// Add global frequency base to tones to get absolute Hz.
	hz1 := *FREQ + toneBegin
	hz2 := *FREQ + toneEnd

	// Determine number of ticks and lenghts of ramps up & down.
	numTicks := *RATE * duration.Seconds()
	rampTicks := *RATE * *RAMP
	if rampTicks*2 > numTicks {
		// Ramp faster when tone is too short for usual ramp.
		rampTicks = numTicks / 2
	}

	// Output (to channel volts) one Volt (sound amplitude) per tick.
	for t := 0; t < int(numTicks); t++ {
		// Portion ranges 0.0 to almost 1.0.
		portion := float64(t) / float64(numTicks)
		// Interpolate part of the way between hz1 and hz2.
		hz := hz1 + portion*(hz2-hz1)
		// log.Printf("%06d: %8.0f hz (%5.1f, %5.1f)", t, hz, toneBegin, toneEnd)

		// Apply a raised-cosine envelope to the first and last RampTicks ticks.
		var envelopeGain float64
		switch {
		case (ramp == BEGIN || ramp == BOTH) && t < int(rampTicks):
			// First RampTicks, gain goes from 0.0 to 1.0
			{
				x := (float64(t) / rampTicks) * math.Pi
				y := math.Cos(x) // Cosine to shape envelope.
				envelopeGain = 0.5 - y/2.0
			}
		case (ramp == END || ramp == BOTH) && int(numTicks)-t < int(rampTicks):
			// Last RampTicks, gain goes from 1.0 to 0.0.
			{
				x := ((numTicks - float64(t)) / rampTicks) * math.Pi
				y := math.Cos(x) // Cosine to shape envelope.
				envelopeGain = 0.5 - y/2.0
			}
		default:
			// Middle of the Boop has full envelopeGain 1.0.
			{
				envelopeGain = 1.0
			}
		}

		// TODO: Observe old theta to eliminate exaggerated chirps.

		// The angle theta depends on the ticks and the frequency hz.
		theta := float64(t) * hz * (2.0 * math.Pi) / *RATE
		// Take the sin of the angle, and multiply by the envelopeGain.
		v := envelopeGain * math.Sin(theta)
		volts <- Volt(v)
	}
}
