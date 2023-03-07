package dimm

type DNSClient interface {
	Send(Query) (Response, error)
}
