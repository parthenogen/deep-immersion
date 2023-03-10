package errorhandlers

import (
	"github.com/rs/zerolog/log"
)

type exitingErrorHandler struct {
}

func NewExitingErrorHandler() *exitingErrorHandler {
	return new(exitingErrorHandler)
}

func (h *exitingErrorHandler) Handle(e error) {
	log.Fatal().
		Err(e).
		Msg("Error encountered. Exiting with code 1.")
}
