package lib

import (
	"flag"
	"log"
	"time"
)

var DIT = flag.Duration("dit", 1.0*time.Second, "seconds per dit")
var RATE = flag.Float64("rate", 44100, "samples per second")
var GAIN = flag.Float64("gain", 1.0, "global audio gain")
var RAMP = flag.Float64("ramp", 0.1, "seconds per ramp up or down")
var FREQ = flag.Float64("freq", 1000, "global frequency offset")
var WIDTH = flag.Float64("width", 10, "bandwidth")
var STEP = flag.Float64("step", 500, "step frequency")
var CLIP = flag.Bool("clip", false, "clip silently (instead of panicking)")

type DiDah byte         // Just dits & dahs, like '.' '-' ' ' etc.
type DiDahSlice []DiDah // Just dits & dahs, like '.' '-' ' ' etc.
type Volt float64       // Normalize from -1.0 to +1.0 (like Sin() and Cos()).

func (o DiDahSlice) String() string {
	bb := make([]byte, len(o))
	for i, e := range o {
		bb[i] = byte(e)
	}
	return string(bb)
}

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

func Must(pred bool) {
	if !pred {
		log.Fatal("Must FAILED")
	}
}
