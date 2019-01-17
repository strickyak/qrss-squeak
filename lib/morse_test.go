package lib

import "testing"

var cases = []struct {
	in  string
	out DiDah
}{
	{"abcdefg", ".- -... -.-. -.. . ..-. --."},
	{"hijklmnop", ".... .. .--- -.- .-.. -- -. --- .--."},
	{"qrstuv", "--.- .-. ... - ..- ...-"},
	{"wxyz", ".-- -..- -.-- --.."},
	{"12345", ".---- ..--- ...-- ....- ....."},
	{"67890", "-.... --... ---.. ----. -----"},
	{"Paris", ".--. .- .-. .. ..."},
	{"VVV de", "...- ...- ...-   -.. ."},
	{".,?/=+", ".-.-.- --..-- ..--.. -..-. -...- .-.-."},
}

func TestMorse(t *testing.T) {
	for i, e := range cases {
		got := Morse(e.in, false)
		if got != e.out {
			t.Errorf("Test #%da in=%q got=%q want=%q", i, e.in, got, e.out)
		}
		got2 := Morse(e.in, true)
		if got2 != got+" " {
			t.Errorf("Test #%db in=%q got2=%q want=%q", i, e.in, got2, got+" ")
		}
	}
}
