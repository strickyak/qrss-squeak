package lib

import (
	"flag"
	"log"
	"time"
)

var DIT = flag.Float64("dit", 0.1, "duration per dit, in seconds")
var RATE = flag.Float64("rate", 44100, "samples per second to output, and for all internal Emitters to send.")
var GAIN = flag.Float64("gain", 1.0, "global audio gain")
var RAMP = flag.Float64("ramp", 0.01, "duration per ramp up or ramp down, in seconds")
var FREQ = flag.Float64("freq", 1000, "global frequency offset")
var BW = flag.Float64("bw", 10, "bandwidth to use, in Hertz")
var CLIP = flag.Bool("clip", false, "clip silently (instead of panicking)")

// DiDah is just dits, dahs, and spaces, represented as '.', '-', and ' '.
type DiDah byte

// DiDahSlice is a sequence of dits, dahs, and spaces.
type DiDahSlice []DiDah // Just dits & dahs, like '.' '-' ' ' etc.
// Volt is a floating point level normalized between -1.0 to +1.0 (like Sin() and Cos()).
type Volt float64

// Convert from a slice of DiDah to an ordinary string, just for debugging & logging.
func (o DiDahSlice) String() string {
	bb := make([]byte, len(o))
	for i, e := range o {
		bb[i] = byte(e)
	}
	return string(bb)
}

// Emitter is something that emits a sequence of Volt representing an audio signal.
// The generator for each mode (cw, dt, etc.) is an Emitter.
type Emitter interface {
	// Send the audio signal as a sequence of Volt on the given channel.
	// Don't close the channel; that's the caller's responsibility.
	Emit(chan Volt)
	// Return the duration of the signal, if it were to be emitted.
	// The number of samples in the signal divided by *RATE should give the Duration.
	Duration() time.Duration
	// Debug string describing the Emitter.
	String() string
}

// InfiniteDuration indicates an Emitter that never quits.
// It's actual value is a little over 100 years.
var InfiniteDuration = 100 * 366 * 24 * time.Hour

// RampType describes an envelope to shape the output sine waves
// so they don't have sharp beginning and ends.
// We can apply a "Raised Cosine" shape to the beginning or the end or both.
// The global flag --ramp specifies the duration of the ramp up or ramp down.
// TODO: Be able to glue together signals so the phase is consistent from one to the next.
type RampType int

const (
	// NONE means sharp beginning and ending of the tone.
	NONE RampType = iota
	// Ramp up the beginning of the tone.
	BEGIN
	// Ramp down the ending of the tone.
	END
	// Ramp up and down the beginning and ending of the tone.
	BOTH
)

// Must is a quick-and-dirty Assert for invarients.
func Must(pred bool) {
	if !pred {
		log.Fatal("FATAL: ASSERTION FAILED")
	}
}

// Convert seconds in float64 to time.Duration.
// Used for our command line flags which are float64 seconds.
func Secs(x float64) time.Duration {
	nanos := (x * float64(time.Second))
	return time.Duration(int64(nanos))
}
