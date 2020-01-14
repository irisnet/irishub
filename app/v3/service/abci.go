package service

import (
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))
	logger := ctx.Logger()

	// Reset the intra-transaction counter.
	k.SetIntraTxCounter(ctx, 0)

	params := k.GetParamSet(ctx)
	slashFraction := params.SlashFraction

	activeIterator := k.ActiveRequestQueueIterator(ctx, ctx.BlockHeight())
	defer activeIterator.Close()

	for ; activeIterator.Valid(); activeIterator.Next() {
		var req SvcRequest
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &req)

		// if not Profiling mode,should slash provider
		slashCoins := sdk.NewCoins()
		if !req.Profiling {
			binding, found := k.GetServiceBinding(ctx, req.DefName, req.Provider)
			if found {
				for _, coin := range binding.Deposit {
					taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(slashFraction).TruncateInt()
					slashCoins.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxAmount)))
				}
			}

			if err := k.Slash(ctx, binding, slashCoins); err != nil {
				panic(err)
			}
		}

		k.AddReturnFee(ctx, req.Consumer, req.ServiceFee)

		k.DeleteActiveRequest(ctx, req)
		k.GetMetrics().ActiveRequests.Add(-1)
		k.DeleteRequestExpiration(ctx, req)

		resTags = resTags.AppendTag(types.TagAction, types.TagActionSvcCallTimeOut)
		resTags = resTags.AppendTag(types.TagRequestID, []byte(req.RequestID()))
		resTags = resTags.AppendTag(types.TagProvider, []byte(req.Provider))
		resTags = resTags.AppendTag(types.TagSlashCoins, []byte(slashCoins.String()))

		logger.Info("Remove timeout request", "request_id", req.RequestID(), "consumer", req.Consumer.String())
	}

	return
}
