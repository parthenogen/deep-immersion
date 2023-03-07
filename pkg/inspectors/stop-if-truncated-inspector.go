package inspectors

import (
	"github.com/rs/zerolog/log"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

type stoppable interface {
	Stop()
}

type stopIfTruncatedInspector struct {
	stoppable
}

func NewStopIfTruncatedInspector(s stoppable) (i *stopIfTruncatedInspector) {
	i = &stopIfTruncatedInspector{
		stoppable: s,
	}

	return
}

func (i *stopIfTruncatedInspector) Inspect(response dimm.Response) {
	if response.Truncated() {
		log.Info().
			Caller().
			Msg("Truncated response observed. Stopping now.")

		i.stoppable.Stop()
	}

	return
}
