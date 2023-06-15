package v1

import (
	"github.com/irisnet/irismod/types/exported"
)

var _ exported.ParamSet = (*Params)(nil)

// parameter keys
var (
	KeyTokenTaxRate      = []byte("TokenTaxRate")
	KeyIssueTokenBaseFee = []byte("IssueTokenBaseFee")
	KeyMintTokenFeeRatio = []byte("MintTokenFeeRatio")
)

func (p *Params) ParamSetPairs() exported.ParamSetPairs {
	return exported.ParamSetPairs{
		exported.NewParamSetPair(KeyTokenTaxRate, &p.TokenTaxRate, validateTaxRate),
		exported.NewParamSetPair(
			KeyIssueTokenBaseFee,
			&p.IssueTokenBaseFee,
			validateIssueTokenBaseFee,
		),
		exported.NewParamSetPair(
			KeyMintTokenFeeRatio,
			&p.MintTokenFeeRatio,
			validateMintTokenFeeRatio,
		),
	}
}

// ParamKeyTable returns the TypeTable for the token module
func ParamKeyTable() exported.KeyTable {
	return exported.NewKeyTable().RegisterParamSet(&Params{})
}
