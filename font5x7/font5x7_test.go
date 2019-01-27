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
	want := []byte("\x11\x11\x11\x1f\x11\x11\x11\x00\x00\x04\x00\f\x04\x04\x04\x0e\x00\x00\x00\x00\x00\x00\x00\f\f\x00\x00")
	ExpectDeepEqual(t, got, want)
}

func TestSevenHorizontalRowsOfBool(t *T) {
	got := SevenHorizontalRowsOfBool("Jó napot!")
	want := [][]bool{
		[]bool{false, false, true, true, true, false, false, true, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, true, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, false, false, false, false, true, false, false, false, false, false, false, false, true, false, true, true, false, false, false, true, true, true, false, false, true, true, true, true, false, false, false, true, true, true, false, false, true, true, true, false, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, false, false, false, true, false, false, false, false, false, false, false, false, true, true, false, false, true, false, false, false, false, false, true, false, true, false, false, false, true, false, true, false, false, false, true, false, false, true, false, false, false, false, false, false, true, false, false, false},
		[]bool{false, false, false, true, false, false, false, false, true, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, true, true, true, true, false, true, true, true, true, false, false, true, false, false, false, true, false, false, true, false, false, false, false, false, false, true, false, false, false},
		[]bool{true, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, true, false, false, false, true, false, true, false, false, false, false, false, true, false, false, false, true, false, false, true, false, false, true, false, false, false, false, false, false, false},
		[]bool{false, true, true, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, true, true, true, true, false, true, false, false, false, false, false, false, true, true, true, false, false, false, false, true, true, false, false, false, false, true, false, false, false},
	}
	ExpectDeepEqual(t, got, want)
}
