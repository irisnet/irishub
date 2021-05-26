package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Farm params default values
var (
	DefaultCreatePoolFee      = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5000)) // 5000stake
	DefaultMaxRewardCategoryN = uint32(2)
)

// Keys for parameter access
// nolint
var (
	KeyCreatePoolFee      = []byte("CreatePoolFee")
	KeyMaxRewardCategoryN = []byte("MaxRewardCategoryN")
)

// NewParams creates a new Params instance
func NewParams(createPoolFee sdk.Coin, maxRewardCategoryN uint32) Params {
	return Params{
		CreatePoolFee:      createPoolFee,
		MaxRewardCategoryN: maxRewardCategoryN,
	}
}

// ParamSetPairs implements paramstypes.ParamSet
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(
			KeyCreatePoolFee, &p.CreatePoolFee, validateCreatePoolFee),
		paramstypes.NewParamSetPair(
			KeyMaxRewardCategoryN, &p.MaxRewardCategoryN, validateMaxRewardCategoryN),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultCreatePoolFee, DefaultMaxRewardCategoryN)
}

// Validate validates a set of params
func (p Params) Validate() error {
	return validateCreatePoolFee(p.CreatePoolFee)
}

func validateCreatePoolFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", v)
	}
	return nil
}

func validateMaxRewardCategoryN(i interface{}) error { return nil }
