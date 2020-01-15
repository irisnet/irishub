package keeper

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

func (k Keeper) Init(ctx sdk.Context, assetKeeper types.AssetKeeper, accountKeeper types.AccountKeeper) {
	tokens := assetKeeper.GetAllTokens(ctx)
	for _, token := range tokens {
		denom := token.GetDenom()
		uniID, err := types.GetUniID(sdk.IrisAtto, denom)
		if err == nil {
			coins := k.getAbandonedPool(ctx, accountKeeper, uniID)
			irisToken := sdk.NewCoin(sdk.IrisAtto, coins.AmountOf(sdk.IrisAtto))
			otherToken := sdk.NewCoin(denom, coins.AmountOf(denom))
			_ = k.SetPool(ctx, types.NewPool(uniID, sdk.NewCoins(irisToken, otherToken)))
		}
	}
}

//Except for the upgrade process from v2 to v3, please do not use this code
func (k Keeper) getAbandonedPool(ctx sdk.Context, accountKeeper types.AccountKeeper, uniId string) (coins sdk.Coins) {
	swapPoolAccAddr := sdk.AccAddress(crypto.AddressHash([]byte(uniId)))
	acc := accountKeeper.GetAccount(ctx, swapPoolAccAddr)
	if acc == nil {
		return coins
	}
	return acc.GetCoins()
}
