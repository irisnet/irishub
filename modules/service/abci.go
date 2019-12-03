package service

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service/internal/types"
)

// EndBlocker handles block ending logic
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))
	logger := ctx.Logger()
	// Reset the intra-transaction counter.
	keeper.SetIntraTxCounter(ctx, 0)

	params := keeper.GetParamSet(ctx)
	slashFraction := params.SlashFraction

	activeIterator := keeper.ActiveRequestQueueIterator(ctx, ctx.BlockHeight())
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		var req SvcRequest
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &req)

		// if not Profiling mode,should slash provider
		slashCoins := sdk.Coins{}
		if !req.Profiling {
			binding, found := keeper.GetServiceBinding(ctx, req.DefChainID, req.DefName, req.BindChainID, req.Provider)
			if found {
				for _, coin := range binding.Deposit {
					taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(slashFraction).TruncateInt()
					slashCoins = append(slashCoins, sdk.NewCoin(coin.Denom, taxAmount))
				}
			}

			slashCoins = slashCoins.Sort()

			_, err := keeper.ck.BurnCoins(ctx, auth.ServiceDepositCoinsAccAddr, slashCoins)
			if err != nil {
				panic(err)
			}
			err = keeper.Slash(ctx, binding, slashCoins)
			if err != nil {
				panic(err)
			}
		}

		keeper.AddReturnFee(ctx, req.Consumer, req.ServiceFee)

		keeper.DeleteActiveRequest(ctx, req)
		keeper.metrics.ActiveRequests.Add(-1)
		keeper.DeleteRequestExpiration(ctx, req)
	
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSvcCallTimeout,
				sdk.NewAttribute(types.AttributeKeyRequestID, req.RequestID()),
				sdk.NewAttribute(types.AttributeKeyProvider, req.Provider.String()),
				sdk.NewAttribute(types.AttributeKeySlashCoins, slashCoins.String()),
			),
		)

		logger.Info("Remove timeout request", "request_id", req.RequestID(), "consumer", req.Consumer.String())
	}

	return 
}