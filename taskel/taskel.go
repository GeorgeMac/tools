package taskel

import "github.com/GeorgeMac/tools/conduit"

var closed chan struct{}

func init() {
	closed = make(chan struct{})
	close(closed)
}

type Trigger <-chan struct{}

type Scheduler struct {
	t                Task
	timeout, between func() Trigger
	term             Trigger
}

// New returns a new instance of the Scheduler, with options
// applied.
// Just calling taskel.New(t) will return a Scheduler, which on a call
// to Begin will start an infinite for loop which calls t.Run(...).
func New(t Task, opts ...option) *Scheduler {
	s := &Scheduler{
		t:       t,
		timeout: func() Trigger { return nil },
		between: func() Trigger { return closed },
	}

	// apply any other options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Begin
func (s Scheduler) Begin(notify chan<- struct{}) {
	for {
		// channel to signal that Task t is complete
		done := make(chan struct{})
		// begin task
		go s.t.Run(done)

		// conduit for case where term recvs
		t := conduit.New()
		// conduit for either the task or timeout
		e := conduit.New()
		// block until t conduit recvs
		_, ok := <-t.Trigger(e.Either(done, s.timeout()), s.term)
		if !ok {
			// wait until e finishes, then send on notify
			terminate(notify, e)
			return
		}

		// otherwise either we're done or we timed out.

		// conduit for case where term recvs
		t = conduit.New()
		// block until either between sends to t, or term closes it.
		_, ok = <-t.Trigger(s.between(), s.term)
		if !ok {
			// send on notify immediately
			terminate(notify)
			return
		}
	}
}

// terminate is a useful function for common shutdown procedure
func terminate(notify chan<- struct{}, after ...<-chan struct{}) {
	// block on all channels in after
	for _, a := range after {
		<-a
	}
	close(notify)
}

// Task
type Task interface {
	Run(chan<- struct{})
}

// TaskFunc implements Task.
// On a call to Run it calls the underlying function.
type TaskFunc func(chan<- struct{})

// Run calls the underlying TaskFunc with the provided notify channel.
func (t TaskFunc) Run(notify chan<- struct{}) {
	t(notify)
}

// AfterTaskFunc implements Task.
// On a call to Run it calls the underlying function.
// When the function completes it closes the notify channel.
type AfterTaskFunc func()

// Run calls the underlying func. Once the call finishes
// it closes the provided notify channel.
func (t AfterTaskFunc) Run(notify chan<- struct{}) {
	t()
	close(notify)
}
