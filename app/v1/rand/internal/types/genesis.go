package types

// GenesisState contains all rand state that must be provided at genesis
type GenesisState struct {
	GeneratedRands     []Rand    // generated rands
	PendingRandRequest []Request // pending rand requests
}
