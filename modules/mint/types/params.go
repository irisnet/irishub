package types

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// default paramspace for params keeper
const (
	DefaultParamSpace = "mint"
)

// Parameter store key
var (
	// params store for inflation params
	KeyInflation = []byte("Inflation")
	KeyMintDenom = []byte("MintDenom")
)

// ParamTable for mint module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, inflation math.LegacyDec) Params {
	return Params{
		MintDenom: mintDenom,
		Inflation: inflation,
	}
}

// DefaultParams returns default minting module parameters
func DefaultParams() Params {
	return Params{
		Inflation: math.LegacyNewDecWithPrec(4, 2),
		MintDenom: sdk.DefaultBondDenom,
	}
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyInflation, &p.Inflation, validateInflation),
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
	}
}

// GetParamSpace implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

// Validate returns err if the Params is invalid
func (p Params) Validate() error {
	if p.Inflation.GT(math.LegacyNewDecWithPrec(2, 1)) || p.Inflation.LT(math.LegacyZeroDec()) {
		return sdkerrors.Wrapf(
			ErrInvalidMintInflation,
			"Mint inflation [%s] should be between [0, 0.2] ",
			p.Inflation.String(),
		)
	}
	if len(p.MintDenom) == 0 {
		return sdkerrors.Wrapf(
			ErrInvalidMintDenom,
			"Mint denom [%s] should not be empty",
			p.MintDenom,
		)
	}
	return nil
}

func validateInflation(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(math.LegacyNewDecWithPrec(2, 1)) || v.LT(math.LegacyZeroDec()) {
		return fmt.Errorf("Mint inflation [%s] should be between [0, 0.2] ", v.String())
	}

	return nil
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
