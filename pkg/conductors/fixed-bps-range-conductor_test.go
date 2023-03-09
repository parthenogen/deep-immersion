package conductors

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestFixedBPSRangeConductor(t *testing.T) {
	const (
		minBPS      = 1 << 6
		maxBPS      = 1 << 8
		logInterval = 100 * time.Millisecond
		logLabel    = "qps"
	)

	var (
		conductor *fixedBPSRangeConductor
		nBeats    uint
		nBeatsFin uint
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conductor = NewFixedBPSRangeConductor(minBPS, maxBPS, logInterval, logLabel)

	go countBeats(
		&nBeats,
		conductor.Beats(),
	)

	time.Sleep(time.Second)

	conductor.Stop()

	nBeatsFin = nBeats

	if nBeatsFin < minBPS || nBeatsFin > maxBPS {
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
