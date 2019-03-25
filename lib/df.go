package lib

import (
	"fmt"
	"log"
	"time"
)

type DFConf struct {
	ToneWhenOff bool
	Dit         time.Duration
	Freq        float64
	Bandwidth   float64
	Morse       []DiDah
	Text        string
	Tail        bool
}
type DFEmitter struct {
	DFConf
}

func NewDFEmitter(conf *DFConf) *DFEmitter {
	o := &DFEmitter{
		DFConf: *conf,
	}
	// If Text is provided but not Morse, convert Text & Tail to Morse.
	if len(o.Morse) == 0 {
		o.Morse = Morse(o.Text, o.Tail)
	}
	return o
}

func (o *DFEmitter) DurationInDits() float64 {
	return float64(len(o.Morse))
}

func (o *DFEmitter) Duration() time.Duration {
	return time.Duration(o.DurationInDits()) * o.Dit
}

func (o *DFEmitter) DitPtr() *time.Duration {
	return &o.Dit
}

func (o *DFEmitter) String() string {
	if o.ToneWhenOff {
		return fmt.Sprintf(
			"DFEmitter{ToneWhenOff,text=%q,morse=%q,freq=%.1f,width=%.1f,dit=%v,total=%v}",
			o.Text, o.Morse, o.Freq, o.Bandwidth, o.Dit, o.Duration())
	} else {
		return fmt.Sprintf(
			"DFEmitter{text=%q,morse=%q,freq=%.1f,width=%.1f,dit=%v,total=%v}",
			o.Text, o.Morse, o.Freq, o.Bandwidth, o.Dit, o.Duration())
	}
}

func (o *DFEmitter) Emit(out chan Volt) {
	log.Printf("DFEmitter Start: %v", o)
	f1, f2 := o.Freq, o.Freq+o.Bandwidth // OFF, Dit, Dah frequencies.
	f0 := (f1 + f2) / 2
	gap := func() {
		if o.ToneWhenOff {
			PlayTone(f0, f0, BOTH, o.Dit, AutomaticRampDuration, out)
		} else {
			PlayGap(o.Dit, out)
		}
	}

	for _, didah := range o.Morse {
		switch didah {
		case '.':
			PlayTone(f1, f1, BOTH, o.Dit, AutomaticRampDuration, out)
		case '-':
			PlayTone(f2, f2, BOTH, o.Dit, AutomaticRampDuration, out)
		case ' ':
			gap()
		default:
			log.Fatalf("bad didah: %d in %q", didah, o.Text)
		}
	}
	log.Printf("DFEmitter Finish: %v", o)
}
