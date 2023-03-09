package ipaddresses

import (
	"net"
	"testing"
)

func TestRandomIPAddr(t *testing.T) {
	const (
		cidr   = "172.16.0.0/12"
		nCases = 1 << 20
	)

	var (
		block  *net.IPNet
		random net.IP
		n      uint
		e      error
	)

	_, block, e = net.ParseCIDR(cidr)
	if e != nil {
		t.Error(e)
	}

	for n = 0; n < nCases; n++ {
		random, e = RandomIPAddr(block)
		if e != nil {
			t.Error(e)
		}

		if !block.Contains(random) {
			t.Log(random)
			t.Fail()
		}
	}
}
