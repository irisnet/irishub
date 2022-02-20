package exported

// MT multi token interface
type MT interface {
	GetID() string
	GetSupply() uint64
	GetData() []byte
}
