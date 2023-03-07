package dnsclients

import (
	"net"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

type typeADNSClient struct {
	client        *dns.Client
	serverUDPAddr string

	counter     uint
	logInterval uint
}

func NewTypeADNSClient(clientAddr, serverAddr *net.UDPAddr, logInterval uint) (
	c *typeADNSClient, e error,
) {
	c = &typeADNSClient{
		client: &dns.Client{
			Dialer: &net.Dialer{
				LocalAddr: clientAddr,
			},
		},
		serverUDPAddr: serverAddr.String(),
		logInterval:   logInterval,
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

	c.counter += 1

	if c.counter%c.logInterval == 0 {
		log.Debug().
			Caller().
			Uint("counter", c.counter).
			Str("query.name", outgoing.Question[0].Name).
			Bool("response.truncated", incoming.MsgHdr.Truncated). // XXX *
			Msg("Query sent; response received.")

		// * highly specific
	}

	return
}
