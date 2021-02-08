package types

import (
	"strconv"
)

// NewGenesisState constructs a new GenesisState instance
func NewGenesisState(pendingRequests map[string]Requests) *GenesisState {
	return &GenesisState{
		PendingRandomRequests: pendingRequests,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		PendingRandomRequests: map[string]Requests{},
	}
}

// ValidateGenesis validates the given random genesis state
func ValidateGenesis(data GenesisState) error {
	for height := range data.PendingRandomRequests {
		if _, err := strconv.ParseUint(height, 10, 64); err != nil {
			return err
		}
	}
	return nil
}
