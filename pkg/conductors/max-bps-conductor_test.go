package conductors

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestMaxBeatsPerSecConductor(t *testing.T) {
	const (
		maxBPS      = 4096
		logInterval = 256
		logLabel    = "actualQPS"
	)

	var (
		conductor *maxBeatsPerSecConductor
		nBeats    uint
		nBeatsFin uint
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conductor = NewMaxBeatsPerSecConductor(maxBPS, logInterval, logLabel)

	go countBeats(
		&nBeats,
		conductor.Beats(),
	)

	time.Sleep(time.Second)

	conductor.Stop()

	nBeatsFin = nBeats

	if nBeatsFin < 1 || nBeatsFin > maxBPS {
		t.Fail()
	}

	nBeats = 0

	go countBeats(
		&nBeats,
		conductor.Beats(),
	)

	time.Sleep(time.Second)

	if nBeats != 0 {
		t.Fail()
	}
}

func countBeats(counter *uint, beats <-chan struct{}) {
	for {
		select {
		case <-beats:
			*counter += 1
		}
	}
}
