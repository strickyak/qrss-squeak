package lib

import (
	"fmt"
	"time"
	"log"
)

type CWEmitter struct {
	Dit	time.Duration
	Freq	float64
	Morse	[]DiDah
	Text	string
	Total time.Duration
}

func NewCWEmitter(text string, freq float64, tail bool) *CWEmitter {
	o := &CWEmitter{
		Dit: *DIT,
		Freq: freq,
		Morse: Morse(text, tail),
		Text: text,
	}
	o.Total = time.Duration(o.DurationInDits()) * o.Dit
	return o
}

func (o *CWEmitter) DurationInDits() float64 {
	var n float64
	for _, didah := range o.Morse {
		switch didah {
		case '.', ' ': n += 1
		case '-': n += 3
		default:
			log.Fatalf("bad didah: %d in %q", didah, o.Text)
		}
	}
	return n
}

func (o *CWEmitter) Duration() time.Duration {
	return o.Total
}

func (o *CWEmitter) String() string {
	return fmt.Sprintf("CWEmitter{text=%q,morse=%q,freq=%.1f,dit=%v,total=%v}", o.Text, o.Morse, o.Freq, o.Dit, o.Total)
}

func (o *CWEmitter) Emit(out chan Volt) {
	for _, didah := range o.Morse {
		switch didah {
		case '.': {
			PlayTone(o.Freq, o.Freq, BOTH, o.Dit, out)
		}
		case '-': {
			PlayTone(o.Freq, o.Freq, BOTH, 3 * o.Dit, out)
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
