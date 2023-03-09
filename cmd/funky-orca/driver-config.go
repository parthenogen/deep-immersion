package main

import (
	"flag"
	"net"
	"time"

	"github.com/parthenogen/deep-immersion/pkg/conductors"
	"github.com/parthenogen/deep-immersion/pkg/dimm"
	"github.com/parthenogen/deep-immersion/pkg/dnsclients"
	"github.com/parthenogen/deep-immersion/pkg/errorhandlers"
	"github.com/parthenogen/deep-immersion/pkg/inspectors"
	"github.com/parthenogen/deep-immersion/pkg/sources"
)

const (
	minQPSDefault     = 1 << 14
	maxQPSDefault     = 1 << 16
	domainDefault     = "example.org."
	clientCIDRDefault = "127.0.0.0/8"
	serverAddrDefault = "127.46.140.94:5353"
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

		minQPSFlag  = "min-qps"
		minQPSUsage = "Lower limit to number of queries per second"

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

		clientAddrFlag    = "client-addr"
		clientAddrDefault = "127.37.98.54:35353"
		clientAddrUsage   = "UDP host:port from which queries would be sent"

		serverAddrFlag  = "server-addr"
		serverAddrUsage = "UDP host:port to which queries would be sent"

		nInspectorsFlag    = "inspectors"
		nInspectorsDefault = 1
		nInspectorsUsage   = "Number of response inspectors to initialise"

		nErrorHandlersFlag    = "error-handlers"
		nErrorHandlersDefault = 1
		nErrorHandlersUsage   = "Number of error handlers to initialise"
	)

	var (
		minQPS         uint
		maxQPS         uint
		checkInterval  time.Duration
		domain         string
		nSources       uint
		nDNSClients    uint
		clientAddr     string
		serverAddr     string
		nInspectors    uint
		nErrorHandlers uint

		i uint

		clientUDPAddr *net.UDPAddr
		serverUDPAddr *net.UDPAddr
	)

	flag.UintVar(&maxQPS,
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

	flag.StringVar(&clientAddr,
		clientAddrFlag,
		clientAddrDefault,
		clientAddrUsage,
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

	flag.Parse()

	c = &driverConfig{
		conductor: conductors.NewFixedBPSRangeConductor(
			minQPS,
			maxQPS,
			checkInterval,
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

	clientUDPAddr, e = net.ResolveUDPAddr(network, clientAddr)
	if e != nil {
		return
	}

	serverUDPAddr, e = net.ResolveUDPAddr(network, serverAddr)
	if e != nil {
		return
	}

	for i = 0; i < nDNSClients; i++ {
		c.dnsClients[i], e = dnsclients.NewTypeADNSClient(
			clientUDPAddr,
			serverUDPAddr,
		)
		if e != nil {
			return
		}
	}

	for i = 0; i < nInspectors; i++ {
		c.inspectors[i] = inspectors.NewExitIfTruncatedInspector()
	}

	for i = 0; i < nErrorHandlers; i++ {
		c.errorHandlers[i] = errorhandlers.NewExitingErrorHandler()
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
