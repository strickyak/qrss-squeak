// Produce an image of a 5-channel Baudot Paper tape, the same encoding used in RTTY mode in ham radio.
// Usage:  go run baudot.go > baudot.svg
package main

import . "fmt"

const HEIGHT = 650

var WIDTH = len(Message)*100 + 700

const (
	LTRS = "oo.ooo"
	W    = "oo.  o"
	FIGS = "oo. oo"
	_6   = "o .o o"
	R    = " o. o "
	E    = "o .   "
	K    = "oo.oo "
	C    = " o.oo "
	M    = "  .ooo"
	SP   = "  .o  "
	_9   = "  . oo"
	_7   = "oo.o  "
)

var Message = []string{LTRS, W, FIGS, _6, LTRS, R, E, K, SP, C, M, FIGS, _9, _7}

func Circ(x, y, r, w int) {
	Printf("<circle cx='%d' cy='%d' r='%d' stroke='white' stroke-width='%d' fill='white' />\n", x, y, r, w)
}

func main() {
	Printf("<svg height='%d' width='%d'>\n", HEIGHT, WIDTH)
	Printf("<rect width='100%%' height='100%%' fill='black'/>\n")

	for i, code := range Message {
		for j := 0; j < 6; j++ {
			if code[j] == 'o' {
				Circ(100*i+400, 100*j+75, 35, 1)
			}
		}
	}

	for k := 100; k < len(Message)*100+700; k += 100 {
		Circ(k, 100*2+75, 20, 1)
	}

	Printf("<line x1='25' y1='2' x2='%d' y2='2' style='stroke:white;stroke-width:3' />\n", WIDTH-25)
	Printf("<line x1='25' y1='%d' x2='%d' y2='%d' style='stroke:white;stroke-width:3' />\n", HEIGHT-2, WIDTH-25, HEIGHT-2)

	Printf("</svg>\n")
}
