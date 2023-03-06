package conductors

import (
	"time"

	"github.com/rs/zerolog/log"
)

type maxBeatsPerSecConductor struct {
	ticker *time.Ticker
	beats  chan struct{}
	stop   chan struct{}

	logInterval       uint
	actualBPSLogLabel string

	nBeats uint
	tStart time.Time
}

func NewMaxBeatsPerSecConductor(
	maxBPS, logInterval uint, actualBPSLogLabel string,
) (
	c *maxBeatsPerSecConductor,
) {
	c = &maxBeatsPerSecConductor{
		ticker: time.NewTicker(
			time.Second / time.Duration(maxBPS),
		),
		beats: make(chan struct{}),
		stop:  make(chan struct{}),

		logInterval:       logInterval,
		actualBPSLogLabel: actualBPSLogLabel,
	}

	go c.run()

	return
}

func (c *maxBeatsPerSecConductor) Beats() chan struct{} {
	return c.beats
}

func (c *maxBeatsPerSecConductor) Stop() {
	select {
	case <-c.beats:
	}

	close(c.stop)

	return
}

func (c *maxBeatsPerSecConductor) run() {
	const (
		rateConvFactor = float64(time.Second / time.Nanosecond)
	)

	var (
		t time.Time
	)

	for {
		select {
		case <-c.stop:
			return

		case t = <-c.ticker.C:
			if c.nBeats == 0 {
				c.tStart = t
			}

			c.beats <- struct{}{}

			c.nBeats += 1

			if c.nBeats%c.logInterval == 0 {
				log.Info().
					Caller().
					Uint(c.actualBPSLogLabel,
						uint(
							float64(c.nBeats)/float64(t.Sub(c.tStart))*
								rateConvFactor,
						),
					).
					Msg("")
			}
		}
	}
}
