// Play sounds for a CW or QRSS beacon to a transmitter in Upper Sideband mode.
// Emits raw s16be mono audio to stdout.
//
// Usage:
//   go run qrss.go [flags] [words] | paplay --rate=44100 --channels=1 --format=s16le --raw --latency-msec=30 /dev/stdin
//
// Flags:
//   --mode=cw     (or any other mode defined in Modes below)
//   --rate=44100  (matches the --rate flag to paplay)
//   --dit=0.1s    (dit time.  Add "s" for seconds.)
//
// Demo:
// go run ../qrss.go -mode=cw -dit=0.1 -loop=5 hi hi | paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin

package main

// TODO -- fix zero values to defaults.

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	. "github.com/strickyak/qrss-squeak/lib"
)

var MODE = flag.String("mode", "", "Which mode to use.")
var LOOP = flag.Int64("loop", 0, "Repeat using this many seconds.  If 0, do not repeat.  Synchronizes to UNIX time modulo this many seconds.")
var LOOP_OFFSET = flag.Int64("loop_offset", 0, "Offset seconds within the loop.")

type ModeSpec struct {
	Func    func(text string, flusher func()) Emitter
	Explain string
}

var Modes = map[string]ModeSpec{
	// Usual modes:
	"cw": ModeSpec{mainCW, "normal CW (single tone)"},
	"fs": ModeSpec{mainFSCW, "Frequency Shift CW (low tone for gaps)"},
	"df": ModeSpec{mainDFCW, "Dual Frequency CW (high tone for dahs)"},
	"tf": ModeSpec{mainTFCW, "Three Frequency CW (mid-tone for OFF state)"},

	// Weird modes:
	"par": ModeSpec{mainParallelCW, "CW letters in parallel (polyphonic)"},

	// Not useful execpt as demos:
	"demo-clock": ModeSpec{mainDemoClock, "demo of ticking clock"},
	"demo-junk":  ModeSpec{mainDemoJunk, "demo of cron & async"},
}

func mainCW(text string, flusher func()) Emitter {
	o := &CWConf{
		ToneWhenOff: false,
		Dit:         Secs(*DIT),
		Freq:        0,
		Bandwidth:   0,
		Text:        text,
		Tail:        true,
	}
	return NewCWEmitter(o)
}

func mainFSCW(text string, flusher func()) Emitter {
	o := &CWConf{
		ToneWhenOff: true,
		Dit:         Secs(*DIT),
		Freq:        0,
		Bandwidth:   *BW,
		Text:        text,
		Tail:        true,
	}
	return NewCWEmitter(o)
}

func mainDFCW(text string, flusher func()) Emitter {
	o := &DFConf{
		ToneWhenOff: false,
		Dit:         Secs(*DIT),
		Freq:        0,
		Bandwidth:   *BW,
		Text:        text,
		Tail:        true,
	}
	return NewDFEmitter(o)
}

func mainTFCW(text string, flusher func()) Emitter {
	o := &DFConf{
		ToneWhenOff: true,
		Dit:         Secs(*DIT),
		Freq:        0,
		Bandwidth:   *BW,
		Text:        text,
		Tail:        true,
	}
	return NewDFEmitter(o)
}

func mainParallelCW(text string, flusher func()) Emitter {
	text = strings.TrimSpace(text)
	var inputs []Emitter
	n := 0
	for _, _ = range text {
		n++
	}
	delta := *BW / float64(n) // Difference between to neighboring CW frequenices, in Hertz.
	for i, r := range text {
		inputs = append(inputs, NewCWEmitter(&CWConf{
			ToneWhenOff: false,
			Dit:         Secs(*DIT),
			Freq:        *BW - ((float64(i) + 0.5) * delta),
			Bandwidth:   0,
			Text:        string(r),
			Tail:        false,
		}))
	}
	return &Sum{Gain: 1 / float64(len(inputs)), Inputs: inputs}
}

func mainDemoClock(text string, flusher func()) Emitter {
	sum := NewAsyncSum(flusher)

	p1 := &Cron{
		ModuloSeconds:    10,
		RemainderSeconds: 0,
		Run: func() {
			sum.Add(NewCWEmitter(&CWConf{
				ToneWhenOff: false,
				Dit:         500 * time.Millisecond,
				Freq:        0,
				Bandwidth:   0,
				Text:        "e",
				Tail:        false,
			}))
		}}
	p1.Start()

	for i := 1; i <= 3; i++ {
		p2 := &Cron{
			ModuloSeconds:    10,
			RemainderSeconds: 10 - int64(i),
			Run: func() {
				sum.Add(NewCWEmitter(&CWConf{
					ToneWhenOff: false,
					Dit:         100 * time.Millisecond,
					Freq:        0,
					Bandwidth:   0,
					Text:        "e",
					Tail:        false,
				}))
			}}
		p2.Start()
	}

	return sum

}

func mainDemoJunk(text string, flusher func()) Emitter {
	sum := NewAsyncSum(flusher)
	p1 := &Cron{
		ModuloSeconds:    20,
		RemainderSeconds: 0,
		Run: func() {
			cw := NewCWEmitter(&CWConf{
				ToneWhenOff: false,
				Dit:         100 * time.Millisecond,
				Freq:        0,
				Bandwidth:   0,
				Text:        text,
				Tail:        true,
			})
			sum.Add(&Gain{0.5, cw})
		}}
	p1.Start()
	p2 := &Cron{
		ModuloSeconds:    20,
		RemainderSeconds: 15,
		Run: func() {
			cw := NewCWEmitter(&CWConf{
				ToneWhenOff: false,
				Dit:         200 * time.Millisecond,
				Freq:        100,
				Bandwidth:   0,
				Text:        "k",
				Tail:        true,
			})
			sum.Add(&Gain{0.5, cw})
		}}
	p2.Start()

	for i := 0; i < 4; i++ {
		p2 := &Cron{
			ModuloSeconds:    15,
			RemainderSeconds: 15 - int64(i),
			Run: func() {
				cw := NewCWEmitter(&CWConf{
					ToneWhenOff: false,
					Dit:         100 * time.Millisecond,
					Freq:        500,
					Bandwidth:   0,
					Text:        "e",
					Tail:        false,
				})
				sum.Add(&Gain{0.5, cw})
			}}
		p2.Start()
	}

	return sum
}

func main() {
	flag.Parse()
	m, ok := Modes[*MODE]
	if !ok {
		for k, v := range Modes {
			log.Printf("Available --mode: %q  # %s", k, v.Explain)
		}
		log.Fatalf("Unknown --mode requested: %q", *MODE)
	}
	text := strings.Join(flag.Args(), " ")
	w := bufio.NewWriter(os.Stdout)
	flusher := func() { w.Flush() }

	if *LOOP > 0 {
		sum := NewAsyncSum(flusher)
		cron := &Cron{
			ModuloSeconds:    *LOOP,
			RemainderSeconds: *LOOP_OFFSET,
			Run: func() {
				sum.Add(m.Func(text, flusher))
			}}
		cron.Start()
		Play(sum, w) // Runs forever.
	} else {
		Play(m.Func(text, flusher), w)
		w.Flush()
	}
}
