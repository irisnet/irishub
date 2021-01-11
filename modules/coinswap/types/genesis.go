package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params Params, denom string) *GenesisState {
	return &GenesisState{
		Params:        params,
		StandardDenom: denom,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), StandardDenom)
}

// ValidateGenesis - placeholder function
func ValidateGenesis(data GenesisState) error {
	if data.StandardDenom == "" {
		return fmt.Errorf("coinswap standard denom can't be an empty string")
	}
	if err := sdk.ValidateDenom(data.StandardDenom); err != nil {
		return err
	}
	return data.Params.Validate()
}
