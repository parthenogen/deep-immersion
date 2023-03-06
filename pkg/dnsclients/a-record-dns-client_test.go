package dnsclients

import (
	"testing"
	"time"

	"github.com/miekg/dns"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

func TestARecordDNSClient(t *testing.T) {
	const (
		domainName = "example.org."
		serverAddr = "127.29.170.213:5353"
		timeout    = time.Second
	)

	var (
		client   *aRecordDNSClient
		server   *mockServer
		query    *stubQuery
		response dimm.Response
		e        error
	)

	server = newMockServer(serverAddr, domainName)

	defer server.shutdown()

	client, e = NewARecordDNSClient(serverAddr, timeout)
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

type mockServer struct {
	backend *dns.Server
	expects string
}

func newMockServer(address, expects string) (s *mockServer) {
	const (
		network = "udp"
	)

	s = &mockServer{
		backend: &dns.Server{Addr: address, Net: network},
		expects: expects,
	}

	s.backend.Handler = s

	go s.backend.ListenAndServe()

	return
}

func (s *mockServer) ServeDNS(writer dns.ResponseWriter, request *dns.Msg) {
	const (
		canned = "192.0.2.0"
	)

	var (
		response = new(dns.Msg)
	)

	response.SetReply(request)

	if request.Question[0].Name == s.expects {
		response.Truncated = true

		writer.WriteMsg(response)
	}

	return
}

func (s *mockServer) shutdown() {
	s.backend.Shutdown()

	return
}

type stubQuery struct {
	domainName string
}

func (q *stubQuery) FQDN() string {
	return q.domainName
}
