package conductor

import (
	"time"

	"github.com/GeorgeMac/tools/wrappers"
)

type option func(s *Scheduler)

// SetTimeout uses standard time.Ticker functionality
func SetTimeout(dur time.Duration) option {
	return func(s *Scheduler) {
		s.timeout = func() <-chan struct{} { return wrappers.StructTick(dur) }
	}
}

// SetTimeoutChannel allows you to pass a signalling
// channel of your own.
func SetTimeoutChannel(c <-chan struct{}) option {
	return func(s *Scheduler) {
		s.timeout = func() <-chan struct{} {
			return c
		}
	}
}

// SetTimeBetween allows you to define a wait period
// betwen running tasks.
func SetTimeBetween(dur time.Duration) option {
	return func(s *Scheduler) {
		s.between = func() <-chan struct{} { return wrappers.StructTick(dur) }
	}
}

func SetTerminate(fn func()) option {
	return func(s *Scheduler) {
		s.term = fn
	}
}
