package lib

import (
	"fmt"
	"time"
)

type AGC struct {
	Volts    []Volt
	Max      Volt
	InputStr string
}

func NewAGC(em Emitter) Emitter {
	ch := make(chan Volt)
	go func() {
		em.Emit(ch)
		close(ch)
	}()
	var vec []Volt
	var max Volt
	for v := range ch {
		if -v > max {
			max = -v
		}
		if v > max {
			max = v
		}
		vec = append(vec, v)
	}
	Must(max > 0.0)
	for i, v := range vec {
		vec[i] = v / max
	}
	return &AGC{
		Volts:    vec,
		Max:      max,
		InputStr: em.String(),
	}
}

func (o *AGC) Duration() time.Duration {
	return Secs(float64(len(o.Volts)) / *RATE)
}
func (o *AGC) DitPtr() *time.Duration {
	return nil
}

func (o *AGC) String() string {
	return fmt.Sprintf("AGC{len=%d,dur=%v,max=%f, <- %s}", len(o.Volts), o.Duration(), o.Max, o.InputStr)
}

func (o *AGC) Emit(out chan Volt) {
	for _, volt := range o.Volts {
		out <- volt
	}
}
