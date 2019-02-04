package lib

import (
	"fmt"
	"time"
)

type Slice struct {
	Slice []float64
}

func (o *Slice) Duration() time.Duration {
	seconds := float64(len(o.Slice)) / (*FREQ)
	return time.Duration(seconds) * time.Second
}

func (o *Slice) DitPtr() *time.Duration {
	return nil
}

func (o *Slice) String() string {
	return fmt.Sprintf("Slice{%v}", o.Slice)
}

func (o *Slice) Emit(out chan Volt) {
	for _, x := range o.Slice {
		out <- Volt(x)
	}
}
