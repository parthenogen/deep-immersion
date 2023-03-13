package dnsclients

import (
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"

	"github.com/parthenogen/deep-immersion/pkg/dimm"
)

type typeADNSClient struct {
	client        *dns.Client
	serverUDPAddr string

	ticker *time.Ticker
}

func NewTypeADNSClient(clientAddr, serverAddr *net.UDPAddr,
	checkInterval time.Duration,
) (
	c *typeADNSClient, e error,
) {
	c = &typeADNSClient{
		client: &dns.Client{
			Dialer: &net.Dialer{
				LocalAddr: clientAddr,
			},
		},
		serverUDPAddr: serverAddr.String(),
		ticker:        time.NewTicker(checkInterval),
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

	select {
	case <-c.ticker.C:
		log.Debug().
			Str("client.addr", c.client.Dialer.LocalAddr.String()).
			Uint16("query.id", outgoing.MsgHdr.Id).
			Str("query.name", outgoing.Question[0].Name).
			Msg("Sampled outgoing DNS query.")

	default:
	}

	return
}
