// Play sounds for a CW or QRSS beacon to a transmitter in Upper Sideband mode.
// Emits raw s16be mono audio to stdout.
//
// Usage:
//   go run qrss.go [flags] [words] | paplay --rate=44100 --channels=1 --format=s16le --raw /dev/stdin
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
	"github.com/strickyak/qrss-squeak/marimba"
)

var MODE = flag.String("mode", "", "Which mode to use.")
var LOOP = flag.Int64("loop", 0, "Repeat using this many seconds.  If 0, do not repeat.  Synchronizes to UNIX time modulo this many seconds.")
var LOOP_OFFSET = flag.Int64("loop_offset", 0, "Offset seconds within the loop.")
var DURATION = flag.Float64("duration", 0, "Specify total Duration (in seconds) instead of --dit time.")
var IMAGE = flag.String("image", "", "Image for marimba mode")

type ModeSpec struct {
	Func    func(text string) Emitter
	Explain string
}

var Modes = map[string]ModeSpec{
	// Usual modes:
	"cw":      ModeSpec{mainCW, "normal CW (single tone)"},
	"fs":      ModeSpec{mainFSCW, "Frequency Shift CW (low tone for gaps)"},
	"df":      ModeSpec{mainDFCW, "Dual Frequency CW (high tone for dahs)"},
	"tf":      ModeSpec{mainTFCW, "Three Frequency CW (mid-tone for OFF state)"},
	"hell":    ModeSpec{mainHell, "Hellschreiber: print human-readable ASCII directly in spectrum"},
	"marimba": ModeSpec{mainMarimba, "draw bitmap in spectrum with many tiny oscillators"},

	// Weird modes:
	"par":     ModeSpec{mainParallelCW, "CW letters in parallel (polyphonic)"},
	"decon5":  ModeSpec{mainDecon5, "Five deconstructed lines of / and gap"},
	"fractal": ModeSpec{mainFractal, "4-level fractal Dual Tone"},

	// Not useful execpt as demos:
	"demo-clock": ModeSpec{mainDemoClock, "demo of ticking clock"},
	"demo-4":     ModeSpec{mainDemoFour, "demo of four modes every 20 minutes"},
	"demo-junk":  ModeSpec{mainDemoJunk, "demo of cron & async"},
}

func mainCW(text string) Emitter {
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

func mainFSCW(text string) Emitter {
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

func mainDFCW(text string) Emitter {
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

func mainTFCW(text string) Emitter {
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

func mainHell(text string) Emitter {
	o := &HellConf{
		Dit:       Secs(*DIT),
		Freq:      0,
		Bandwidth: *BW,
		Text:      text,
		Tail:      true,
	}
	return NewHellEmitter(o)
}

func mainMarimba(text string) Emitter {
	o := &marimba.Conf{
		Dit:       Secs(*DIT),
		Freq:      0,
		Bandwidth: *BW,
		Filename:  *IMAGE,
		Gain:      0.25,
	}
	return NewAGC(marimba.NewEmitter(o))
}

func mainParallelCW(text string) Emitter {
	texts := strings.Split(strings.TrimSpace(text), ";")
	n := len(texts)
	delta := *BW / float64(n-1) // Difference between neighboring CW frequenices, in Hertz.
	var inputs []Emitter
	for i, s := range texts {
		inputs = append(inputs, NewCWEmitter(&CWConf{
			ToneWhenOff: false,
			Dit:         Secs(*DIT),
			Freq:        *BW - ((float64(i)) * delta),
			Bandwidth:   0,
			Text:        s,
			Tail:        false,
		}))
	}
	return &Mixer{Gain: 1 / float64(len(inputs)), Inputs: inputs}
}

func mainDemoClock(text string) Emitter {
	sum := NewAsyncMixer()

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

func mainFractal(_ string) Emitter {
	// Hardwire 12-element pattern here.
	pattern := []DiDah(" .-. . -.-  ") // (W6)REK
	mix := NewAsyncMixer()
	dit, _ := time.ParseDuration("694444444ns") // 100/144 seconds.

	// Total 14400-second daily transmisison.
	f := 0.0
	for ia, _ := range pattern { // 1200-second loop
		if pattern[ia] != ' ' {
			fa := f
			if pattern[ia] == '-' {
				fa = f + 20
			}
			for ib, _ := range pattern { // 100-second loop
				if pattern[ib] != ' ' {
					fb := fa
					if pattern[ib] == '-' {
						fb = fa + 10
					}

					func(_ia, _ib int) { // Capture loop vars
						p := &Cron{ // CW
							ModuloSeconds:    24 * 60 * 60, // diurnal.
							RemainderSeconds: 3600 + 1200*int64(_ia) + 100*int64(_ib),
							Run: func() {
								fractal := NewFractalEmitter(&FractalConf{
									ToneWhenOff: false,
									Dit:         dit,
									Freq:        fb,
									Bandwidth:   10,
									Morse:       pattern,
									Tail:        false,
								})
								mix.Add(&Gain{0.99, fractal})
							}}
						p.Start()
					}(ia, ib)

				}
			}
		}
	}
	return mix
}

func mainDecon5(keying string) Emitter {
	rows := strings.Split(keying, ",")
	if len(rows) != 5 {
		log.Fatalf("mainDecon5: expected 5 rows with commas in between: %s", keying)
	}

	sum := NewAsyncMixer()
	for i, row := range rows {
		func(_i int, _row string) { // Capture i as _i, row as _row
			p := &Cron{ // CW
				ModuloSeconds:    25 * 60,             // 25-minute cycles.
				RemainderSeconds: int64(_i)*5*60 + 15, // 5 minutes per row.
				Run: func() {
					cw := NewCWEmitter(&CWConf{
						ToneWhenOff: false,
						Dit:         time.Duration(*DIT * 1000000000), // try 5 seconds.
						Freq:        float64(_i) * (*BW / 5),
						Bandwidth:   1, // TODO: (*BW / 8),
						Morse:       []DiDah(_row),
						Tail:        false,
						NoGap:       true,
					})
					sum.Add(&Gain{0.99, cw})
				}}
			p.Start()
		}(i, row)
	}
	return sum
}

func mainDemoFour(text string) Emitter {
	sum := NewAsyncMixer()

	p1 := &Cron{ // CW
		ModuloSeconds:    20 * 60,    // twenty minutes
		RemainderSeconds: 0*300 + 10, // 00:00:10
		Run: func() {
			cw := NewCWEmitter(&CWConf{
				ToneWhenOff: false,
				Dit:         100 * time.Millisecond,
				Freq:        0,
				Bandwidth:   0,
				Text:        text,
				Tail:        true,
			})
			AdjustDuration(cw, 150*time.Second)
			sum.Add(&Gain{0.99, cw})
		}}
	p1.Start()

	p2 := &Cron{ // Dual Frequency
		ModuloSeconds:    20 * 60,    // twenty minutes
		RemainderSeconds: 1*300 + 10, // 00:05:10
		Run: func() {
			df := NewDFEmitter(&DFConf{
				ToneWhenOff: false,
				Dit:         100 * time.Millisecond,
				Freq:        0,
				Bandwidth:   8,
				Text:        text,
				Tail:        true,
			})
			AdjustDuration(df, 150*time.Second)
			sum.Add(&Gain{0.99, df})
		}}
	p2.Start()

	p3 := &Cron{ // Frequency Shift
		ModuloSeconds:    20 * 60,    // twenty minutes
		RemainderSeconds: 2*300 + 10, // 00:10:10
		Run: func() {
			fs := NewCWEmitter(&CWConf{
				ToneWhenOff: true,
				Dit:         100 * time.Millisecond,
				Freq:        0,
				Bandwidth:   8,
				Text:        text,
				Tail:        true,
			})
			AdjustDuration(fs, 150*time.Second)
			sum.Add(&Gain{0.99, fs})
		}}
	p3.Start()

	p4 := &Cron{ // Frequency Shift
		ModuloSeconds:    20 * 60,    // twenty minutes
		RemainderSeconds: 3*300 + 10, // 00:15:10
		Run: func() {
			hell := NewHellEmitter(&HellConf{
				Dit:       Secs(*DIT),
				Freq:      0,
				Bandwidth: 15,
				Text:      text,
				Tail:      true,
			})
			AdjustDuration(hell, 150*time.Second)
			sum.Add(&Gain{0.99, hell})
		}}
	p4.Start()

	return sum
}

func mainDemoJunk(text string) Emitter {
	sum := NewAsyncMixer()
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
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.SetPrefix("#")

	m, ok := Modes[*MODE]
	if !ok {
		for k, v := range Modes {
			log.Printf("Available --mode: %q  # %s", k, v.Explain)
		}
		log.Fatalf("Unknown --mode requested: %q", *MODE)
	}

	text := strings.Join(flag.Args(), " ")
	w := bufio.NewWriter(os.Stdout)

	em := m.Func(text)
	if *DURATION > 0 {
		AdjustDuration(em, SecondsToDuration(*DURATION))
	}
	if *LOOP > 0 && !*JUST_PRINT {
		sum := NewAsyncMixer()
		cron := &Cron{
			ModuloSeconds:    *LOOP,
			RemainderSeconds: *LOOP_OFFSET,
			Run: func() {
				sum.Add(em)
			}}
		cron.Start()
		Play(sum, w) // Runs forever.
	} else {
		Play(em, w)
		w.Flush()
	}
}
