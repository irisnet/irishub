package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"
)

// NewParams is the coinswap params constructor
func NewParams(fee, taxRate, unilateralLiquidityFee sdk.Dec, poolCreationFee sdk.Coin) Params {
	return Params{
		Fee:                    fee,
		TaxRate:                taxRate,
		PoolCreationFee:        poolCreationFee,
		UnilateralLiquidityFee: unilateralLiquidityFee,
	}
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	fee := sdk.NewDecWithPrec(3, 3)
	unilateralFee := sdk.NewDecWithPrec(2, 3)
	return Params{
		Fee:                    fee,
		PoolCreationFee:        sdk.NewInt64Coin(sdk.DefaultBondDenom, 5000),
		TaxRate:                sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		UnilateralLiquidityFee: unilateralFee,
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate returns err if Params is invalid
func (p Params) Validate() error {
	if !p.Fee.GT(sdk.ZeroDec()) || !p.Fee.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", p.Fee.String())
	}

	if !p.PoolCreationFee.IsPositive() {
		return fmt.Errorf("poolCreationFee must be positive: %s", p.PoolCreationFee.String())
	}

	if !p.TaxRate.GT(sdk.ZeroDec()) || !p.TaxRate.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", p.TaxRate.String())
	}

	if !p.UnilateralLiquidityFee.GTE(sdk.ZeroDec()) || !p.UnilateralLiquidityFee.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", p.TaxRate.String())
	}
	return nil
}

func validateFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.GT(sdk.ZeroDec()) || !v.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", v.String())
	}

	return nil
}

func validatePoolCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("poolCreationFee must be positive: %s", v.String())
	}
	return nil
}

func validateTaxRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.GT(sdk.ZeroDec()) || !v.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", v.String())
	}
	return nil
}

func validateUnilateraLiquiditylFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// unilateral fee should be in range of [0, 1)
	if !v.GTE(sdk.ZeroDec()) || !v.LT(sdk.OneDec()) {
		return fmt.Errorf(
			"unilateral liquidity fee must be positive and less than 1: %s",
			v.String(),
		)
	}

	return nil
}
