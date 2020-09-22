package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/service/types"
)

// MockTokenKeeper defines a mock implementation for types.TokenKeeper
type MockTokenKeeper struct{}

// GetToken gets the specified token
func (k MockTokenKeeper) GetToken(ctx sdk.Context, denom string) (types.TokenI, error) {
	if denom == sdk.DefaultBondDenom {
		return types.MockToken{
			MinUnit: sdk.DefaultBondDenom,
			Scale:   0,
		}, nil
	}

	return nil, fmt.Errorf("token %s does not exist", denom)
}
