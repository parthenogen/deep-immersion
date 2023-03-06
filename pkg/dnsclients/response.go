package dnsclients

import (
	"github.com/miekg/dns"
)

type response struct {
	message *dns.Msg
}

func (r *response) Truncated() bool {
	return r.message.MsgHdr.Truncated
}
