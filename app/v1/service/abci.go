package service

import (
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/service/tags"
	"github.com/irisnet/irishub/types"
)

func EndBlocker(ctx types.Context, keeper Keeper) (resTags types.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))
	logger := ctx.Logger()
	// Reset the intra-transaction counter.
	keeper.SetIntraTxCounter(ctx, 0)

	resTags = types.NewTags()
	params := keeper.GetParamSet(ctx)
	slashFraction := params.SlashFraction

	activeIterator := keeper.ActiveRequestQueueIterator(ctx, ctx.BlockHeight())
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		var req SvcRequest
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &req)

		// if not Profiling mode,should slash provider
		slashCoins := types.Coins{}
		if !req.Profiling {
			binding, found := keeper.GetServiceBinding(ctx, req.DefChainID, req.DefName, req.BindChainID, req.Provider)
			if found {
				for _, coin := range binding.Deposit {
					taxAmount := types.NewDecFromInt(coin.Amount).Mul(slashFraction).TruncateInt()
					slashCoins = append(slashCoins, types.NewCoin(coin.Denom, taxAmount))
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

		resTags = resTags.AppendTag(tags.Action, tags.ActionSvcCallTimeOut)
		resTags = resTags.AppendTag(tags.RequestID, []byte(req.RequestID()))
		resTags = resTags.AppendTag(tags.Provider, []byte(req.Provider))
		resTags = resTags.AppendTag(tags.SlashCoins, []byte(slashCoins.String()))
		logger.Info("Remove timeout request", "request_id", req.RequestID(), "consumer", req.Consumer.String())
	}

	return resTags
}
