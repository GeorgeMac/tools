package wrappers

import "time"

type StructTicker struct {
	t *time.Ticker
	C <-chan struct{}
}

func NewStructTicker(dur time.Duration) StructTicker {
	c := make(chan struct{})
	s := StructTicker{
		t: time.NewTicker(dur),
		C: c,
	}
	go func() {
		for range s.t.C {
			c <- struct{}{}
		}
	}()
	return s
}

func (s StructTicker) Stop() {
	s.t.Stop()
}

func StructTick(dur time.Duration) <-chan struct{} {
	return NewStructTicker(dur).C
}
