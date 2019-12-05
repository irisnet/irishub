package types

// GenesisState contains all rand state that must be provided at genesis
type GenesisState struct {
	PendingRandRequests map[string][]Request // pending rand requests: height->[]Request
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(pendingRequests map[string][]Request) GenesisState {
	return GenesisState{
		PendingRandRequests: pendingRequests,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingRandRequests: map[string][]Request{},
	}
}

// ValidateGenesis validates the given rand genesis state
func ValidateGenesis(data GenesisState) error {
	return nil
}
