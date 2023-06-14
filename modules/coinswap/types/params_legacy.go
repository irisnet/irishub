package types

import (
	"github.com/irisnet/irismod/types/exported"
)

// Parameter store keys
var (
	KeyFee                    = []byte("Fee")                    // fee key
	KeyPoolCreationFee        = []byte("PoolCreationFee")        // fee key
	KeyTaxRate                = []byte("TaxRate")                // fee key
	KeyStandardDenom          = []byte("StandardDenom")          // standard token denom key
	KeyUnilateralLiquidityFee = []byte("UnilateralLiquidityFee") // fee key
)

// ParamKeyTable returns the TypeTable for coinswap module
func ParamKeyTable() exported.KeyTable {
	return exported.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() exported.ParamSetPairs {
	return exported.ParamSetPairs{
		exported.NewParamSetPair(KeyFee, &p.Fee, validateFee),
		exported.NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		exported.NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
		exported.NewParamSetPair(
			KeyUnilateralLiquidityFee,
			&p.UnilateralLiquidityFee,
			validateUnilateraLiquiditylFee,
		),
	}
}
