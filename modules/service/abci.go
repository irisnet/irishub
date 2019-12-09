package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service/internal/types"
)

// EndBlocker handles block ending logic
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := keeper.Logger(ctx)

	// Reset the intra-transaction counter.
	keeper.SetIntraTxCounter(ctx, 0)

	params := keeper.GetParams(ctx)
	slashFraction := params.SlashFraction

	activeIterator := keeper.ActiveRequestQueueIterator(ctx, ctx.BlockHeight())
	defer activeIterator.Close()

	for ; activeIterator.Valid(); activeIterator.Next() {
		var req SvcRequest
		keeper.GetCdc().MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &req)

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

			err := keeper.Slash(ctx, binding, slashCoins)
			if err != nil {
				panic(err)
			}
		}

		keeper.AddReturnFee(ctx, req.Consumer, req.ServiceFee)

		keeper.DeleteActiveRequest(ctx, req)
		keeper.GetMetrics().ActiveRequests.Add(-1)
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
