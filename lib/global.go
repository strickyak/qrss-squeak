package lib

import (
	"flag"
	"time"
)

var DIT = flag.Duration("dit", 1.0 * time.Second, "seconds per dit")
var RATE = flag.Float64("rate", 44100, "samples per second")
var GAIN = flag.Float64("gain", 1.0, "global audio gain")
var RAMP = flag.Float64("ramp", 0.1, "seconds per ramp up or down")
var FREQ = flag.Float64("freq", 1000, "global frequency offset")
var STEP = flag.Float64("step", 500, "step frequency")
var CLIP = flag.Bool("clip", false, "clip silently (instead of panicking)")

type DiDah byte // Just dits & dahs, like '.' '-' ' ' etc.
type Volt float64 // Normalize from -1.0 to +1.0 (like Sin() and Cos()).

type Emitter interface {
	Emit(chan Volt)
	Duration() time.Duration
	String() string
}

type RampType int

const (
	NONE RampType = iota
	BEGIN
	END
	BOTH
)
