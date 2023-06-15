package v1

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams constructs a new Params instance
func NewParams(tokenTaxRate sdk.Dec, issueTokenBaseFee sdk.Coin,
	mintTokenFeeRatio sdk.Dec,
) Params {
	return Params{
		TokenTaxRate:      tokenTaxRate,
		IssueTokenBaseFee: issueTokenBaseFee,
		MintTokenFeeRatio: mintTokenFeeRatio,
	}
}

// DefaultParams return the default params
func DefaultParams() Params {
	defaultToken := GetNativeToken()
	return Params{
		TokenTaxRate:      sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(60000)),
		MintTokenFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

// ValidateParams validates the given params
func (p Params) Validate() error {
	if err := validateTaxRate(p.TokenTaxRate); err != nil {
		return err
	}
	if err := validateMintTokenFeeRatio(p.MintTokenFeeRatio); err != nil {
		return err
	}
	if err := validateIssueTokenBaseFee(p.IssueTokenBaseFee); err != nil {
		return err
	}

	return nil
}

func validateTaxRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("token tax rate [%s] should be between [0, 1]", v.String())
	}
	return nil
}

func validateMintTokenFeeRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("fee ratio for minting tokens [%s] should be between [0, 1]", v.String())
	}
	return nil
}

func validateIssueTokenBaseFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("base fee for issuing token should not be negative")
	}
	return nil
}
