package gopool

import (
"time"
"sync"
)


// TimerPool provides GC-able pooling of *time.Timer's.
// can be used by multiple goroutines concurrently.
type TimerPool struct {
	p sync.Pool
}

// Get returns a timer that completes after the given duration.
func (tp *TimerPool) Get(d time.Duration) *time.Timer {
	if t, _ := tp.p.Get().(*time.Timer); t != nil {
		t.Reset(d)
		return t
	}

	return time.NewTimer(d)
}

// Put pools the given timer.
//
// There is no need to call t.Stop() before calling Put.
//
// Put will try to stop the timer before pooling. If the
// given timer already expired, Put will read the unreceived
// value if there is one.
func (tp *TimerPool) Put(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}

	tp.p.Put(t)
}
