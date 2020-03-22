package lib

import (
	"fmt"
	"log"
	"time"
)

type FractalConf struct {
	ToneWhenOff bool
	Dit         time.Duration
	Freq        float64
	Bandwidth   float64
	Morse       []DiDah
	Text        string
	Tail        bool
}
type FractalEmitter struct {
	Expand []int // 0=gap, 1 2 3 4 =tones
	FractalConf
}

func NewFractalEmitter(conf *FractalConf) *FractalEmitter {
	o := &FractalEmitter{
		FractalConf: *conf,
	}
	for _, c := range o.Morse {
		for _, d := range o.Morse {
			x := 0
			if c != ' ' {
				if c == '-' {
					x = 3
				} else {
					x = 1
				}
			}
			if d != ' ' {
				if d == '-' {
					x++
				}
			} else {
				x = 0
			}
			o.Expand = append(o.Expand, x)
		}
	}
	// Now trim 0's off the tail.
	for o.Expand[len(o.Expand)-1] == ' ' {
		o.Expand = o.Expand[:len(o.Expand)-1]
	}
	return o
}

func (o *FractalEmitter) DurationInDits() float64 {
	return float64(len(o.Expand))
}

func (o *FractalEmitter) Duration() time.Duration {
	return time.Duration(o.DurationInDits()) * o.Dit
}

func (o *FractalEmitter) DitPtr() *time.Duration {
	return &o.Dit
}

func (o *FractalEmitter) String() string {
	return fmt.Sprintf(
		"FractalEmitter{morse=%q,freq=%.1f,width=%.1f,dit=%v,total=%v}",
		o.Morse, o.Freq, o.Bandwidth, o.Dit, o.Duration())
}

var fractalRamp = time.Duration(20 * 1000 * 1000)

func (o *FractalEmitter) Emit(out chan Volt) {
	log.Printf("FractalEmitter Start: %v", o)

	for _, x := range o.Expand {
		if x == 0 {
			PlayGap(o.Dit, out)
		} else {
			f := o.Freq + ((float64(x) - 1) * (o.Bandwidth / 4))
			PlayTone(f, f+1, BOTH, o.Dit, fractalRamp, out)
		}
	}
	log.Printf("FractalEmitter Finish: %v", o)
}
