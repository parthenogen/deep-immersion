package errorhandlers

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestExpectingErrorHandlerBefDelay(t *testing.T) {
	const (
		argument  = "-test.run=TestExpectingErrorHandlerBefDelay"
		envFormat = "%s=%s" // https://pkg.go.dev/os/exec#Cmd
		envKey    = "INCEPTION"
		envValue  = "1"
		exitCode  = 1
		delay     = time.Second
		quote     = "A wizard is never late, Frodo Baggins. Nor is he early."
	)

	var (
		handler *expectingErrorHandler

		command *exec.Cmd
		e       error
		exitErr *exec.ExitError
		ok      bool
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if os.Getenv(envKey) == envValue {
		handler = NewExpectingErrorHandler(delay)

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

func TestExpectingErrorHandlerAftDelay(t *testing.T) {
	const (
		argument  = "-test.run=TestExpectingErrorHandlerAftDelay"
		envFormat = "%s=%s" // https://pkg.go.dev/os/exec#Cmd
		envKey    = "INCEPTION"
		envValue  = "1"
		delay     = time.Second
		quote     = "He arrives precisely when he means to."
	)

	var (
		handler *expectingErrorHandler

		command *exec.Cmd
		e       error
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if os.Getenv(envKey) == envValue {
		handler = NewExpectingErrorHandler(delay)

		time.Sleep(delay)

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
	if e != nil {
		t.Error(e)
	}
}
