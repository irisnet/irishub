package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keys for parameter access
// nolint
var (
	KeyPoolCreationFee     = []byte("CreatePoolFee")
	KeyTaxRate             = []byte("TaxRate") // fee key
	KeyMaxRewardCategories = []byte("MaxRewardCategories")
)

// ParamSetPairs implements paramstypes.ParamSet
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(
			KeyPoolCreationFee,
			&p.PoolCreationFee,
			validatePoolCreationFee,
		),
		paramstypes.NewParamSetPair(
			KeyMaxRewardCategories,
			&p.MaxRewardCategories,
			validateMaxRewardCategories,
		),
		paramstypes.NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
	}
}

// ParamKeyTable for farm module
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}
