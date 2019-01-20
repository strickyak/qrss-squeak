package lib

import (
	"flag"
	"log"
	"strings"
)

var speed = flag.Float64("wpm", 12.0, "Words Per Minute")

var MorseCodeMap = map[rune][]DiDah{
	' ':  []DiDah(" "),
	'a':  []DiDah(".-"),
	'b':  []DiDah("-..."),
	'c':  []DiDah("-.-."),
	'd':  []DiDah("-.."),
	'e':  []DiDah("."),
	'f':  []DiDah("..-."),
	'g':  []DiDah("--."),
	'h':  []DiDah("...."),
	'i':  []DiDah(".."),
	'j':  []DiDah(".---"),
	'k':  []DiDah("-.-"),
	'l':  []DiDah(".-.."),
	'm':  []DiDah("--"),
	'n':  []DiDah("-."),
	'o':  []DiDah("---"),
	'p':  []DiDah(".--."),
	'q':  []DiDah("--.-"),
	'r':  []DiDah(".-."),
	's':  []DiDah("..."),
	't':  []DiDah("-"),
	'u':  []DiDah("..-"),
	'v':  []DiDah("...-"),
	'w':  []DiDah(".--"),
	'x':  []DiDah("-..-"),
	'y':  []DiDah("-.--"),
	'z':  []DiDah("--.."),
	'1':  []DiDah(".----"),
	'2':  []DiDah("..---"),
	'3':  []DiDah("...--"),
	'4':  []DiDah("....-"),
	'5':  []DiDah("....."),
	'6':  []DiDah("-...."),
	'7':  []DiDah("--..."),
	'8':  []DiDah("---.."),
	'9':  []DiDah("----."),
	'0':  []DiDah("-----"),
	'/':  []DiDah("-..-."),
	'.':  []DiDah(".-.-.-"),
	',':  []DiDah("--..--"),
	'?':  []DiDah("..--.."),
	'!':  []DiDah("-.-.--"),
	'=':  []DiDah("-...-"),
	'+':  []DiDah(".-.-."),
	'-':  []DiDah("-....-"),
	'(':  []DiDah("-.--."),
	')':  []DiDah("-.--.-"),
	':':  []DiDah("---..."),
	'@':  []DiDah(".--.-."),
	'&':  []DiDah(".-..."),
	'$':  []DiDah("...-..-"),
	'"':  []DiDah(".-..-."),
	'\'': []DiDah(".----."),
}

// Morse converts a string into morse code
// using '.' (1 dit time), '-' (3 dit times), and ' ' (1 dit time).
// It's fatal if we don't know one of the symbols.
func Morse(s string, tail bool) DiDahSlice {
	var z []DiDah
	n := len(s)
	for i, r := range strings.ToLower(s) {
		didahs, ok := MorseCodeMap[r]
		if !ok {
			log.Panicf("Morse code not known for rune: %d in %q", r, s)
		}
		for _, didah := range didahs {
			switch didah {
			case '.', '-', ' ':
				z = append(z, didah)
			default:
				log.Panicf("Bad char in morse code: %d in %q in %q", didah, didahs, s)
			}
		}
		if tail || i < n-1 {
			z = append(z, ' ')
		}
	}
	return z
}
