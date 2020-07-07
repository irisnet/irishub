package types

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	Minter Minter `json:"minter"` // minter object
	Params Params `json:"params"` // inflation params
}

// NewGenesisState constructs a GenesisState
func NewGenesisState(minter Minter, params Params) GenesisState {
	return GenesisState{
		Minter: minter,
		Params: params,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minter: DefaultMinter(),
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	return ValidateMinter(data.Minter)
}
