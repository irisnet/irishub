package mint

import (
	"fmt"

	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

// mint parameters
type Params struct {
	MintDenom         string  `json:"mint_denom"` // type of coin to mint
	Inflation         sdk.Dec `json:"inflation"`  // inflation rate
	InflationBasement sdk.Int `json:"inflation_basement"`
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:         stakeTypes.StakeDenom,
		Inflation:         sdk.NewDecWithPrec(4, 2),
		InflationBasement: sdk.NewIntWithDecimal(2, 9).Mul(sdk.NewIntWithDecimal(1, 18)), // 2*(10^9)iris, 2*(10^9)*(10^18)iris-atto
	}
}

func validateParams(params Params) error {
	if params.Inflation.LT(sdk.ZeroDec()) {
		return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	}
	if params.MintDenom != stakeTypes.StakeDenom {
		return fmt.Errorf("invalid mint token denom")
	}
	return nil
}
