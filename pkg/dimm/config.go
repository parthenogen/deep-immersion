package dimm

type config interface {
	Conductor() conductor
	Sources() []source
	DNSClients() []dnsClient
	Inspectors() []inspector
	ErrorHandlers() []errorHandler
}
