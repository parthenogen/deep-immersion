package inspectors

import (
	"testing"
)

func TestStopIfTruncatedInspector(t *testing.T) {
	var (
		inspector *stopIfTruncatedInspector
		spy       spyStoppable
		stub      stubResponse
	)

	inspector = NewStopIfTruncatedInspector(&spy)

	if spy.stopped {
		t.Fail()
	}

	inspector.Inspect(&stub)

	if spy.stopped {
		t.Fail()
	}

	stub.truncated = true

	inspector.Inspect(&stub)

	if !spy.stopped {
		t.Fail()
	}
}

type spyStoppable struct {
	stopped bool
}

func (s *spyStoppable) Stop() {
	s.stopped = true

	return
}

type stubResponse struct {
	truncated bool
}

func (r *stubResponse) Truncated() bool {
	return r.truncated
}
