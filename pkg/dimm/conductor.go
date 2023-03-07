package dimm

type Conductor interface {
	Beats() <-chan struct{}
}
