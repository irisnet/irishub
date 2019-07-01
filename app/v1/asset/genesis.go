package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params   Params          `json:"params"`   // asset params
	Tokens   []FungibleToken `json:"tokens"`   // issued tokens
	Gateways []Gateway       `json:"gateways"` // created gateways
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	// init gateways
	for _, gateway := range data.Gateways {
		k.SetGateway(ctx, gateway)
		k.SetOwnerGateway(ctx, gateway.Owner, gateway.Moniker)
	}

	// TODO: init tokens with data.Tokens
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// export created gateways
	var gateways []Gateway
	k.IterateGateways(ctx, func(gw Gateway) (stop bool) {
		gateways = append(gateways, gw)
		return false
	})

	var tokens []FungibleToken // TODO: extract existing tokens from app state

	return GenesisState{
		Params:   k.GetParamSet(ctx),
		Tokens:   tokens,
		Gateways: gateways,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:   DefaultParams(),
		Tokens:   []FungibleToken{},
		Gateways: []Gateway{},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		Params:   DefaultParamsForTest(),
		Tokens:   []FungibleToken{},
		Gateways: []Gateway{},
	}
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}

	// validate gateways
	err = validateGateways(data.Gateways)
	if err != nil {
		return err
	}

	return nil
}

// ValidateGateways validates the provided gateways
func validateGateways(gateways []Gateway) error {
	for _, gateway := range gateways {
		if err := gateway.Validate(); err != nil {
			return err
		}
	}

	return nil
}
