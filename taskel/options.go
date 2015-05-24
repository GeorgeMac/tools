package taskel

import (
	"syscall"
	"time"

	"github.com/GeorgeMac/tools/wrappers"
)

type option func(s *Scheduler)

// SetTimeout uses standard time.Ticker functionality
func SetTimeout(dur time.Duration) option {
	return func(s *Scheduler) {
		s.timeout = func() Trigger { return wrappers.StructTick(dur) }
	}
}

// SetTimeoutChannel allows you to pass a signalling
// channel of your own.
func SetTimeoutChannel(c Trigger) option {
	return func(s *Scheduler) {
		s.timeout = func() Trigger {
			return c
		}
	}
}

// SetTimeBetween allows you to define a wait period
// betwen running tasks.
func SetTimeBetween(dur time.Duration) option {
	return func(s *Scheduler) {
		s.between = func() Trigger { return wrappers.StructTick(dur) }
	}
}

// SetTerminateOn will trigger termination when it recvs
// on Trigger t.
func SetTerminateOn(t Trigger) option {
	return func(s *Scheduler) {
		s.term = t
	}
}

// SetTerminateOnSignal will trigger termination on the capture
// of a syscall.SIGTERM.
func SetTerminatOnSignal() option {
	return func(s *Scheduler) {
		ch := make(chan struct{})
		wrappers.Notify(ch, syscall.SIGTERM)
		s.term = ch
	}
}