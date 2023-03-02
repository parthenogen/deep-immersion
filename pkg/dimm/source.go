package dimm

type source interface {
	GenerateFQDN() string
}
