package dnsclients

import (
	"time"

	"github.com/miekg/dns"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

type aRecordDNSClient struct {
	client        *dns.Client
	serverUDPAddr string
	timeout       time.Duration
}

func NewARecordDNSClient(serverUDPAddr string, timeout time.Duration) (
	c *aRecordDNSClient, e error,
) {
	c = &aRecordDNSClient{
		client:        new(dns.Client),
		serverUDPAddr: serverUDPAddr,
		timeout:       timeout,
	}

	return
}

func (c *aRecordDNSClient) Send(query dimm.Query) (r dimm.Response, e error) {
	var (
		outgoing *dns.Msg = new(dns.Msg)
		incoming *dns.Msg

		timeIsUp <-chan time.Time
	)

	outgoing.SetQuestion(query.FQDN(),
		dns.TypeA,
	)

	timeIsUp = time.After(c.timeout)

	for {
		select {
		case <-timeIsUp:
			return

		default:
			incoming, _, e = c.client.Exchange(outgoing, c.serverUDPAddr)
		}

		if e == nil {
			break
		}
	}

	r = &response{incoming}

	return
}
