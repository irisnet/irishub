package types

// NewGenesisState constructs a GenesisState
func NewGenesisState(supers []Super) *GenesisState {
	return &GenesisState{
		Supers: supers,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}
