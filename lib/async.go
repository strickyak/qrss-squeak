package lib

import (
	"log"
	"time"
)

var InfiniteDuration = 100 * 366 * 24 * time.Hour // A hundred years.

// AsyncSum must be created with NewAsyncSum().
// Then you can add Emitters to it with Add().
// Those Emitters may be finite or infinite.
// You can add more Emitters later at any time.
// The AsyncSum will pause writing voltages to its output channel
// while there are no inputs, so this should be the final thing
// before Output().   AsyncSum is an infinite Emitter,
// but may have long (or infinite) pauses.  The flusher input
// is a function that gets called when the AsyncSum pauses.
// It should flush pending output so it is not held in buffers
// during pauses.
type AsyncSum struct {
	channels []chan Volt
	adding   chan chan Volt
	flusher  func()
}

func NewAsyncSum(flusher func()) *AsyncSum {
	return &AsyncSum{
		adding:  make(chan chan Volt, small),
		flusher: flusher,
	}
}

func (o *AsyncSum) Add(e Emitter) {
	ch := make(chan Volt)
	go e.Emit(ch)
	o.adding <- ch
}

func (o *AsyncSum) Duration() time.Duration {
	return InfiniteDuration // Never stops -- but may have many long gaps.
}

func (o *AsyncSum) String() string {
	return "AsyncSum{}"
}

func sliceContainsInt(vec []int, x int) bool {
	for _, e := range vec {
		if e == x {
			return true
		}
	}
	return false
}

func (o *AsyncSum) Emit(out chan Volt) {
	for {
		// First look for a new channel being added.
		// Just take one; we can get another on the next loop.
		// If there are no channels, block waiting on one to arrive.
		if len(o.channels) > 0 {
			select {
			case e := <-o.adding:
				o.channels = append(o.channels, e)
			default:
				// Don't block.
			}
		} else {
			o.flusher()         // Flush before blocking!
			e, ok := <-o.adding // Block.
			if !ok {
				log.Fatalf("Not expecting `adding` to close.")
			}
			o.channels = append(o.channels, e)
		}

		// Now get a voltage from each channel.
		// If the channel is closed, remember it, so we can remove it.
		// Output the sum of the voltages, unless all channels were closed.
		var channelsToBeDropped []int
		var sum Volt
		gotSomething := false
		for i, ch := range o.channels {
			x, ok := <-ch
			if !ok {
				channelsToBeDropped = append(channelsToBeDropped, i)
			} else {
				gotSomething = true
				sum += x
			}
		}
		if gotSomething {
			out <- sum
		}

		if len(channelsToBeDropped) > 0 {
			// Rebuild o.channels using tmp; omit those to be dropped.
			var tmp []chan Volt
			for i, ch := range o.channels {
				if !sliceContainsInt(channelsToBeDropped, i) {
					tmp = append(tmp, ch)
				}

			}
			o.channels = tmp
		}
	}
}
