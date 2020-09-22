package types

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// StandardDenom for coinswap
	StandardDenom = sdk.DefaultBondDenom
)

// Parameter store keys
var (
	KeyFee           = []byte("Fee")           // fee key
	KeyStandardDenom = []byte("StandardDenom") // standard token denom key
)

// NewParams coinswap paramtypes constructor
func NewParams(fee sdk.Dec, feeDenom string) Params {
	return Params{
		Fee:           fee,
		StandardDenom: feeDenom,
	}
}

// ParamTypeTable returns the TypeTable for coinswap module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// KeyValuePairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFee, &p.Fee, validateFee),
		paramtypes.NewParamSetPair(KeyStandardDenom, &p.StandardDenom, validateStandardDenom),
	}
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	fee := sdk.NewDecWithPrec(3, 3)
	return Params{
		Fee:           fee,
		StandardDenom: StandardDenom,
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
	if p.StandardDenom == "" {
		return fmt.Errorf("coinswap parameter standard denom can't be an empty string")
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

func validateStandardDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("standard denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
