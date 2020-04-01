package exported

type TokenI interface {
	GetDenom() string
	GetDecimal() uint8
}
