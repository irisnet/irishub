package types

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // service params
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}
