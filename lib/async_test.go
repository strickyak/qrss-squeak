package lib

import (
	"log"
	"testing"
	. "time"
)

func TestAsyncSum1(t *testing.T) {
	count := make(chan struct{}, 1)
	nop := func() {
		count <- struct{}{}
	}

	a := NewAsyncSum(nop)

	ch := make(chan Volt, small)
	go a.Emit(ch)

	a.Add(&Slice{Slice: []float64{30, 40, 50}})
	a.Add(&Slice{Slice: []float64{30, 40, 50}})
	Sleep(Microsecond)
	a.Add(&Slice{Slice: []float64{30, 40, 50}})
	a.Add(&Slice{Slice: []float64{30, 40, 50}})
	Sleep(Microsecond)
	a.Add(&Slice{Slice: []float64{30, 40, 50}})
	a.Add(&Slice{Slice: []float64{30, 40, 50}})
	Sleep(Microsecond)

	counter := 0
	var sum Volt
Loop:
	for {
		select {
		case _, ok := <-count:
			Must(ok)
			counter++
			log.Printf("Count: %d", counter)
		case n, ok := <-ch:
			Must(ok)
			sum += n
			log.Printf("Sum: %g", sum)
		case <-After(5 * Millisecond):
			log.Print("After")
			if sum == 720 {
				break Loop
			}
		}
	}
	if sum != 720 {
		t.Errorf("sum got %g, want 360", sum)
	}
}
