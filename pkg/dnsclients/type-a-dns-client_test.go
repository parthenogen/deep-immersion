package dnsclients

import (
	"net"
	"net/netip"
	"testing"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
	"github.com/parthenogen/deep-immersion/pkg/dnsservers"
)

func TestTypeADNSClient(t *testing.T) {
	const (
		domainName = "example.org."
		clientAddr = "127.171.180.45:35353"
		serverAddr = "127.29.170.213:5353"

		truncateAfter = 0
	)

	var (
		client   *typeADNSClient
		server   stoppable
		query    *stubQuery
		response dimm.Response
		e        error
	)

	server, e = dnsservers.NewMockDNSServer(
		net.UDPAddrFromAddrPort(netip.MustParseAddrPort(serverAddr)),
		clientAddr,
		domainName,
		truncateAfter,
	)
	if e != nil {
		t.Error(e)
	}

	defer server.Stop()

	client, e = NewTypeADNSClient(
		net.UDPAddrFromAddrPort(netip.MustParseAddrPort(clientAddr)),
		net.UDPAddrFromAddrPort(netip.MustParseAddrPort(serverAddr)),
	)
	if e != nil {
		t.Error(e)
	}

	query = &stubQuery{domainName}

	response, e = client.Send(query)
	if e != nil {
		t.Error(e)

	} else if !response.Truncated() {
		t.Fail()
	}
}

type stubQuery struct {
	domainName string
}

func (q *stubQuery) FQDN() string {
	return q.domainName
}

type stoppable interface {
	Stop()
}
