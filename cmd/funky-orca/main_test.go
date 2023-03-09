package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"testing"

	"github.com/parthenogen/deep-immersion/pkg/dnsservers"
)

func TestMain(t *testing.T) {
	const (
		argument  = "-test.run=TestMain"
		envFormat = "%s=%s" // https://pkg.go.dev/os/exec#Cmd
		envKey    = "INCEPTION"
		envValue  = "1"
		exitCode  = 0
	)

	var (
		server  interface{ Stop() }
		command *exec.Cmd
		e       error
	)

	if os.Getenv(envKey) == envValue {
		main()
	}

	server, e = dnsservers.NewMockDNSServer(
		net.UDPAddrFromAddrPort(netip.MustParseAddrPort(serverAddrDefault)),
		clientCIDRDefault,
		domainDefault,
		minQPSDefault,
	)
	if e != nil {
		t.Error(e)
	}

	defer server.Stop()

	command = exec.Command(os.Args[0],
		argument,
	)

	command.Stderr = os.Stderr

	command.Env = append(os.Environ(),
		fmt.Sprintf(envFormat, envKey, envValue),
	)

	e = command.Run()
	if e != nil {
		t.Fail()
	}
}
