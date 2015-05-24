package conductor

import (
	"syscall"

	"github.com/GeorgeMac/tools/conduit"
	"github.com/GeorgeMac/tools/wrappers"
)

var closed chan struct{}

func init() {
	closed = make(chan struct{})
	close(closed)
}

type Scheduler struct {
	t                Task
	timeout, between func() <-chan struct{}
	term             func()
}

func New(t Task, opts ...option) *Scheduler {
	s := &Scheduler{
		t:       t,
		between: func() <-chan struct{} { return closed },
	}

	// apply any other options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s Scheduler) Begin(notify chan<- struct{}) {
	term := make(chan struct{}, 1)
	if s.term != nil {
		wrappers.Notify(term, syscall.SIGTERM)
	}

	for {
		// channel to signal that Task t is complete
		done := make(chan struct{}, 0)
		// begin task
		go s.t.Run(done)

		// conduit for case where term recvs
		t := conduit.New()
		// conduit for either the task or timeout
		e := conduit.New()
		// block until t conduit recvs
		_, ok := <-t.Trigger(e.Either(done, s.timeout()), term)
		if !ok {
			// wait until e finishes, then send on notify
			s.terminate(notify, e)
			return
		}

		// otherwise either we're done or we timed out.

		// conduit for case where term recvs
		t = conduit.New()
		// block until either between sends to t, or term closes it.
		_, ok = <-t.Trigger(s.between(), term)
		if !ok {
			// send on notify immediately
			s.terminate(notify)
			return
		}
	}
}

// terminate is a useful function for common shutdown procedure
func (s *Scheduler) terminate(notify chan<- struct{}, after ...<-chan struct{}) {
	if s.term != nil {
		s.term()
	}
	for _, a := range after {
		<-a
	}
	close(notify)
}

// Task
type Task interface {
	Run(chan<- struct{})
}

// TaskFunc implements Task
// On a call to Run it calls the underlying function.
// When the function completes it clsoes the notify channel.
type TaskFunc func()

func (t TaskFunc) Run(notify chan<- struct{}) {
	t()
	close(notify)
}
