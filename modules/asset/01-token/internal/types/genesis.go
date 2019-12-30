package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/types"
)

// DefaultToken return the definition of iris token

func IrisToken() FungibleToken {
	return FungibleToken{
		Symbol:        types.Iris,
		Name:          "IRIS Network",
		Scale:         18,
		MinUnit:       types.IrisAtto,
		InitialSupply: sdk.NewIntWithDecimal(20, 9),
		Mintable:      true,
	}
}

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
	return NewGenesisState(DefaultParams(), DefaultTokens())
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	for _, token := range data.Tokens {
		if err := ValidateName(token.Name); err != nil {
			return err
		}
		if err := ValidateScale(token.Scale); err != nil {
			return err
		}

		symbolLen := len(token.Symbol)
		if symbolLen < MinimumTokenSymbolSize || symbolLen > MaximumTokenSymbolSize ||
			!IsBeginWithAlpha(token.Symbol) || !IsAlphaNumeric(token.Symbol) {
			return fmt.Errorf("invalid token symbol %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", token.Symbol, MinimumTokenSymbolSize, MaximumTokenSymbolSize)
		}

		minUnitsLen := len(token.MinUnit)
		if minUnitsLen < MinimumTokenMinUnitSize ||
			minUnitsLen > MaximumTokenMinUnitSize ||
			!IsAlphaNumericDash(token.MinUnit) ||
			!IsBeginWithAlpha(token.MinUnit) {
			return fmt.Errorf("invalid token min_unit %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", token.MinUnit, MinimumTokenMinUnitSize, MaximumTokenMinUnitSize)
		}
	}
	return data.Params.Validate()
}

func DefaultTokens() Tokens {
	return Tokens{
		IrisToken(),
	}
}
