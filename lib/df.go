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
	Width       float64
	Morse       []DiDah
	Text        string
	Tail        bool
}
type DFEmitter struct {
	DFConf
	Total time.Duration
}

func NewDFEmitter(conf *DFConf) *DFEmitter {
	o := &DFEmitter{
		DFConf: *conf,
	}
	// If Text is provided but not Morse, convert Text & Tail to Morse.
	if len(o.Morse) == 0 {
		o.Morse = Morse(o.Text, o.Tail)
	}
	o.Total = time.Duration(o.DurationInDits()) * o.Dit
	return o
}

func (o *DFEmitter) DurationInDits() float64 {
	return float64(len(o.Morse))
}

func (o *DFEmitter) Duration() time.Duration {
	return o.Total
}

func (o *DFEmitter) String() string {
	if o.ToneWhenOff {
		return fmt.Sprintf("DFEmitter{ToneWhenOff,text=%q,morse=%q,freq=%.1f,width=%.1f,dit=%v,total=%v}", o.Text, o.Morse, o.Freq, o.Width, o.Dit, o.Total)
	} else {
		return fmt.Sprintf("DFEmitter{text=%q,morse=%q,freq=%.1f,width=%.1f,dit=%v,total=%v}", o.Text, o.Morse, o.Freq, o.Width, o.Dit, o.Total)
	}
}

func (o *DFEmitter) Emit(out chan Volt) {
	log.Printf("DFEmitter Start: %v", o)
	f0, f1, f2 := o.Freq, o.Freq, o.Freq+o.Width // OFF, Dit, Dah frequencies.
	if o.ToneWhenOff {
		f0 = (f1 + f2) / 2
	}
	gap := func() {
		if o.ToneWhenOff {
			PlayTone(f0, f0, BOTH, o.Dit, out)
		} else {
			PlayGap(o.Dit, out)
		}
	}

	lenMorse, i := len(o.Morse), 0
	for _, didah := range o.Morse {
		i++
		switch didah {
		case '.':
			PlayTone(f1, f1, BOTH, o.Dit, out)
		case '-':
			PlayTone(f2, f2, BOTH, o.Dit, out)
		case ' ':
			gap()
		default:
			log.Fatalf("bad didah: %d in %q", didah, o.Text)
		}

		if o.Tail || i < lenMorse {
			gap()
		}
	}
	log.Printf("DFEmitter Finish: %v", o)
}
