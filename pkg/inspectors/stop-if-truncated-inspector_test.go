package inspectors

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestStopIfTruncatedInspector(t *testing.T) {
	var (
		inspector *stopIfTruncatedInspector
		spy       spyStoppable
		stub      stubResponse
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
