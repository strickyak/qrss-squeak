package lib

import (
	"fmt"
	"time"
)

const small = 10

type Gain struct {
	Gain  float64
	Input Emitter
}

func (o *Gain) Duration() time.Duration {
	return o.Input.Duration()
}

func (o *Gain) String() string {
	return fmt.Sprintf("Gain{%f, %v}", o.Gain, o.Input)
}

func (o *Gain) Emit(out chan Volt) {
	in := make(chan Volt, small)
	go o.Input.Emit(in)
	for {
		x, ok := <-in
		if !ok {
			break
		}
		out <- Volt(o.Gain) * x
	}
	close(out)
}
