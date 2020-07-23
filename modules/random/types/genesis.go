package types

import (
	"strconv"
)

// GenesisState contains all rand state that must be provided at genesis
type GenesisState struct {
	PendingRandomRequests map[string][]Request `json:"pending_rand_requests" yaml:"pending_rand_requests"` // pending rand requests: height->[]Request
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(pendingRequests map[string][]Request) GenesisState {
	return GenesisState{
		PendingRandomRequests: pendingRequests,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingRandomRequests: map[string][]Request{},
	}
}

// ValidateGenesis validates the given rand genesis state
func ValidateGenesis(data GenesisState) error {
	for height := range data.PendingRandomRequests {
		if _, err := strconv.ParseUint(height, 10, 64); err != nil {
			return err
		}
	}
	return nil
}
