package types

// NewGenesisState constructs a GenesisState
func NewGenesisState(profilers, trustees []Guardian) *GenesisState {
	return &GenesisState{
		Profilers: profilers,
		Trustees:  trustees,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}
