// Package gopool contains tools for goroutine reuse.
package gopool

import (
	"fmt"
	"time"
)

var GlobalTimerpool Timerpool

// ErrScheduleTimeout returned by Grpool to indicate that there no free
// goroutines during some period of time.
var ErrScheduleTimeout = fmt.Errorf("schedule error: timed out")

// Grpool contains logic of goroutine reuse.
type Grpool struct {
	sem  chan struct{}
	work chan func()
}

// NewGrpool creates new goroutine pool with given size. It also creates a work
// queue of given size. Finally, it spawns given amount of goroutines
// immediately.
func NewGrpool(size, queue, spawn int) *Grpool {
	if spawn <= 0 && queue > 0 {
		panic("dead queue configuration detected")
	}
	if spawn > size {
		panic("spawn > workers")
	}
	p := &Grpool{
		sem:  make(chan struct{}, size),
		work: make(chan func(), queue),
	}
	for i := 0; i < spawn; i++ {
		p.sem <- struct{}{}
		go p.worker(func() {})
	}

	return p
}

// Schedule schedules task to be executed over pool's workers.
func (p *Grpool) Schedule(task func()) {
	p.schedule(task, nil)
}

// ScheduleTimeout schedules task to be executed over pool's workers.
// It returns ErrScheduleTimeout when no free workers met during given timeout.
func (p *Grpool) ScheduleTimeout(timeout time.Duration, task func()) error {
	return p.schedule(task, GlobalTimerpool.Get(timeout).C)
}

func (p *Grpool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.work <- task:
		return nil
	case p.sem <- struct{}{}:
		go p.worker(task)
		return nil
	}
}

func (p *Grpool) worker(task func()) {
	defer func() { <-p.sem }()

	task()

	for task := range p.work {
		task()
	}
}
