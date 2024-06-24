package v1

import (
	"mods.irisnet.org/modules/token/types"
)

var _ types.ParamSet = (*Params)(nil)

// parameter keys
var (
	KeyTokenTaxRate      = []byte("TokenTaxRate")
	KeyIssueTokenBaseFee = []byte("IssueTokenBaseFee")
	KeyMintTokenFeeRatio = []byte("MintTokenFeeRatio")
)

func (p *Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{
		types.NewParamSetPair(KeyTokenTaxRate, &p.TokenTaxRate, validateTaxRate),
		types.NewParamSetPair(
			KeyIssueTokenBaseFee,
			&p.IssueTokenBaseFee,
			validateIssueTokenBaseFee,
		),
		types.NewParamSetPair(
			KeyMintTokenFeeRatio,
			&p.MintTokenFeeRatio,
			validateMintTokenFeeRatio,
		),
	}
}

// ParamKeyTable returns the TypeTable for the token module
func ParamKeyTable() types.KeyTable {
	return types.NewKeyTable().RegisterParamSet(&Params{})
}
