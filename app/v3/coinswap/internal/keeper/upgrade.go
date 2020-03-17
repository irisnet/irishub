package keeper

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

func (k Keeper) Init(ctx sdk.Context, assetKeeper types.AssetKeeper, accountKeeper types.AccountKeeper) {
	logger := k.Logger(ctx).With("handler", "Init")

	tokens := assetKeeper.GetAllTokens(ctx)
	logger.Info("Begin execute upgrade method")

	for _, token := range tokens {
		denom := token.GetDenom()
		voucherCoinName, err := types.GetVoucherCoinName(sdk.IrisAtto, denom)
		if err == nil {
			poolAcc, existed := k.getAccount(ctx, accountKeeper, voucherCoinName)
			if !existed {
				continue
			}
			balance := poolAcc.GetCoins()
			voucherDenom, _ := types.GetVoucherDenom(voucherCoinName)

			liquidity := balance.AmountOf(voucherDenom)
			if liquidity.LTE(sdk.ZeroInt()) {
				continue
			}

			irisToken := sdk.NewCoin(sdk.IrisAtto, balance.AmountOf(sdk.IrisAtto))
			otherToken := sdk.NewCoin(denom, balance.AmountOf(denom))
			voucherToken := sdk.NewCoin(voucherDenom, liquidity)

			coins := sdk.NewCoins(irisToken, otherToken, voucherToken)

			//create pool for voucherCoinName
			logger.Info("Create liquidity pool", "poolName", voucherCoinName)
			_ = k.SetPool(ctx, types.NewPool(voucherCoinName, nil))
			_ = k.SendCoinsFromAccountToPool(ctx, poolAcc.GetAddress(), voucherCoinName, coins)
			logger.Info(
				"Transfer coin to liquidity pool",
				"from", poolAcc.GetAddress().String(),
				"amount", coins.String(),
				"poolName", voucherCoinName,
			)
		}
	}

	logger.Info("End execute upgrade method")
}

//Except for the upgrade process from v2 to v3, please do not use this code
func (k Keeper) getAccount(ctx sdk.Context, accountKeeper types.AccountKeeper, voucherCoinName string) (account auth.Account, existed bool) {
	swapPoolAccAddr := sdk.AccAddress(crypto.AddressHash([]byte(voucherCoinName)))

	account = accountKeeper.GetAccount(ctx, swapPoolAccAddr)
	if account == nil {
		return account, false
	}

	return account, true
}
