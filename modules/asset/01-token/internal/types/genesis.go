package types

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params" yaml:"params"` // asset params
	Tokens Tokens `json:"tokens" yaml:"tokens"` // issued tokens
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens Tokens) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

// DefaultGenesisState return the default asset genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), []FungibleToken{})
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	return data.Tokens.Validate()
}
