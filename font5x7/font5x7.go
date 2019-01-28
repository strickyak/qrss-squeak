/*
	Package font5x7 provides a simple 5x8 bitmap.
*/
package font5x7

import "log"

/*
EightHorizontalRowsOfBool returns a slice of length 8 where each element is a slice of bool.
The outer slice goes down, and the inner slices go across, to print like this (with '#' for true and '_' for false):

#_____________#____________________________________________________##_____________________________#___
#_____________#_____________________________________________________#_____________________________#___
#______###__#####_______##_#__#___#_______#_##___###___###__#_##____#____###_________###___###____#___
#_____#___#___#_________#_#_#_#___#_______##__#_#___#_#___#_##__#___#___#___#_______#__##_#___#___#___
#_____#####___#_________#_#_#__####_______##__#_#####_#___#_##__#___#___#####_______#__##_#___#___#___
#_____#_______#_#_______#_#_#_____#_______#_##__#_____#___#_#_##____#___#____________##_#_#___#_______
#####__###_____#________#_#_#_#___#_______#______###___###__#______###___###____________#__###____#___
_______________________________###________#_________________#________________________###______________

*/
func EightHorizontalRowsOfBool(s string) [][]bool {
	var z [][]bool
	for row := 0; row < 8; row++ { // For row in 0..6 inclusive
		var y []bool
		for _, ascii := range s {
			if ascii > 255 {
				ascii = '?' // Replace non-printable-ascii with '?'
			}
			for col := 0; col < 5; col++ {
				rowBits := Font[5*(int(ascii))+col]
				y = append(y, bool(0 != ((1<<uint(row))&rowBits)))
			}
			y = append(y, false)
		}
		z = append(z, y)
	}
	return z
}

// VerticalStringFiveBitsWide converts the ASCII string into bytes for a vertical banner
// with the letters written on top of each other going downward.
// See demo.go for an example of how to consume them.
func VerticalStringFiveBitsWide(s string) []byte {
	var z []byte
	for _, ascii := range s {
		if ascii > 255 {
			ascii = '?' // Replace non-printable-ascii with '?'
		}
		for row := 0; row < 8; row++ {
			var x byte
			for col := 0; col < 5; col++ {
				x = (x << 1)
				if Pixel(byte(ascii), row, col) {
					x |= 1
				}
			}
			z = append(z, x)
		}
		z = append(z, byte(0)) // Empty row after 7 rows of char pixels.
	}
	return z
}

func Pixel(ascii byte, row int, col int) bool {
	if row < 0 || row > 7 {
		log.Panicf("Expected 0 <= row <= 7; got row=%d", row)
	}
	if col < 0 || col > 4 {
		log.Panicf("Expected 0 <= col <= 4; got col=%d", col)
	}
	if ascii > 255 {
		return 0 != ((row & 1) ^ (col & 1)) // Checkerboard for bad ascii.
	}
	rowBits := Font[5*(int(ascii))+col]
	return 0 != ((rowBits >> uint(row)) & 1)
}
