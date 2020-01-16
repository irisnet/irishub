package keeper

import (
	"github.com/irisnet/irishub/app/v1/auth"
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
			account := k.getAccount(ctx, accountKeeper, uniID)
			balance := account.GetCoins()
			irisToken := sdk.NewCoin(sdk.IrisAtto, balance.AmountOf(sdk.IrisAtto))
			otherToken := sdk.NewCoin(denom, balance.AmountOf(denom))
			coins := sdk.NewCoins(irisToken, otherToken)
			if _, _, err := k.bk.SubtractCoins(ctx, account.GetAddress(), coins); err == nil {
				_ = k.SetPool(ctx, types.NewPool(uniID, coins))
			}
		}
	}
}

//Except for the upgrade process from v2 to v3, please do not use this code
func (k Keeper) getAccount(ctx sdk.Context, accountKeeper types.AccountKeeper, uniId string) auth.Account {
	swapPoolAccAddr := sdk.AccAddress(crypto.AddressHash([]byte(uniId)))
	return accountKeeper.GetAccount(ctx, swapPoolAccAddr)
}
