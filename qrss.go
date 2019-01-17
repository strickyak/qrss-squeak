package main

import (
	"bufio"
	"flag"
	"log"
	"strings"
	"os"

	. "github.com/strickyak/qrss-squawk/lib"
)

var MODE = flag.String("mode", "", "Which mode to use")
var X = flag.Bool("x", false, "Just print")

var Modes = map[string]func(*bufio.Writer) {
	"cw": mainCW,
	"df": mainDF,
}

func mainCW(w *bufio.Writer) {
	text := strings.Join(flag.Args(), " ") + " "
	cw := NewCWEmitter(text, 0.0 /*freq*/, true /*tail*/)
	if *X {
		log.Printf("%v", cw)
		os.Exit(0)
	}

	//  Output(volts chan Volt, w io.Writer, done chan bool)
	volts := make(chan Volt, int(*RATE))  // for 1.0s
	done := make(chan bool)

	go Output(volts, w, done)
	cw.Emit(volts)
	close(volts)
	<-done
}

func mainDF(w *bufio.Writer) {
	text := strings.Join(flag.Args(), " ") + " "
	df := NewDFEmitter(text, 0.0 /*freq*/, *STEP, true /*tail*/)
	if *X {
		log.Printf("%v", df)
		os.Exit(0)
	}

	//  Output(volts chan Volt, w io.Writer, done chan bool)
	volts := make(chan Volt, int(*RATE))  // for 1.0s
	done := make(chan bool)

	go Output(volts, w, done)
	df.Emit(volts)
	close(volts)
	<-done
}

func main() {
	flag.Parse()
	f, ok := Modes[*MODE]
	if !ok {
		for k := range Modes {
			log.Printf("Available --mode: %q", k)
		}
		log.Fatalf("Unknown --mode requested: %q", *MODE)
	}
	w := bufio.NewWriter(os.Stdout)
	f(w)
	w.Flush()
}
