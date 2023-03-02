package dimm

type Query interface {
	FQDN() string
}

type query struct {
	fqdn string
}

func (q *query) FQDN() string {
	return q.fqdn
}
