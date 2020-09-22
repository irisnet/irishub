package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// parameter keys
var (
	KeyTokenTaxRate      = []byte("TokenTaxRate")
	KeyIssueTokenBaseFee = []byte("IssueTokenBaseFee")
	KeyMintTokenFeeRatio = []byte("MintTokenFeeRatio")
)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTokenTaxRate, &p.TokenTaxRate, validateTaxRate),
		paramtypes.NewParamSetPair(KeyIssueTokenBaseFee, &p.IssueTokenBaseFee, validateIssueTokenBaseFee),
		paramtypes.NewParamSetPair(KeyMintTokenFeeRatio, &p.MintTokenFeeRatio, validateMintTokenFeeRatio),
	}
}

// NewParams token params constructor
func NewParams(tokenTaxRate sdk.Dec, issueTokenBaseFee sdk.Coin,
	mintTokenFeeRatio sdk.Dec,
) Params {
	return Params{
		TokenTaxRate:      tokenTaxRate,
		IssueTokenBaseFee: issueTokenBaseFee,
		MintTokenFeeRatio: mintTokenFeeRatio,
	}
}

// ParamTypeTable returns the TypeTable for coinswap module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// default token module params
func DefaultParams() Params {
	defaultToken := GetNativeToken()
	return Params{
		TokenTaxRate:      sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewIntWithDecimal(60000, int(defaultToken.Scale))),
		MintTokenFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func ValidateParams(p Params) error {
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
