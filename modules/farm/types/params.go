package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Farm params default values
var (
	DefaultPoolCreationFee     = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5000)) // 5000stake
	DefaulttTaxRate            = sdk.NewDecWithPrec(4, 1)                            // 0.4 (40%)
	DefaultMaxRewardCategories = uint32(2)
)

// NewParams creates a new Params instance
func NewParams(createPoolFee sdk.Coin, maxRewardCategories uint32, taxRate sdk.Dec) Params {
	return Params{
		PoolCreationFee:     createPoolFee,
		TaxRate:             taxRate,
		MaxRewardCategories: maxRewardCategories,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultPoolCreationFee, DefaultMaxRewardCategories, DefaulttTaxRate)
}

// Validate validates a set of params
func (p Params) Validate() error {
	return validatePoolCreationFee(p.PoolCreationFee)
}

func validatePoolCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", v)
	}
	return nil
}

func validateTaxRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.GT(sdk.ZeroDec()) || !v.LT(sdk.OneDec()) {
		return fmt.Errorf("tax rate must be positive and less than 1: %s", v.String())
	}
	return nil
}

func validateMaxRewardCategories(i interface{}) error { return nil }
