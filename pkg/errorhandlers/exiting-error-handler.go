package errorhandlers

import (
	"os"
)

type exitingErrorHandler struct {
}

func NewExitingErrorHandler() *exitingErrorHandler {
	return new(exitingErrorHandler)
}

func (h *exitingErrorHandler) Handle(error) {
	const (
		exitCode = 1
	)

	os.Exit(exitCode)
}
