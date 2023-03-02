package dimm

type dnsClient interface {
	Send(Query) (Response, error)
}
