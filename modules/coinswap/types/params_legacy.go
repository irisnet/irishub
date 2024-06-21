package types

// Parameter store keys
var (
	KeyFee                    = []byte("Fee")                    // fee key
	KeyPoolCreationFee        = []byte("PoolCreationFee")        // fee key
	KeyTaxRate                = []byte("TaxRate")                // fee key
	KeyStandardDenom          = []byte("StandardDenom")          // standard token denom key
	KeyUnilateralLiquidityFee = []byte("UnilateralLiquidityFee") // fee key
)

// ParamKeyTable returns the TypeTable for coinswap module
func ParamKeyTable() KeyTable {
	return NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() ParamSetPairs {
	return ParamSetPairs{
		NewParamSetPair(KeyFee, &p.Fee, validateFee),
		NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
		NewParamSetPair(
			KeyUnilateralLiquidityFee,
			&p.UnilateralLiquidityFee,
			validateUnilateraLiquiditylFee,
		),
	}
}
