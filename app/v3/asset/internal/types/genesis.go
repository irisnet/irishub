package types

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // asset params
	Tokens Tokens `json:"tokens"` // issued tokens
}
