package dnsservers

import (
	"net"
)

type mockDNSServer struct {
	udpConn  *net.UDPConn
	incoming chan message
	outgoing chan message
	stop     chan struct{}
}

func newMockDNSServer(address *net.UDPAddr, handler func(message) message) (
	s *mockDNSServer, e error,
) {
	const (
		network = "udp"
	)

	s = &mockDNSServer{
		incoming: make(chan message),
		outgoing: make(chan message),
		stop:     make(chan struct{}),
	}

	s.udpConn, e = net.ListenUDP(network, address)
	if e != nil {
		return
	}

	go s.runReader()
	go s.runWriter()
	go s.runHandler(handler)

	return
}

func (s *mockDNSServer) Stop() {
	close(s.stop)

	return
}

func (s *mockDNSServer) runReader() {
	const (
		udpMsgSizeUppLim = 1 << 16
	)

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

func (s *mockDNSServer) runHandler(handler func(message) message) {
	var (
		incoming message
	)

	for {
		select {
		case <-s.stop:
			return

		default:
			incoming = <-s.incoming

			s.outgoing <- handler(incoming)
		}
	}
}
