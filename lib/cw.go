package lib

import (
	"fmt"
	"log"
	"time"
)

type CWConf struct {
	ToneWhenOff bool
	Dit         time.Duration
	Freq        float64
	Bandwidth   float64
	Morse       []DiDah
	Text        string
	Tail        bool
	NoGap       bool // TODO: get rid of quick hack
}
type CWEmitter struct {
	CWConf
	Total time.Duration
}

func NewCWEmitter(conf *CWConf) *CWEmitter {
	o := &CWEmitter{
		CWConf: *conf,
	}
	// If Text is provided but not Morse, convert Text & Tail to Morse.
	if len(o.Morse) == 0 {
		o.Morse = Morse(o.Text, o.Tail)
	}
	o.Total = time.Duration(o.DurationInDits()) * o.Dit
	return o
}

func (o *CWEmitter) DurationInDits() float64 {
	var n float64
	for _, didah := range o.Morse {
		switch didah {
		case '.', ' ':
			if o.NoGap {
				n += 1 // 1 for ON, 1 for OFF.
			} else {
				n += 2 // 1 for ON, 1 for OFF.
			}
		case '-':
			if o.NoGap {
				n += 3 // 3 for ON, 1 for OFF.
			} else {
				n += 4 // 3 for ON, 1 for OFF.
			}
		default:
			log.Fatalf("bad didah: %d in %q", didah, o.Text)
		}
	}
	if !o.Tail {
		n-- // Delete final OFF.
	}
	return n
}

func (o *CWEmitter) Duration() time.Duration {
	return o.Total
}

func (o *CWEmitter) String() string {
	if o.ToneWhenOff {
		return fmt.Sprintf("CWEmitter{ToneWhenOff,text=%q,morse=%q,freq=%.1f,width=%.1f,dit=%v,total=%v}", o.Text, o.Morse, o.Freq, o.Bandwidth, o.Dit, o.Total)
	} else {
		return fmt.Sprintf("CWEmitter{text=%q,morse=%q,freq=%.1f,dit=%v,total=%v}", o.Text, o.Morse, o.Freq, o.Dit, o.Total)
	}
}

func (o *CWEmitter) Emit(out chan Volt) {
	f0, f1 := o.Freq, o.Freq // OFF, ON frequencies.
	if o.ToneWhenOff {
		f1 += o.Bandwidth
	}
	log.Printf("CWEmitter Start: %v", o)
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
			PlayTone(f1, f1, BOTH, 3*o.Dit, out)
		case ' ':
			gap()
		default:
			log.Fatalf("bad didah: %d in %q", didah, o.Text)
		}

		if o.Tail || i < lenMorse {
			if !o.NoGap {
				gap()
			}
		}
	}
	log.Printf("CWEmitter Finish: %v", o)
}
