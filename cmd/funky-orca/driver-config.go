package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/parthenogen/deep-immersion/pkg/conductors"
	"github.com/parthenogen/deep-immersion/pkg/dimm"
	"github.com/parthenogen/deep-immersion/pkg/dnsclients"
	"github.com/parthenogen/deep-immersion/pkg/errorhandlers"
	"github.com/parthenogen/deep-immersion/pkg/inspectors"
	"github.com/parthenogen/deep-immersion/pkg/ipaddresses"
	"github.com/parthenogen/deep-immersion/pkg/sources"
)

const (
	maxQPSDefault     = 1 << 16
	domainDefault     = "example.org."
	clientCIDRDefault = "127.0.0.0/8"
)

type driverConfig struct {
	conductor     dimm.Conductor
	sources       []dimm.Source
	dnsClients    []dimm.DNSClient
	inspectors    []dimm.Inspector
	errorHandlers []dimm.ErrorHandler
}

func newDriverConfig() (c *driverConfig, e error) {
	const (
		actualBPSLogLabel = "qps"
		network           = "udp"
		failDelay         = 3 * time.Second // DNS client default timeout + 1s

		minQPSFlag    = "min-qps"
		minQPSDefault = 1 << 12
		minQPSUsage   = "Lower limit to number of queries per second"

		maxQPSFlag  = "max-qps"
		maxQPSUsage = "Upper limit to number of queries per second"

		checkIntervalFlag    = "check-interval"
		checkIntervalDefault = 250 * time.Millisecond
		checkIntervalUsage   = "Time interval at which limits are enforced"

		domainFlag  = "domain"
		domainUsage = "Domain for which sub-domains would be generated"

		nSourcesFlag    = "sources"
		nSourcesDefault = 1
		nSourcesUsage   = "Number of sub-domain generators to initialise"

		nDNSClientsFlag    = "dns-clients"
		nDNSClientsDefault = 1
		nDNSClientsUsage   = "Number of concurrent DNS clients to initialise"

		clientCIDRFlag  = "client-cidr"
		clientCIDRUsage = "CIDR block from which queries would be sent"

		serverAddrFlag    = "server-addr"
		serverAddrDefault = "127.46.140.94:5353"
		serverAddrUsage   = "UDP host:port to which queries would be sent"

		nInspectorsFlag    = "inspectors"
		nInspectorsDefault = 1
		nInspectorsUsage   = "Number of response inspectors to initialise"

		nErrorHandlersFlag    = "error-handlers"
		nErrorHandlersDefault = 1
		nErrorHandlersUsage   = "Number of error handlers to initialise"

		expectErrorFlag    = "expect-error"
		expectErrorDefault = false
		expectErrorUsage   = "Exit with code 0 upon encountering client error"

		expectErrDelayFlag    = "expect-error-delay"
		expectErrDelayDefault = 0
		expectErrDelayUsage   = "Duration after which to expect client error"

		acceptRisksFlag    = "accept-risks"
		acceptRisksDefault = false
		acceptRisksUsage   = "Assume all risks related to running this software"
	)

	var (
		minQPS         uint
		maxQPS         uint
		checkInterval  time.Duration
		domain         string
		nSources       uint
		nDNSClients    uint
		clientCIDR     string
		serverAddr     string
		nInspectors    uint
		nErrorHandlers uint
		expectError    bool
		expectErrDelay time.Duration
		acceptRisks    bool

		i uint

		clientIPBlock *net.IPNet
		clientUDPAddr *net.UDPAddr
		serverUDPAddr *net.UDPAddr
	)

	flag.UintVar(&minQPS,
		minQPSFlag,
		minQPSDefault,
		minQPSUsage,
	)

	flag.UintVar(&maxQPS,
		maxQPSFlag,
		maxQPSDefault,
		maxQPSUsage,
	)

	flag.DurationVar(&checkInterval,
		checkIntervalFlag,
		checkIntervalDefault,
		checkIntervalUsage,
	)

	flag.StringVar(&domain,
		domainFlag,
		domainDefault,
		domainUsage,
	)

	flag.UintVar(&nSources,
		nSourcesFlag,
		nSourcesDefault,
		nSourcesUsage,
	)

	flag.UintVar(&nDNSClients,
		nDNSClientsFlag,
		nDNSClientsDefault,
		nDNSClientsUsage,
	)

	flag.StringVar(&clientCIDR,
		clientCIDRFlag,
		clientCIDRDefault,
		clientCIDRUsage,
	)

	flag.StringVar(&serverAddr,
		serverAddrFlag,
		serverAddrDefault,
		serverAddrUsage,
	)

	flag.UintVar(&nInspectors,
		nInspectorsFlag,
		nInspectorsDefault,
		nInspectorsUsage,
	)

	flag.UintVar(&nErrorHandlers,
		nErrorHandlersFlag,
		nErrorHandlersDefault,
		nErrorHandlersUsage,
	)

	flag.BoolVar(&expectError,
		expectErrorFlag,
		expectErrorDefault,
		expectErrorUsage,
	)

	flag.DurationVar(&expectErrDelay,
		expectErrDelayFlag,
		expectErrDelayDefault,
		expectErrDelayUsage,
	)

	flag.BoolVar(&acceptRisks,
		acceptRisksFlag,
		acceptRisksDefault,
		acceptRisksUsage,
	)

	flag.Parse()

	if !acceptRisks {
		e = fmt.Errorf("Accept risks using command-line flag `-%s`.",
			acceptRisksFlag,
		)

		return
	}

	c = &driverConfig{
		conductor: conductors.NewFixedBPSRangeConductor(
			minQPS,
			maxQPS,
			checkInterval,
			failDelay,
			actualBPSLogLabel,
		),
		sources:       make([]dimm.Source, nSources),
		dnsClients:    make([]dimm.DNSClient, nDNSClients),
		inspectors:    make([]dimm.Inspector, nInspectors),
		errorHandlers: make([]dimm.ErrorHandler, nErrorHandlers),
	}

	for i = 0; i < nSources; i++ {
		c.sources[i] = sources.NewUUIDSource(domain)
	}

	serverUDPAddr, e = net.ResolveUDPAddr(network, serverAddr)
	if e != nil {
		return
	}

	_, clientIPBlock, e = net.ParseCIDR(clientCIDR)
	if e != nil {
		return
	}

	for i = 0; i < nDNSClients; i++ {
		clientUDPAddr = new(net.UDPAddr)

		clientUDPAddr.IP, e = ipaddresses.RandomIPAddr(clientIPBlock)
		if e != nil {
			return
		}

		c.dnsClients[i], e = dnsclients.NewTypeADNSClient(
			clientUDPAddr,
			serverUDPAddr,
			checkInterval,
		)
		if e != nil {
			return
		}
	}

	for i = 0; i < nInspectors; i++ {
		c.inspectors[i] = inspectors.NewExitIfTruncatedInspector()
	}

	for i = 0; i < nErrorHandlers; i++ {
		switch expectError {
		case true:
			c.errorHandlers[i] = errorhandlers.NewExpectingErrorHandler(
				expectErrDelay,
			)

		case false:
			c.errorHandlers[i] = errorhandlers.NewExitingErrorHandler()
		}
	}

	return
}

func (c *driverConfig) Conductor() dimm.Conductor {
	return c.conductor
}

func (c *driverConfig) Sources() []dimm.Source {
	return c.sources
}

func (c *driverConfig) DNSClients() []dimm.DNSClient {
	return c.dnsClients
}

func (c *driverConfig) Inspectors() []dimm.Inspector {
	return c.inspectors
}

func (c *driverConfig) ErrorHandlers() []dimm.ErrorHandler {
	return c.errorHandlers
}
