package dnsservers

import (
	"net"
	"strings"
)

type truncatingMockDNSServer struct {
	*mockDNSServer
	queriesHandled uint
}

func NewTruncatingMockDNSServer(
	address *net.UDPAddr, expectedClientCIDR, expectedDomainName string,
	truncateAfterQuery uint,
) (
	s *truncatingMockDNSServer, e error,
) {
	var (
		handler func(message) message
	)

	s = new(truncatingMockDNSServer)

	handler = newTruncatingHandler(
		expectedClientCIDR,
		expectedDomainName,
		truncateAfterQuery,
		&s.queriesHandled,
	)

	s.mockDNSServer, e = newMockDNSServer(address, handler)
	if e != nil {
		return
	}

	return
}

func newTruncatingHandler(expectedClientCIDR, expectedDomainName string,
	truncateAfterQuery uint, nQueriesHandled *uint,
) (
	handler func(message) message,
) {
	var (
		expectedClientIPNet *net.IPNet
		e                   error
	)

	_, expectedClientIPNet, e = net.ParseCIDR(expectedClientCIDR)
	if e != nil {
		return
	}

	handler = func(incoming message) (outgoing message) {
		var (
			clientAddrOK bool
			domainNameOK bool
		)

		clientAddrOK = expectedClientIPNet.Contains(
			incoming.client.IP,
		)

		domainNameOK = strings.HasSuffix(
			incoming.Questions[0].Name.String(),
			expectedDomainName,
		)

		if clientAddrOK && domainNameOK {
			outgoing.Header.ID = incoming.Header.ID

			if *nQueriesHandled >= truncateAfterQuery {
				outgoing.Header.Truncated = true
			}
		}

		outgoing.client = incoming.client

		*nQueriesHandled += 1

		return
	}

	return
}
