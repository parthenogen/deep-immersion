package ipaddresses

import (
	"crypto/rand"
	"net"
)

func RandomIPAddr(block *net.IPNet) (ip net.IP, e error) {
	var (
		bytes []byte
		i     int
	)

	bytes = make([]byte,
		len(block.IP),
	)

	_, e = rand.Read(bytes)
	if e != nil {
		return
	}

	for i = 0; i < len(block.IP); i++ {
		bytes[i] = []byte(block.IP)[i]&[]byte(block.Mask)[i] |
			bytes[i]&^[]byte(block.Mask)[i] //FIXME: readability
	}

	ip = net.IP(bytes)

	return
}
