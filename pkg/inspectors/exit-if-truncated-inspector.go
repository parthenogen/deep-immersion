package inspectors

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

type exitIfTruncatedInspector struct {
}

func NewExitIfTruncatedInspector() *exitIfTruncatedInspector {
	return new(exitIfTruncatedInspector)
}

func (i *exitIfTruncatedInspector) Inspect(response dimm.Response) {
	const (
		exitCode = 0
	)

	if response.Truncated() {
		log.Info().
			Caller().
			Msg("Truncated response observed. Exiting with code 0.")

		os.Exit(exitCode)
	}

	return
}
