package types

// GenesisState contains all rand state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // rand params
}
