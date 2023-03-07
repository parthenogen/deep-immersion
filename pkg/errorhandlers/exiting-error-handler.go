package errorhandlers

import (
	"os"

	"github.com/rs/zerolog/log"
)

type exitingErrorHandler struct {
}

func NewExitingErrorHandler() *exitingErrorHandler {
	return new(exitingErrorHandler)
}

func (h *exitingErrorHandler) Handle(e error) {
	const (
		exitCode = 1
	)

	log.Fatal().
		Caller().
		Err(e).
		Msg("Error encountered. Exiting with code 1.")

	os.Exit(exitCode)
}
