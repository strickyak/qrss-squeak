package lib

import (
	"fmt"
	"log"
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

func (o *Gain) DitPtr() *time.Duration {
	return o.Input.DitPtr()
}

func (o *Gain) String() string {
	return fmt.Sprintf("Gain{%f, %v}", o.Gain, o.Input)
}

func (o *Gain) Emit(out chan Volt) {
	log.Printf("Gain Start: %v", o)
	in := make(chan Volt, small)
	go func() {
		o.Input.Emit(in)
		close(in)
	}()
	for {
		x, ok := <-in
		if !ok {
			break
		}
		out <- Volt(o.Gain) * x
	}
	log.Printf("Gain Finish: %v", o)
}
