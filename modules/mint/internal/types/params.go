package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/config"
)

var _ params.ParamSet = (*Params)(nil)

// default paramspace for params keeper
const (
	DefaultParamSpace = "mint"
)

//Parameter store key
var (
	// params store for inflation params
	KeyInflation = []byte("Inflation")
	KeyMintDenom = []byte("MintDenom")
)

// ParamTable for mint module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// mint parameters
type Params struct {
	Inflation sdk.Dec `json:"inflation"`  // inflation rate
	MintDenom string  `json:"mint_denom"` // type of coin to mint
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyInflation, &p.Inflation},
		{KeyMintDenom, &p.MintDenom},
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`Mint Params:
  mint/Inflation:  %s,
  mint/MintDenom:  %s`,
		p.Inflation.String(), p.MintDenom)
}

// Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		Inflation: sdk.NewDecWithPrec(4, 2),
		MintDenom: config.StakeDenom,
	}
}

func validateParams(p Params) error {
	if p.Inflation.GT(sdk.NewDecWithPrec(2, 1)) || p.Inflation.LT(sdk.ZeroDec()) {
		return sdk.NewError(params.DefaultCodespace, CodeInvalidMintInflation, fmt.Sprintf("Mint Inflation [%s] should be between [0, 0.2] ", p.Inflation.String()))
	}

	if len(p.MintDenom) == 0 {
		return sdk.NewError(params.DefaultCodespace, CodeInvalidMintDenom, fmt.Sprintf("Mint MintDenom [%s] should not be empty", p.MintDenom))
	}
	return nil
}
