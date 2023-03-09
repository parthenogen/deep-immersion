package dnsclients

import (
	"net"

	"github.com/miekg/dns"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

type typeADNSClient struct {
	client        *dns.Client
	serverUDPAddr string
}

func NewTypeADNSClient(clientAddr, serverAddr *net.UDPAddr) (
	c *typeADNSClient, e error,
) {
	c = &typeADNSClient{
		client: &dns.Client{
			Dialer: &net.Dialer{
				LocalAddr: clientAddr,
			},
		},
		serverUDPAddr: serverAddr.String(),
	}

	return
}

func (c *typeADNSClient) Send(query dimm.Query) (r dimm.Response, e error) {
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
