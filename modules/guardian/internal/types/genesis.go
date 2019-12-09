package types

// GenesisState - all guardian state that must be provided at genesis
type GenesisState struct {
	Profilers []Guardian `json:"profilers"`
	Trustees  []Guardian `json:"trustees"`
}

func NewGenesisState(profilers, trustees []Guardian) GenesisState {
	return GenesisState{
		Profilers: profilers,
		Trustees:  trustees,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	guardian := Guardian{Description: "genesis", AccountType: Genesis}
	return NewGenesisState([]Guardian{guardian}, []Guardian{guardian})
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return DefaultGenesisState()
}
