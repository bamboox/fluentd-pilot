package wait

import (
	"time"
	"math/rand"
)

var NeverStop <-chan struct{} = make(chan struct{})

func JitterUntil(f func(), period time.Duration, jitterFactor float64, sliding bool, stopCh <-chan struct{}) {
	for {

		select {
		case <-stopCh:
			return
		default:
		}

		jitteredPeriod := period
		if jitterFactor > 0.0 {
			jitteredPeriod = Jitter(period, jitterFactor)
		}

		var t *time.Timer
		if !sliding {
			t = time.NewTimer(jitteredPeriod)
		}

		func() {
			f()
		}()

		if sliding {
			t = time.NewTimer(jitteredPeriod)
		}

		// NOTE: b/c there is no priority selection in golang
		// it is possible for this to race, meaning we could
		// trigger t.C and stopCh, and t.C select falls through.
		// In order to mitigate we re-check stopCh at the beginning
		// of every loop to prevent extra executions of f().
		//
		select {
		case <-stopCh:
			return
		case <-t.C:
		}
	}
}
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	JitterUntil(f, period, 0.0, true, stopCh)
}
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64() * maxFactor * float64(duration))
	return wait
}
