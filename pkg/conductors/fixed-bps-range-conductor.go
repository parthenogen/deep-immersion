package conductors

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type fixedBPSRangeConductor struct {
	ticker *time.Ticker
	beats  chan struct{}
	stop   chan struct{}

	nBeats            uint
	checkInterval     time.Duration
	actualBPSLogLabel string
}

func NewFixedBPSRangeConductor(
	minBPS, maxBPS uint, checkInterval time.Duration, actualBPSLogLabel string,
) (
	c *fixedBPSRangeConductor,
) {
	c = &fixedBPSRangeConductor{
		ticker: time.NewTicker(
			time.Second / time.Duration(maxBPS),
		),
		beats: make(chan struct{}),
		stop:  make(chan struct{}),

		checkInterval:     checkInterval,
		actualBPSLogLabel: actualBPSLogLabel,
	}

	go c.run()

	go c.enforce(minBPS, checkInterval)

	return
}

func (c *fixedBPSRangeConductor) Beats() <-chan struct{} {
	return (<-chan struct{})(c.beats)
}

func (c *fixedBPSRangeConductor) Stop() {
	select {
	case <-c.beats:
	}

	close(c.stop)

	return
}

func (c *fixedBPSRangeConductor) run() {
	for {
		select {
		case <-c.stop:
			return

		case <-c.ticker.C:
			c.beats <- struct{}{}

			c.nBeats += 1
		}
	}
}

func (c *fixedBPSRangeConductor) enforce(minBPS uint, interval time.Duration) {
	const (
		exitCode = 1
	)

	var (
		bps uint
	)

	for {
		select {
		case <-c.stop:
			return

		default:
			bps = c.measure(interval)

			if bps < minBPS {
				log.Fatal().
					Uint(c.actualBPSLogLabel, bps).
					Uint("required", minBPS).
					Msgf("Failed to achieve minimum rate. "+
						"Exiting with code %d.",
						exitCode,
					)

				os.Exit(exitCode)
			}

			log.Info().
				Uint(c.actualBPSLogLabel, bps).
				Uint("count", c.nBeats).
				Msg("")
		}
	}
}

func (c *fixedBPSRangeConductor) measure(duration time.Duration) (bps uint) {
	const (
		rateConvFactor = time.Second / time.Nanosecond
	)

	var (
		nBeats0 uint
		nBeats1 uint
	)

	nBeats0 = c.nBeats

	time.Sleep(duration)

	nBeats1 = c.nBeats

	bps = uint(
		float64(nBeats1-nBeats0) / float64(duration) *
			float64(rateConvFactor),
	)

	return
}
