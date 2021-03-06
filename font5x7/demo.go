// +build main

/*
   Usage:    go run demo.go "The quick brown fox"

             go run demo.go | less
	        shows all 256 chars.
*/
package main

import (
	F "github.com/strickyak/qrss-squeak/font5x7"

	"fmt"
	"os"
	"strings"
)

var printf = fmt.Printf

func renderHorizontal(rows [][]bool) {
	for _, row := range rows {
		for _, bit := range row {
			if bit {
				printf("#")
			} else {
				printf("_")
			}
		}
		printf("\n")
	}
}

func renderVertical(bitmap []byte) {
	for _, row := range bitmap {
		var mask byte = 0x10
		for j := 0; j < 5; j++ {
			if 0 != (mask & row) {
				printf(" #")
			} else {
				printf(" _")
			}

			mask >>= 1
		}
		printf("\n")
	}
}

func main() {
	var s string
	if len(os.Args) >= 2 {
		s = strings.Join(os.Args[1:], " ") // String to print.
	} else {
		// Show all 256 chars, if no Args.
		for ch := 0; ch < 256; ch++ {
			s += string(ch)
		}
	}

	bitmap := F.VerticalStringFiveBitsWide(s)
	renderVertical(bitmap)

	rows := F.EightHorizontalRowsOfBool(s)
	renderHorizontal(rows)
}
