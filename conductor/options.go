package conductor

import "time"

type option func(s *Scheduler)

func SetTimeout(dur time.Duration) option {
	return func(s *Scheduler) {
		s.timeout = dur
	}
}

func SetTimeBetween(dur time.Duration) option {
	return func(s *Scheduler) {
		s.between = dur
	}
}

func SetTerminate(fn func()) option {
	return func(s *Scheduler) {
		s.term = fn
	}
}
