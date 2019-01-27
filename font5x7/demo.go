// +build main

/*
   Usage:    go run demo.go "The quick brown fox" | less
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
	s := "Let my people go!" // String to print.
	if len(os.Args) >= 2 {
		s = strings.Join(os.Args[1:], " ") // String to print.
	}

	bitmap := F.VerticalStringFiveBitsWide(s)
	renderVertical(bitmap)

	rows := F.SevenHorizontalRowsOfBool(s)
	renderHorizontal(rows)
}
