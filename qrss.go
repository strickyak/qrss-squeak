package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	. "github.com/strickyak/qrss-squeak/lib"
)

var MODE = flag.String("mode", "", "Which mode to use")

type ModeSpec struct {
	Func    func(text string, flusher func()) Emitter
	Explain string
}

var Modes = map[string]ModeSpec{
	"clock1":  ModeSpec{mainDemoClock1, "demo of ticking clock"},
	"demo1":  ModeSpec{mainDemo1, "demo of cron & async"},
	"cw":   ModeSpec{mainCW, "normal CW (single tone)"},
	"fscw": ModeSpec{mainFSCW, "Frequency Shift CW (low tone for gaps)"},
	"dfcw": ModeSpec{mainDFCW, "Dual Frequency CW (high tone for dahs)"},
	"tfcw": ModeSpec{mainTFCW, "Three Frequency CW (mid-tone for OFF state)"},
}

func mainCW(text string, flusher func()) Emitter {
	o := &CWConf{
		ToneWhenOff: false,
		Dit:         *DIT,
		Freq:        0,
		Width:       0,
		Text:        text,
		Tail:        true,
	}
	return NewCWEmitter(o)
}

func mainFSCW(text string, flusher func()) Emitter {
	o := &CWConf{
		ToneWhenOff: true,
		Dit:         *DIT,
		Freq:        0,
		Width:       *WIDTH,
		Text:        text,
		Tail:        true,
	}
	return NewCWEmitter(o)
}

func mainDFCW(text string, flusher func()) Emitter {
	o := &DFConf{
		ToneWhenOff: false,
		Dit:         *DIT,
		Freq:        0,
		Width:       *STEP,
		Text:        text,
		Tail:        true,
	}
	return NewDFEmitter(o)
}

func mainTFCW(text string, flusher func()) Emitter {
	o := &DFConf{
		ToneWhenOff: true,
		Dit:         *DIT,
		Freq:        0,
		Width:       *STEP,
		Text:        text,
		Tail:        true,
	}
	return NewDFEmitter(o)
}

func mainDemoClock1(text string, flusher func()) Emitter {
	sum := NewAsyncSum(flusher)

	p1 := &Cron{
		ModuloSeconds:    10,
		RemainderSeconds: 0,
		Run: func() {
			sum.Add(NewCWEmitter(&CWConf{
				ToneWhenOff: false,
				Dit:         500 * time.Millisecond,
				Freq:        0,
				Width:       0,
				Text:        "e",
				Tail:        false,
			}))
		}}
	p1.Start()

	for i := 1; i <= 3; i++ {
		p2 := &Cron{
			ModuloSeconds:    10,
			RemainderSeconds: 10-int64(i),
			Run: func() {
				sum.Add(NewCWEmitter(&CWConf{
					ToneWhenOff: false,
					Dit:         100 * time.Millisecond,
					Freq:        0,
					Width:       0,
					Text:        "e",
					Tail:        false,
				}))
			}}
		p2.Start()
	}

	return sum

}

func mainDemo1(text string, flusher func()) Emitter {
	sum := NewAsyncSum(flusher)
	p1 := &Cron{
		ModuloSeconds:    20,
		RemainderSeconds: 0,
		Run: func() {
			cw := NewCWEmitter(&CWConf{
				ToneWhenOff: false,
				Dit:         100 * time.Millisecond,
				Freq:        0,
				Width:       0,
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
				Width:       0,
				Text:        "k",
				Tail:        true,
			})
			sum.Add(&Gain{0.5, cw})
		}}
	p2.Start()

	for i := 0; i < 4; i++ {
		p2 := &Cron{
			ModuloSeconds:    15,
			RemainderSeconds: 15-int64(i),
			Run: func() {
				cw := NewCWEmitter(&CWConf{
					ToneWhenOff: false,
					Dit:         100 * time.Millisecond,
					Freq:        500,
					Width:       0,
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
	spec, ok := Modes[*MODE]
	if !ok {
		for k, v := range Modes {
			log.Printf("Available --mode: %q  # %s", k, v.Explain)
		}
		log.Fatalf("Unknown --mode requested: %q", *MODE)
	}
	text := strings.Join(flag.Args(), " ") + " "
	w := bufio.NewWriter(os.Stdout)
	flusher := func() {w.Flush()}
	Play(spec.Func(text, flusher), w)
	w.Flush()
}
