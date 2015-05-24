package taskel

import "testing"

var errmsg string = "Expected %v, Got %v\n"

func Test_Scheduler_New(t *testing.T) {
	var called bool
	fn := func() {
		called = true
	}

	s := New(TaskFunc(fn))
	// check nil was not returned
	if s == nil {
		t.Fatalf(errmsg, "not nil", s)
	}

	// run scheduler task and check it handles notification
	// and function delegation correctly.
	ch := make(chan struct{})
	s.t.Run(ch)

	// ensure called is now set to true
	if !called {
		t.Fatalf(errmsg, true, called)
	}

	// ensure we don't block on notify channel `ch`
	select {
	case <-ch:
		return
	default:
		t.Fatalf(errmsg, "to not block", "blocked")
	}
}
