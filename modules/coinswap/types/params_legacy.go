package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFee, &p.Fee, validateFee),
		paramtypes.NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		paramtypes.NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
		paramtypes.NewParamSetPair(
			KeyUnilateralLiquidityFee,
			&p.UnilateralLiquidityFee,
			validateUnilateraLiquiditylFee,
		),
	}
}
