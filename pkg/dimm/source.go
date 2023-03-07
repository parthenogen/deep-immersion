package dimm

type Source interface {
	GenerateFQDN() string
}
