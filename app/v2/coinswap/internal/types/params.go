package types

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
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
	Fee sdk.Rat `json:"fee"`
}

// NewParams coinswap params constructor
func NewParams(fee sdk.Rat) Params {
	return Params{
		Fee: fee,
	}
}

// FeeParam defines the numerator and denominator used in calculating the
// amount to be reserved as a liquidity fee.
// TODO: come up with a more descriptive name than Numerator/Denominator
// Fee = 1 - (Numerator / Denominator) TODO: move this to spec
type FeeParam struct {
	Numerator   sdk.Int `json:"fee_numerator"`
	Denominator sdk.Int `json:"fee_denominator"`
}

// NewFeeParam coinswap fee param constructor
func NewFeeParam(numerator, denominator sdk.Int) FeeParam {
	return FeeParam{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

// ParamTypeTable returns the TypeTable for coinswap module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Fee:			%s`, p.Fee.String(),
	)
}

// GetParamSpace Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

// KeyValuePairs  Implements params.KeyValuePairs
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: feeKey, Value: &p.Fee},
	}
}

// Validate Implements params.Validate
func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

// StringFromBytes Implements params.StringFromBytes
func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// ReadOnly Implements params.ReadOnly
func (p *Params) ReadOnly() bool {
	return false
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	feeParam := sdk.NewRat(997, 1000)
	return Params{
		Fee: feeParam,
	}
}

// ValidateParams validates a set of params
func ValidateParams(p Params) error {
	if !p.Fee.GT(sdk.ZeroRat()) {
		return fmt.Errorf("fee is not positive: %s", p.Fee.String())
	}
	return nil
}
