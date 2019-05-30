package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	// StartingGatewayID is the initial number from which the gateway ids start
	StartingGatewayID = 1
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // asset params
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	// set the initial gateway id
	if err := k.setInitialGatewayID(ctx, StartingGatewayID); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return NewGenesisState(k.GetParamSet(ctx))
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params: DefaultParamsForTest(),
	}
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}
