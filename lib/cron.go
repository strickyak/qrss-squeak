package lib

import (
	"log"
	. "time"
)

const pollsPerSecond = 10

var TEST_Scaling int64 = 1 // Always 1 except in speeded-up test.

type Cron struct {
	ModuloSeconds    int64
	RemainderSeconds int64
	Run              func()

	embargo int64
}

func (o *Cron) Start() {
	now := Now().UnixNano() * TEST_Scaling / 1000000000
	o.embargo = now + 1 // Don't trigger mid-second.
	pollTime := Second / Duration(TEST_Scaling) / pollsPerSecond
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				log.Fatalf("FATAL: Cron task %v failed: %v", o, e)
			}
		}()
		for {
			o.step()
			Sleep(pollTime)
		}
	}()
}

func (o *Cron) step() {
	now := Now().UnixNano() * TEST_Scaling / 1000000000
	if now < o.embargo {
		return
	}
	if now%o.ModuloSeconds != o.RemainderSeconds {
		return
	}

	o.embargo = now + 1 // Don't trigger again this second.

	o.Run()
}
