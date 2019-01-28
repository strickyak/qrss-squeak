package lib

import (
	"flag"
	"log"
	"os/exec"
	"strings"
	"time"
)

var TX_ON = flag.String("tx_on", "", "Shell command to run to turn transmitter on.")
var TX_OFF = flag.String("tx_off", "", "Shell command to run to turn transmitter off.  Should sleep first for player buffer to empty.")

// AsyncMixer must be created with NewAsyncMixer().
// Then you can add Emitters to it with Add().
// Those Emitters may be finite or infinite.
// You can add more Emitters later at any time.
// The AsyncMixer will pause writing voltages to its output channel
// while there are no inputs, so this should be the final thing
// before Output().   AsyncMixer is an infinite Emitter,
// but may have long (or infinite) pauses.  The flusher input
// is a function that gets called when the AsyncMixer pauses.
// It should flush pending output so it is not held in buffers
// during pauses.
type AsyncMixer struct {
	channels []chan Volt
	adding   chan chan Volt
	flusher  func()
}

func NewAsyncMixer(flusher func()) *AsyncMixer {
	return &AsyncMixer{
		adding:  make(chan chan Volt, small),
		flusher: flusher,
	}
}

func Transmit(on bool) {
	var command string
	if on {
		command = *TX_ON
	} else {
		command = *TX_OFF
	}
	if len(strings.TrimSpace(command)) == 0 {
		return // No command.
	}
	log.Printf("Transmitter (%v) Command: %q", on, command)
	err := exec.Command("/bin/bash", "-c", "set -x; "+command).Run()
	if err != nil {
		log.Fatalf("FATAL: Transmitter Error: %q: %v", command, err)
	}
}

func (o *AsyncMixer) Add(e Emitter) {
	ch := make(chan Volt)
	go func() {
		e.Emit(ch)
		close(ch)
	}()
	o.adding <- ch
}

func (o *AsyncMixer) Duration() time.Duration {
	return InfiniteDuration // Never stops -- but may have many long gaps.
}

func (o *AsyncMixer) String() string {
	return "AsyncMixer{}"
}

func sliceContainsInt(vec []int, x int) bool {
	for _, e := range vec {
		if e == x {
			return true
		}
	}
	return false
}

func (o *AsyncMixer) Emit(out chan Volt) {
	Transmit(true)
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
			o.flusher() // Flush before blocking!
			Transmit(false)
			log.Printf("AsyncMixer blocking.")
			e, ok := <-o.adding // Block.
			if !ok {
				log.Fatalf("Not expecting `adding` to close.")
			}
			Transmit(true)
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
