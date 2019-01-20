package lib

import (
	"testing"
	"time"
)

func Test3Seconds(t *testing.T) {
	i := 0
	p := &Cron{
		ModuloSeconds:    2,
		RemainderSeconds: 1,
		Run: func() {
			i++
			t.Log("running")
		}}
	p.Start()
	time.Sleep(3 * time.Second)
	if i < 1 || i > 2 {
		t.Errorf("got i = %d", i)
	}
	t.Log("okay: i =", i)
}
