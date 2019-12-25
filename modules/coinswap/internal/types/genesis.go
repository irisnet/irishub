package types

// GenesisState - coinswap genesis state
type GenesisState struct {
	Params Params `json:"params"`
}

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params Params) GenesisState {
	return GenesisState{Params: params}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams())
}

// ValidateGenesis - placeholder function
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
