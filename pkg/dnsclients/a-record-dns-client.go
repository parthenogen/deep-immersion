package dnsclients

import (
	"net"

	"github.com/miekg/dns"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

type aRecordDNSClient struct {
	client        *dns.Client
	serverUDPAddr string
}

func NewARecordDNSClient(clientAddr, serverAddr *net.UDPAddr) (
	c *aRecordDNSClient, e error,
) {
	c = &aRecordDNSClient{
		client: &dns.Client{
			Dialer: &net.Dialer{
				LocalAddr: clientAddr,
			},
		},
		serverUDPAddr: serverAddr.String(),
	}

	return
}

func (c *aRecordDNSClient) Send(query dimm.Query) (r dimm.Response, e error) {
	var (
		outgoing *dns.Msg = new(dns.Msg)
		incoming *dns.Msg
	)

	outgoing.SetQuestion(query.FQDN(),
		dns.TypeA,
	)

	incoming, _, e = c.client.Exchange(outgoing, c.serverUDPAddr)
	if e != nil {
		return
	}

	r = &response{incoming}

	return
}
