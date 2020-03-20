package exported

type Token interface {
	GetDenom() string
	GetDecimal() uint8
}
