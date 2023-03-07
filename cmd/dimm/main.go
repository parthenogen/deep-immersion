package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

func main() {
	const (
		exitCode = 1
	)

	var (
		config *driverConfig
		driver interface {
			Run()
			Stop()
		}
		e error
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config, e = newDriverConfig()
	if e != nil {
		log.Fatal().
			Caller().
			Err(e).
			Msg("Aborting due to bad configuration.")

		os.Exit(exitCode)
	}

	driver = dimm.NewDriver(config)

	log.Debug().
		Caller().
		Msg("Deep Immersion initialised. Now running.")

	driver.Run()

	for {
	}
}
