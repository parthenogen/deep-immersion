package sources

import (
	"regexp"
	"testing"
)

const (
	uuidPattern = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
)

func TestGenerateUUID(t *testing.T) {
	var (
		random string
		valid  bool
	)

	random = generateUUID()

	valid, _ = regexp.MatchString(uuidPattern, random)

	if !valid {
		t.Fail()
	}
}

func TestUUIDSourceGenerateFQDN(t *testing.T) {
	const (
		domain  = "example.org."
		pattern = uuidPattern + `\.` + domain
	)

	var (
		source *uuidSource

		random string
		valid  bool
	)

	source = NewUUIDSource(domain)

	random = source.GenerateFQDN()

	valid, _ = regexp.MatchString(pattern, random)

	if !valid {
		t.Fail()
	}
}
