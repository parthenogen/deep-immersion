package dnsservers

import (
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

const (
	udpMsgSizeUppLim = 1 << 16
)

type mockDNSServer struct {
	udpConn  *net.UDPConn
	incoming chan message
	outgoing chan message
	stop     chan struct{}

	expectedClientAddr string
	expectedDomainName string
}

func NewMockDNSServer(
	address *net.UDPAddr, expectedClientAddr, expectedDomainName string,
) (
	s *mockDNSServer, e error,
) {
	const (
		network = "udp"
	)

	s = &mockDNSServer{
		incoming:           make(chan message),
		outgoing:           make(chan message),
		stop:               make(chan struct{}),
		expectedClientAddr: expectedClientAddr,
		expectedDomainName: expectedDomainName,
	}

	s.udpConn, e = net.ListenUDP(network, address)
	if e != nil {
		return
	}

	go s.runReader()
	go s.runWriter()
	go s.runHandler()

	return
}

func (s *mockDNSServer) Stop() {
	close(s.stop)

	return
}

func (s *mockDNSServer) runReader() {
	var (
		b = make([]byte, udpMsgSizeUppLim)

		msg message
		e   error
	)

	for {
		select {
		case <-s.stop:
			return

		default:
			_, msg.client, e = s.udpConn.ReadFromUDP(b)
			if e != nil {
				panic(e) // only for use in testing
			}

			e = msg.Unpack(b)
			if e != nil {
				panic(e)
			}

			s.incoming <- msg
		}
	}
}

func (s *mockDNSServer) runWriter() {
	var (
		msg message
		b   []byte
		e   error
	)

	for {
		select {
		case <-s.stop:
			return

		default:
			msg = <-s.outgoing

			b, e = msg.Pack()
			if e != nil {
				panic(e)
			}

			_, e = s.udpConn.WriteToUDP(b, msg.client)
			if e != nil {
				panic(e)
			}
		}
	}
}

func (s *mockDNSServer) runHandler() {
	var (
		incoming message
		outgoing message

		clientAddrOK bool
		domainNameOK bool
	)

	for {
		select {
		case <-s.stop:
			return

		default:
			incoming = <-s.incoming

			clientAddrOK = (incoming.client.String() == s.expectedClientAddr)

			domainNameOK = (incoming.Questions[0].Name.String() ==
				s.expectedDomainName)

			if clientAddrOK && domainNameOK {
				outgoing.Header.ID = incoming.Header.ID

				outgoing.Header.Truncated = true
			}

			outgoing.client = incoming.client

			s.outgoing <- outgoing
		}
	}
}

type message struct {
	dnsmessage.Message
	client *net.UDPAddr
}
