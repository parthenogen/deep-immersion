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

func TestMainNotExpectingError(t *testing.T) {
	const (
		serverAddr = "127.51.52.232:5353"

		argument0 = "-test.run=TestMainNotExpectingError"
		argument1 = "-server-addr=" + serverAddr
		argument2 = "-accept-risk"
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
		os.Args = []string{os.Args[0],
			argument1,
			argument2,
		}

		main()
	}

	server, e = dnsservers.NewTruncatingMockDNSServer(
		net.UDPAddrFromAddrPort(netip.MustParseAddrPort(serverAddr)),
		clientCIDRDefault,
		domainDefault,
		minQPSDefault,
	)
	if e != nil {
		t.Error(e)
	}

	defer server.Stop()

	command = exec.Command(os.Args[0],
		argument0,
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

func TestMainExpectingError(t *testing.T) {
	const (
		serverAddr = "127.205.99.251:5353"

		argument0 = "-test.run=TestMainExpectingError"
		argument1 = "-server-addr=" + serverAddr
		argument2 = "-expect-error"
		argument3 = "-expect-error-delay=3s" // default client timeout is 2s
		argument4 = "-accept-risk"
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
		os.Args = []string{os.Args[0],
			argument1,
			argument2,
			argument3,
			argument4,
		}

		main()
	}

	server, e = dnsservers.NewDroppingMockDNSServer(
		net.UDPAddrFromAddrPort(netip.MustParseAddrPort(serverAddr)),
		minQPSDefault,
	)
	if e != nil {
		t.Error(e)
	}

	defer server.Stop()

	command = exec.Command(os.Args[0],
		argument0,
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
