package mint

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// mint parameters
type Params struct {
	Inflation         sdk.Dec `json:"inflation"`  // inflation rate
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		Inflation:         sdk.NewDecWithPrec(4, 2),
	}
}

func validateParams(params Params) error {
	if params.Inflation.LT(sdk.ZeroDec()) {
		return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	}
	return nil
}
