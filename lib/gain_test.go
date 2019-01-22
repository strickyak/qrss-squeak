package lib

import "testing"

func TestGain(t *testing.T) {
	o := &Gain{
		Gain: 3.14,
		Input: &Slice{
			Slice: []float64{1 * 1, 2 * 2, 3 * 3},
		},
	}
	ch := make(chan Volt, 0)
	go func() {
		o.Emit(ch)
		close(ch)
	}()
	for i := 1; i <= 3; i++ {
		x, ok := <-ch
		if !ok || x != 3.14*Volt(i*i) {
			t.Errorf("fail %v %v", x, ok)
		}
	}
	x, ok := <-ch
	if ok {
		t.Errorf("fail2 %v %v", x, ok)
	}
}
