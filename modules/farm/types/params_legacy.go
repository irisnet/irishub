package types

import (
	"github.com/irisnet/irismod/types/exported"
)

// Keys for parameter access
// nolint
var (
	KeyPoolCreationFee     = []byte("CreatePoolFee")
	KeyTaxRate             = []byte("TaxRate") // fee key
	KeyMaxRewardCategories = []byte("MaxRewardCategories")
)

// ParamSetPairs implements paramstypes.ParamSet
func (p *Params) ParamSetPairs() exported.ParamSetPairs {
	return exported.ParamSetPairs{
		exported.NewParamSetPair(
			KeyPoolCreationFee,
			&p.PoolCreationFee,
			validatePoolCreationFee,
		),
		exported.NewParamSetPair(
			KeyMaxRewardCategories,
			&p.MaxRewardCategories,
			validateMaxRewardCategories,
		),
		exported.NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
	}
}

// ParamKeyTable for farm module
func ParamKeyTable() exported.KeyTable {
	return exported.NewKeyTable().RegisterParamSet(&Params{})
}
