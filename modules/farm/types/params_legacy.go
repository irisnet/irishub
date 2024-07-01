package types

// Keys for parameter access
// nolint
var (
	KeyPoolCreationFee     = []byte("CreatePoolFee")
	KeyTaxRate             = []byte("TaxRate") // fee key
	KeyMaxRewardCategories = []byte("MaxRewardCategories")
)

// ParamSetPairs implements paramstypes.ParamSet
func (p *Params) ParamSetPairs() ParamSetPairs {
	return ParamSetPairs{
		NewParamSetPair(
			KeyPoolCreationFee,
			&p.PoolCreationFee,
			validatePoolCreationFee,
		),
		NewParamSetPair(
			KeyMaxRewardCategories,
			&p.MaxRewardCategories,
			validateMaxRewardCategories,
		),
	    NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
	}
}

// ParamKeyTable for farm module
func ParamKeyTable() KeyTable {
	return NewKeyTable().RegisterParamSet(&Params{})
}
