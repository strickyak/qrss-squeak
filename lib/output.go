package lib

/*
Sun Feb  3 13:35:38 PST 2019
Did some debugging on bad buffering in Output.

Problem 1.  Sometimes the output would sound weak and scratchy,
as if an odd number of bytes got flushed, and it was reading
two-byte integers out of sync.

Problem 2.  This loop was sounding twice every 8 seconds, instead of once every 4 seconds:
```
go run qrss.go -mode=cw -dit=0.1 -loop=4 "k" |
paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
```

Solution was to add a timeout of 1ms when not flushed, and flush if that timeout fires.
I've also shortened the channel buffer of Volts.
This assumes the Emitters can emit Volts faster than the player can consume them.
Decrease the --rate=8000 if that becomes a problem.
But I've read that although most sound cards can operate at 8000 samples per second,
they often don't have proper hardware filtering for that, and are best at 44100.
*/

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"
)

var JUST_PRINT = flag.Bool("x", false, "Just print")

const MaxShort = 0x7FFF

// It multiplies by an overall gain, converts to signed int16, and writes to the writer in little-endian format.
// When the input volts channel has no more, we write true to the done channel,
// so the main program can exit.
// Volts below -1.0 or above +1.0 get clipped hard.
func Output(volts chan Volt, w *bufio.Writer, done chan bool) {
	var flushed bool
LOOP:
	for {
		var ok bool
		var volt Volt
		if flushed {
			// OK to block if flushed.
			log.Println("Output: Probably blocking.")
			volt, ok = <-volts
			log.Println("Output: Awake.")
			if !ok {
				break LOOP // Signal done and return.
			}
		} else {
			// Don't block for long if not flushed.
			select {
			case volt, ok = <-volts:
				if !ok {
					break LOOP // Signal done and return.
				}
			case <-time.After(time.Millisecond):
				log.Println("Output: Flushing.")
				w.Flush()
				flushed = true
				continue LOOP // Now that we flushed, we can block.
			}
		}

		y := *GAIN * float64(volt)
		// Clip hard at +/- 1 unit.
		if y > 1.0 {
			// allow clipping 0.5% just for unexplained occurances.
			if y > 1.005 && !*CLIP {
				log.Panicf("Clipping not allowed without --clip flag: y=%g", y)
			}
			y = 1.0
		}
		if y < -1.0 {
			// allow clipping 0.5% just for unexplained occurances.
			if y < -1.005 && !*CLIP {
				log.Panicf("Clipping not allowed without --clip flag: y=%g", y)
			}
			y = -1.0
		}
		yShort := int(MaxShort * y)
		w.WriteByte(byte(yShort))
		w.WriteByte(byte(yShort >> 8))
		flushed = false
	}
	done <- true
}

func Play(em Emitter, w *bufio.Writer) {
	log.Printf("Play: %v", em)
	if *JUST_PRINT {
		os.Exit(0)
	}

	volts := make(chan Volt, int(*RATE/100)) // buffer chan for 10ms
	done := make(chan bool)

	go Output(volts, w, done)
	em.Emit(volts)
	close(volts)
	<-done
}
