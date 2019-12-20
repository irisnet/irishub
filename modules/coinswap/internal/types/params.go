package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultParamSpace for coinswap
	DefaultParamspace = ModuleName
	StandardDenom     = sdk.DefaultBondDenom
)

// Parameter store keys
var (
	KeyFee           = []byte("Fee")
	KeyStandardDenom = []byte("StandardDenom")
)

// Params defines the fee and native denomination for coinswap
type Params struct {
	Fee           sdk.Dec `json:"fee" yaml:"fee"`
	StandardDenom string  `json:"standard_denom" yaml:"standard_denom"`
}

// NewParams coinswap params constructor
func NewParams(fee sdk.Dec, feeDenom string) Params {
	return Params{
		Fee:           fee,
		StandardDenom: feeDenom,
	}
}

// ParamTypeTable returns the TypeTable for coinswap module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// KeyValuePairs  Implements params.KeyValuePairs
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{{
		Key:   KeyFee,
		Value: &p.Fee,
	}, {
		Key:   KeyStandardDenom,
		Value: &p.StandardDenom,
	}}
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	fee := sdk.NewDecWithPrec(3, 3)
	return Params{
		Fee:           fee,
		StandardDenom: StandardDenom,
	}
}

// ValidateParams validates a set of params
func (p Params) Validate() error {
	if !p.Fee.GT(sdk.ZeroDec()) || !p.Fee.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be positive and less than 1: %s", p.Fee.String())
	}
	if p.StandardDenom == "" {
		return fmt.Errorf("coinswap parameter standard denom can't be an empty string")
	}
	return nil
}
