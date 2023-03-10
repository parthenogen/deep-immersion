package dnsservers

import (
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type message struct {
	dnsmessage.Message
	client *net.UDPAddr
}
