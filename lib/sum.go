package lib

import (
	"bytes"
	"fmt"
	"time"
)

type Sum struct {
	Gain   float64
	Inputs []Emitter
}

func (o *Sum) Duration() time.Duration {
	var z time.Duration
	for _, e := range o.Inputs {
		d := e.Duration()
		if d > z {
			z = d
		}
	}
	return z
}

func (o *Sum) String() string {
	var bb bytes.Buffer
	fmt.Fprintf(&bb, "Sum{%f, ", o.Gain)
	for _, e := range o.Inputs {
		fmt.Fprintf(&bb, "%v, ", e)
	}
	fmt.Fprintf(&bb, "}")
	return bb.String()
}

func (o *Sum) Emit(out chan Volt) {
	numInputs := len(o.Inputs)
	done := make([]bool, numInputs)
	ch := make([]chan Volt, numInputs)

	for i, e := range o.Inputs {
		ch[i] = make(chan Volt, small)
		go e.Emit(ch[i])
	}

	numDone := 0
	for numDone < numInputs {
		var sum Volt
		for i, _ := range o.Inputs {
			if done[i] {
				continue
			}
			x, ok := <-ch[i]
			if !ok {
				numDone++
				done[i] = true
				continue
			}
			sum += x
		}
		out <- Volt(o.Gain) * sum
	}
	close(out)
}
