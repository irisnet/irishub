package asset

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params   Params    `json:"params"`   // asset params
	Tokens   Tokens    `json:"tokens"`   // issued tokens
	Gateways []Gateway `json:"gateways"` // created gateways
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

	//init tokens
	for _, token := range data.Tokens {
		_, _, err := k.AddToken(ctx, token)
		if err != nil {
			panic(err.Error())
		}
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// export created gateways
	var gateways []Gateway
	k.IterateGateways(ctx, func(gw Gateway) (stop bool) {
		gateways = append(gateways, gw)
		return false
	})

	// export created token
	var tokens Tokens
	k.IterateTokens(ctx, func(token FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})
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
	if err := validateGateways(data.Gateways); err != nil {
		return err
	}
	// validate tokens
	if err := data.Tokens.Validate(); err != nil {
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
