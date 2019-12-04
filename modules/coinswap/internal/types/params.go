package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultParamSpace for coinswap
	DefaultParamSpace = ModuleName
)

// Parameter store keys
var (
	feeKey = []byte("fee")
)

// Params defines the fee and native denomination for coinswap
type Params struct {
	Fee sdk.Dec `json:"fee"`
}

// NewParams coinswap params constructor
func NewParams(fee sdk.Dec) Params {
	return Params{
		Fee: fee,
	}
}

// ParamTypeTable returns the TypeTable for coinswap module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Coinswap Params:
  Fee:			%s`, p.Fee.String(),
	)
}

// KeyValuePairs  Implements params.KeyValuePairs
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{
			Key:   feeKey,
			Value: &p.Fee,
		},
	}
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	fee := sdk.NewDecWithPrec(3, 3)
	return Params{
		Fee: fee,
	}
}

// ValidateParams validates a set of params
func ValidateParams(p Params) error {
	if !p.Fee.GT(sdk.ZeroDec()) {
		return fmt.Errorf("fee is not positive: %s", p.Fee.String())
	}
	if !p.Fee.LT(sdk.OneDec()) {
		return fmt.Errorf("fee must be less than 1: %s", p.Fee.String())
	}
	return nil
}
