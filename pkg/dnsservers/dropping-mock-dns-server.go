package dnsservers

import (
	"net"
)

type droppingMockDNSServer struct {
	*mockDNSServer
	queriesHandled uint
}

func NewDroppingMockDNSServer(address *net.UDPAddr, dropAfterQuery uint) (
	s *droppingMockDNSServer, e error,
) {
	var (
		handler func(message) message
	)

	s = new(droppingMockDNSServer)

	handler = newDroppingHandler(dropAfterQuery, &s.queriesHandled)

	s.mockDNSServer, e = newMockDNSServer(address, handler)
	if e != nil {
		return
	}

	return
}

func newDroppingHandler(dropAfterQuery uint, nQueriesHandled *uint) (
	handler func(message) message,
) {
	handler = func(incoming message) (outgoing message) {
		if *nQueriesHandled >= dropAfterQuery {
			return
		}

		outgoing.Header.ID = incoming.Header.ID

		outgoing.client = incoming.client

		*nQueriesHandled += 1

		return
	}

	return
}
