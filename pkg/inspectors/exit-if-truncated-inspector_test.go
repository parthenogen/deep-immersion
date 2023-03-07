package inspectors

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestExitIfTruncatedInspector(t *testing.T) {
	const (
		argument  = "-test.run=TestExitIfTruncatedInspector"
		envFormat = "%s=%s" // https://pkg.go.dev/os/exec#Cmd
		envKey    = "INCEPTION"
		envValue  = "1"
		exitCode  = 0
	)

	var (
		inspector *exitIfTruncatedInspector

		command *exec.Cmd
		stub    stubResponse
		e       error
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if os.Getenv(envKey) == envValue {
		inspector = NewExitIfTruncatedInspector()

		stub.truncated = true

		inspector.Inspect(&stub)
	}

	command = exec.Command(os.Args[0],
		argument,
	)

	command.Stderr = os.Stderr

	command.Env = append(os.Environ(),
		fmt.Sprintf(envFormat, envKey, envValue),
	)

	e = command.Run()
	if e != nil {
		t.Fail()
	}
}

type stubResponse struct {
	truncated bool
}

func (r *stubResponse) Truncated() bool {
	return r.truncated
}
