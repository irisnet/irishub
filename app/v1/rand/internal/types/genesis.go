package types

// GenesisState contains all rand state that must be provided at genesis
type GenesisState struct {
	PendingRandRequests map[int64][]Request // pending rand requests
}
