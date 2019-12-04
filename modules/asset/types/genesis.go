package types

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // asset params
	Tokens Tokens `json:"tokens"` // issued tokens
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens Tokens) GenesisState {
	return GenesisState{Params: params, Tokens: tokens}
}
