package lib

import (
	"log"
	"time"
)

const pollsPerSecond = 10

type Cron struct {
	ModuloSeconds    uint64
	RemainderSeconds uint64
	Run              func()

	embargo uint64
}

func (o *Cron) Start() {
	now := uint64(time.Now().Unix())
	o.embargo = now + 1 // Don't trigger mid-second.
	tick := time.Second / pollsPerSecond
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				log.Fatalf("FATAL: Cron task %v failed: %v", o, e)
			}
		}()
		for {
			o.step()
			time.Sleep(tick)
		}
	}()
}

func (o *Cron) step() {
	now := uint64(time.Now().Unix())
	if now < o.embargo {
		return
	}
	if now%o.ModuloSeconds != o.RemainderSeconds {
		return
	}

	o.embargo = now + 1 // Don't trigger again this second.

	o.Run()
}
