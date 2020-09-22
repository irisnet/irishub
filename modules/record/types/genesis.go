package types

// NewGenesisState constructs a GenesisState
func NewGenesisState(records []Record) *GenesisState {
	return &GenesisState{
		Records: records,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}
