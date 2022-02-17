package v152

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/coinswap/types"
)

type CoinswapKeeper interface {
	GetParams(ctx sdk.Context) types.Params
	SetParams(ctx sdk.Context, params types.Params)
}

func Migrate(ctx sdk.Context, k CoinswapKeeper) error {
	params := k.GetParams(ctx)
	params.PoolCreationFee = sdk.NewCoin("uiris", sdk.NewIntWithDecimal(5000, 6))
	params.TaxRate = sdk.NewDecWithPrec(4, 1)
	k.SetParams(ctx, params)
	return nil
}
