package dimm

type errorHandler interface {
	Handle(error)
}
