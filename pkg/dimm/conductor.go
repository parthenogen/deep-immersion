package dimm

type conductor interface {
	Beats() <-chan struct{}
}
