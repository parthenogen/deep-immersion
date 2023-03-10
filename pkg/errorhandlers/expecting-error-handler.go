package errorhandlers

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type expectingErrorHandler struct {
	timesUp <-chan time.Time
	handler func(error)
}

func NewExpectingErrorHandler(delay time.Duration) (h *expectingErrorHandler) {
	h = &expectingErrorHandler{
		timesUp: time.After(delay),
	}

	h.handler = h.HandleBefDelay

	return
}

func (h *expectingErrorHandler) Handle(e error) {
	h.handler(e)

	return
}

func (h *expectingErrorHandler) HandleBefDelay(e error) {
	const (
		exitCode = 1
	)

	select {
	case <-h.timesUp:
		h.handler = h.HandleAftDelay

		h.handler(e)

		return

	default:
		// do not wait for blocking channel
	}

	log.Fatal().
		Caller().
		Err(e).
		Msg("Error encountered earlier than expected. Exiting with code 1.")

	os.Exit(exitCode)
}

func (h *expectingErrorHandler) HandleAftDelay(e error) {
	const (
		exitCode = 0
	)

	log.Info().
		Caller().
		Err(e).
		Msg("Error encountered as expected. Exiting with code 0.")

	os.Exit(exitCode)
}
