package lib

import (
	"bytes"
	"fmt"
	"time"
)

type Seq struct {
	Inputs []Emitter
}

func (o *Seq) Duration() time.Duration {
	var z time.Duration
	for _, e := range o.Inputs {
		z += e.Duration()
	}
	return z
}
func (o *Seq) DitPtr() *time.Duration {
	return nil
}

func (o *Seq) String() string {
	var bb bytes.Buffer
	fmt.Fprintf(&bb, "Seq{[")
	for i, e := range o.Inputs {
		fmt.Fprintf(&bb, "{%d}(%.3fs) %v ;;;; ", i, e.Duration().Seconds(), e)
	}
	fmt.Fprintf(&bb, "]}\n")
	return bb.String()
}

func (o *Seq) Emit(out chan Volt) {
	for _, e := range o.Inputs {
		e.Emit(out)
	}
}
