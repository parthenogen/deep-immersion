package sources

import (
	"fmt"

	"github.com/google/uuid"
)

type uuidSource struct {
	domain string
}

func NewUUIDSource(domain string) (s *uuidSource) {
	s = &uuidSource{
		domain: domain,
	}

	return
}

func (s *uuidSource) GenerateFQDN() (fqdn string) {
	const (
		format = "%s.%s"
	)

	fqdn = fmt.Sprintf(format,
		generateUUID(),
		s.domain,
	)

	return
}

func generateUUID() string {
	return uuid.New().String()
}
