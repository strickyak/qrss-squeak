package lib

import (
	"fmt"
	"time"
	"log"
)

type DFEmitter struct {
	Dit	time.Duration
	Freq	float64
	DeltaFreq	float64
	Morse	[]DiDah
	Text	string
	Total time.Duration
}

func NewDFEmitter(text string, freq, deltaFreq float64, tail bool) *DFEmitter {
	o := &DFEmitter{
		Dit: *DIT,
		Freq: freq,
		DeltaFreq: deltaFreq,
		Morse: Morse(text, tail),
		Text: text,
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
	return fmt.Sprintf("DFEmitter{text=%q,morse=%q,freq=%.1f,deltaFreq=%.1f,dit=%v,total=%v}", o.Text, o.Morse, o.Freq, o.DeltaFreq, o.Dit, o.Total)
}

func (o *DFEmitter) Emit(out chan Volt) {
	for _, didah := range o.Morse {
		switch didah {
		case '.': {
			f := o.Freq
			PlayTone(f, f, BOTH, o.Dit, out)
		}
		case '-': {
			f := o.Freq + o.DeltaFreq
			PlayTone(f, f, BOTH, o.Dit, out)
		}
		case ' ': {
			PlayGap(o.Dit, out)
		}
		default:
			log.Fatalf("bad didah: %d in %q", didah, o.Text)
		}
		PlayGap(o.Dit, out)
	}
}
