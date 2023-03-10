package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

var (
	Commit = ""
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

	log.Debug().
		Str("source", "https://github.com/parthenogen/deep-immersion").
		Str("commit", Commit).
		Msg("Deep Immersion (Funky Orca)")

	log.Warn().
		Msg(
			"WARNING: Abuse, misuse, or incompetent use of this software " +
				"may result in practical and legal consequences. " +
				"Please ensure proper authorisation and containment " +
				"are in place before installation and execution.",
		)

	log.Warn().
		Msg(
			`THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, ` +
				`EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ` +
				`THE WARRANTIES OF MERCHANTABILITY, ` +
				`FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. ` +
				`IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS ` +
				`BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, ` +
				`WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ` +
				`ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE ` +
				`OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`,
		)

	config, e = newDriverConfig()
	if e != nil {
		log.Fatal().
			Err(e).
			Msg("Aborting due to bad configuration.")
	}

	driver = dimm.NewDriver(config)

	driver.Run()

	for {
	}
}
