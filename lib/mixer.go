package lib

import (
	"bytes"
	"fmt"
	"time"
)

type Mixer struct {
	Gain   float64
	Inputs []Emitter
}

func (o *Mixer) Duration() time.Duration {
	var z time.Duration
	for _, e := range o.Inputs {
		d := e.Duration()
		if d > z {
			z = d
		}
	}
	return z
}

func (o *Mixer) String() string {
	var bb bytes.Buffer
	fmt.Fprintf(&bb, "Mixer{%f, ", o.Gain)
	for _, e := range o.Inputs {
		fmt.Fprintf(&bb, "%v, ", e)
	}
	fmt.Fprintf(&bb, "}")
	return bb.String()
}

func (o *Mixer) Emit(out chan Volt) {
	numInputs := len(o.Inputs)
	done := make([]bool, numInputs)
	ch := make([]chan Volt, numInputs)

	for i, e := range o.Inputs {
		ch[i] = make(chan Volt, small)
		go func(j int, x Emitter) {
			x.Emit(ch[j])
			close(ch[j])
		}(i, e)
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
}
