package lib

import (
	"testing"
	. "time"
)

func Test3Seconds(t *testing.T) {
	TEST_Scaling = 5 // Speed up the test 5x.
	i := 0
	p := &Cron{
		ModuloSeconds:    2,
		RemainderSeconds: 1,
		Run: func() {
			i++
			t.Log("running")
		}}
	p.Start()
	Sleep(3 * Second / Duration(TEST_Scaling))
	if i < 1 || i > 2 {
		t.Errorf("got i = %d", i)
	} else {
		t.Log("okay: i =", i)
	}
}
