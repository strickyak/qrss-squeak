package lib

import (
	"fmt"
	"github.com/strickyak/qrss-squeak/font5x7"
	"log"
	"time"
)

type HellConf struct {
	Dit       time.Duration
	Freq      float64
	Bandwidth float64
	Morse     []DiDah
	Text      string
	Tail      bool
}
type HellEmitter struct {
	HellConf

	Total time.Duration
}

func NewHellEmitter(conf *HellConf) *HellEmitter {
	o := &HellEmitter{
		HellConf: *conf,
	}
	o.Total = o.Duration()
	return o
}

func (o *HellEmitter) Duration() time.Duration {
	var dur time.Duration
	for _ = range o.Text {
		dur += 6 * o.Dit
	}
	if !o.Tail {
		dur -= o.Dit
	}
	return dur
}

func (o *HellEmitter) String() string {
	return fmt.Sprintf("HellEmitter{text=%q,freq=%.1f,width=%.1f,dit=%v,total=%v}", o.Text, o.Freq, o.Bandwidth, o.Dit, o.Total)
}

func boolSliceToDiDahSlice(bools []bool) DiDahSlice {
	var z DiDahSlice
	for _, b := range bools {
		if b {
			z = append(z, '.')
		} else {
			z = append(z, ' ')
		}
	}
	return z
}

func (o *HellEmitter) Emit(out chan Volt) {
	log.Printf("HellEmitter Start: %v", o)

	var inputs []Emitter
	for i, row := range font5x7.EightHorizontalRowsOfBool(o.Text) {
		inputs = append(inputs, NewCWEmitter(&CWConf{
			Dit:   o.Dit,
			Freq:  o.Bandwidth * (1.0 - float64(i)/7),
			Morse: boolSliceToDiDahSlice(row),
			Tail:  o.Tail,
			NoGap: true, // TODO: get rid of this hack.
		}))
	}
	mixer := &Mixer{
		Gain:   1.0 / 8.0,
		Inputs: inputs,
	}
	mixer.Emit(out)
	log.Printf("HellEmitter Finish: %v", o)
}
