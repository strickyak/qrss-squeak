package lib

import "testing"

var cases = []struct {
	in  string
	out []DiDah
}{
	{"abcdefg", []DiDah(".- -... -.-. -.. . ..-. --.")},
	{"hijklmnop", []DiDah(".... .. .--- -.- .-.. -- -. --- .--.")},
	{"qrstuv", []DiDah("--.- .-. ... - ..- ...-")},
	{"wxyz", []DiDah(".-- -..- -.-- --..")},
	{"12345", []DiDah(".---- ..--- ...-- ....- .....")},
	{"67890", []DiDah("-.... --... ---.. ----. -----")},
	{"Paris", []DiDah(".--. .- .-. .. ...")},
	{"VVV de", []DiDah("...- ...- ...-   -.. .")},
	{".,?/=+", []DiDah(".-.-.- --..-- ..--.. -..-. -...- .-.-.")},
}

func TestMorse(t *testing.T) {
	for i, e := range cases {
		got := Morse(e.in, false).String()
		if got != string(e.out) {
			t.Errorf("Test #%da in=%q got=%q want=%q", i, e.in, got, e.out)
		}
		got2 := Morse(e.in, true).String()
		if got2 != got+" " {
			t.Errorf("Test #%db in=%q got2=%q want=%q", i, e.in, got2, got+" ")
		}
	}
}
