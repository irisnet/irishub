package types

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params" yaml:"params"` // service params
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided service genesis state
func ValidateGenesis(data GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
