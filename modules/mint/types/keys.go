package types

// nolint
const (
	// ModuleName defines the module name
	ModuleName = "mint"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// Query endpoints supported by the minting querier
	QueryParameters = "parameters"
	QueryInflation  = "inflation"
)

var (
	// use for the keeper store
	MinterKey = []byte{0x00}
	ParamsKey = []byte{0x01}
)
