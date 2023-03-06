package inspectors

import (
	"testing"
)

func TestStopIfTruncatedInspector(t *testing.T) {
	var (
		inspector *stopIfTruncatedInspector
		mock      mockStoppable
		dummy     dummyResponse
	)

	inspector = NewStopIfTruncatedInspector(&mock)

	if mock.stopped {
		t.Fail()
	}

	inspector.Inspect(&dummy)

	if mock.stopped {
		t.Fail()
	}

	dummy.truncated = true

	inspector.Inspect(&dummy)

	if !mock.stopped {
		t.Fail()
	}
}

type mockStoppable struct {
	stopped bool
}

func (s *mockStoppable) Stop() {
	s.stopped = true

	return
}

type dummyResponse struct {
	truncated bool
}

func (r *dummyResponse) Truncated() bool {
	return r.truncated
}
