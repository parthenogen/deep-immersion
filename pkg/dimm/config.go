package dimm

type config interface {
	Conductor() Conductor
	Sources() []Source
	DNSClients() []DNSClient
	Inspectors() []Inspector
	ErrorHandlers() []ErrorHandler
}
