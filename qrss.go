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
	Func    func(string, *bufio.Writer)
	Explain string
}

var Modes = map[string]ModeSpec{
	"exp":  ModeSpec{mainExp, "experimental"},
	"cw":   ModeSpec{mainCW, "normal CW (single tone)"},
	"fscw": ModeSpec{mainFSCW, "Frequency Shift CW (low tone for gaps)"},
	"dfcw": ModeSpec{mainDFCW, "Dual Frequency CW (high tone for dahs)"},
	"tfcw": ModeSpec{mainTFCW, "Three Frequency CW (mid-tone for OFF state)"},
}

func mainCW(text string, w *bufio.Writer) {
	o := &CWConf{
		ToneWhenOff: false,
		Dit:         *DIT,
		Freq:        0,
		Width:       0,
		Text:        text,
		Tail:        true,
	}
	Play(NewCWEmitter(o), w)
}

func mainFSCW(text string, w *bufio.Writer) {
	o := &CWConf{
		ToneWhenOff: true,
		Dit:         *DIT,
		Freq:        0,
		Width:       *WIDTH,
		Text:        text,
		Tail:        true,
	}
	Play(NewCWEmitter(o), w)
}

func mainDFCW(text string, w *bufio.Writer) {
	o := &DFConf{
		ToneWhenOff: false,
		Dit:         *DIT,
		Freq:        0,
		Width:       *STEP,
		Text:        text,
		Tail:        true,
	}
	Play(NewDFEmitter(o), w)
}

func mainTFCW(text string, w *bufio.Writer) {
	o := &DFConf{
		ToneWhenOff: true,
		Dit:         *DIT,
		Freq:        0,
		Width:       *STEP,
		Text:        text,
		Tail:        true,
	}
	Play(NewDFEmitter(o), w)
}

func mainExp(text string, w *bufio.Writer) {
	flusher := func() {
		w.Flush()
	}
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
	Play(sum, w)
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
	spec.Func(text, w)
	w.Flush()
}
