package utils

import "time"

// Timer is an interface to be able to mock time.Timer in the unit tests
type Timer interface {
	Stop()
	C() <-chan time.Time
}

type TimerWrapper struct {
	timer *time.Timer
}

func NewTimerWrapper(d time.Duration) Timer {
	return &TimerWrapper{
		timer: time.NewTimer(d),
	}
}

func (t *TimerWrapper) Stop() {
	t.timer.Stop()
}

func (t *TimerWrapper) C() <-chan time.Time {
	return t.timer.C
}
