package lib

import (
	"log"
	"time"
)

func SecondsToDuration(secs float64) time.Duration {
	return time.Duration(1000000000*secs) * time.Nanosecond
}

func AdjustDuration(em Emitter, newDur time.Duration) {
	oldDit := *em.DitPtr()
	oldDur := em.Duration()

	var ratio float64 = newDur.Seconds() / oldDur.Seconds()
	var newDitSec float64 = ratio * oldDit.Seconds()
	var newDit time.Duration = time.Duration(1000000000*newDitSec) * time.Nanosecond
	*em.DitPtr() = newDit

	log.Printf("AdjustDuration: %v old: dit=%v dur=%v new: dit=%v dur=%v: %v", newDur, oldDit, oldDur, newDit, em.Duration(), em)
}
