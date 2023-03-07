package errorhandlers

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestExitingErrorHandler(t *testing.T) {
	const (
		argument  = "-test.run=TestExitingErrorHandler"
		envFormat = "%s=%s" // https://pkg.go.dev/os/exec#Cmd
		envKey    = "INCEPTION"
		envValue  = "1"
		exitCode  = 1
		quote     = "That many dreams within dreams is too unstable."
	)

	var (
		handler *exitingErrorHandler

		command *exec.Cmd
		e       error
		exitErr *exec.ExitError
		ok      bool
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if os.Getenv(envKey) == envValue {
		handler = NewExitingErrorHandler()

		handler.Handle(
			fmt.Errorf(quote),
		)
	}

	command = exec.Command(os.Args[0],
		argument,
	)

	command.Stderr = os.Stderr

	command.Env = append(os.Environ(),
		fmt.Sprintf(envFormat, envKey, envValue),
	)

	e = command.Run()

	exitErr, ok = e.(*exec.ExitError)
	if !ok {
		t.Fail()
	}

	if exitErr.ExitCode() != exitCode {
		t.Fail()
	}
}
