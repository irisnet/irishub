package types

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // service params
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(params Params) GenesisState {
	return GenesisState{Params: params}
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{Params: DefaultParams()}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params: DefaultParamsForTest(),
	}
}

// ValidateGenesis validates the provided service genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}
