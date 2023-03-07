package dimm

type ErrorHandler interface {
	Handle(error)
}
