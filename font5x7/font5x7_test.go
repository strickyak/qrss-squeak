package font5x7

import . "testing"
import "reflect"

func ExpectDeepEqual(t *T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Logf("GOT : %#v", got)
		t.Logf("WANT: %#v", want)
		t.Error("FAIL, they are different")
	}
}

func TestVerticalStringFiveBitsWide(t *T) {
	got := VerticalStringFiveBitsWide("Hi.")
	want := []byte{0x11, 0x11, 0x11, 0x1f, 0x11, 0x11, 0x11, 0x0, 0x0, 0x4, 0x0, 0xc, 0x4, 0x4, 0x4, 0xe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6, 0x6, 0x0, 0x0}
	ExpectDeepEqual(t, got, want)
}

func TestEightHorizontalRowsOfBool(t *T) {
	got := EightHorizontalRowsOfBool("Jó napot!")
	want := [][]bool{
		[]bool{false, false, true, true, true, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, false, true, false, false, false, false, false, false, false, false, false, false, true, false, true, true, false, false, false, true, true, false, false, false, true, false, true, true, false, false, false, true, true, true, false, false, true, true, true, true, true, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, false, false, true, false, false, false, false, false, false, false, false, false, true, true, false, false, true, false, false, false, false, true, false, false, true, true, false, false, true, false, true, false, false, false, true, false, false, false, true, false, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, false, false, false, true, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, true, true, true, false, false, true, true, false, false, true, false, true, false, false, false, true, false, false, false, true, false, false, false, false, false, true, false, false, false},
		[]bool{true, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, true, false, false, true, false, false, true, false, true, true, false, false, true, false, false, false, true, false, false, false, true, false, true, false, false, false, false, false, false, false},
		[]bool{false, true, true, false, false, false, true, true, true, true, true, false, false, false, false, false, false, false, true, false, false, false, true, false, false, true, true, true, true, false, true, false, false, false, false, false, false, true, true, true, false, false, false, false, false, true, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	}
	ExpectDeepEqual(t, got, want)
}
