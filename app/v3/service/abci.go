package service

import (
	sdk "github.com/irisnet/irishub/types"
)

// BeginBlocker handles block beginning logic for service
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// reset the tx counter
	k.SetIntraTxCounter(ctx, 0)
}

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))
	// logger := ctx.Logger()

	params := k.GetParamSet(ctx)
	slashFraction := params.SlashFraction

	// handle expired request batch queue

	expiredReqBatchIterator := k.ExpiredRequestBatchIterator(ctx, ctx.BlockHeight())
	defer expiredReqBatchIterator.Close()

	for ; expiredReqBatchIterator.Valid(); expiredReqBatchIterator.Next() {
		var reqContextID []byte
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(expiredReqBatchIterator.Value(), &reqContextID)

		reqContext, _ := k.GetRequestContext(ctx, reqContextID)

		if reqContext.BatchState != BATCHCOMPLETED {
			// iterate active requests
			activeReqIterator := k.ActiveRequestsIteratorByReqCtx(ctx, reqContextID, reqContext.BatchCounter)
			defer activeReqIterator.Close()

			for ; activeReqIterator.Valid(); activeReqIterator.Next() {
				requestID := activeReqIterator.Value()
				k.GetCdc().MustUnmarshalBinaryLengthPrefixed(activeReqIterator.Value(), &requestID)

				request, _ := k.GetRequest(ctx, requestID)

				if !request.Profiling {
					binding, found := k.GetServiceBinding(ctx, request.ServiceName, request.Provider)
					if found {
						slashedCoins := sdk.NewCoins()

						for _, coin := range binding.Deposit {
							taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(slashFraction).TruncateInt()
							slashedCoins.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxAmount)))
						}

						if err := k.Slash(ctx, binding, slashedCoins); err != nil {
							panic(err)
						}

						if err := k.RefundServiceFee(ctx, request.Consumer, request.ServiceFee); err != nil {
							panic(err)
						}
					}
				}

				k.DeleteActiveRequest(ctx, reqContext.ServiceName, request.Provider, ctx.BlockHeight(), requestID)
				k.DeleteActiveRequestByID(ctx, requestID)

				k.GetMetrics().ActiveRequests.Add(-1)
			}

			// callback
			if len(reqContext.ModuleName) != 0 {
				if reqContext.BatchResponseCount >= reqContext.ResponseThreshold {
					respCallback, _ := k.GetResponseCallback(reqContext.ModuleName)
					respCallback(ctx, reqContextID, k.GetResponsesOutput(ctx, reqContextID, reqContext.BatchCounter))
				}
			}

			reqContext.BatchState = BATCHCOMPLETED
			k.SetRequestContext(ctx, reqContextID, reqContext)
		}

		k.DeleteRequestBatchExpiration(ctx, reqContextID, ctx.BlockHeight())

		if reqContext.State == RUNNING && reqContext.Repeated && (reqContext.RepeatedTotal < 0 || int64(reqContext.BatchCounter) < reqContext.RepeatedTotal) {
			k.AddNewRequestBatch(ctx, reqContextID, ctx.BlockHeight()-reqContext.Timeout+int64(reqContext.RepeatedFrequency))
		}
	}

	// handle new request batch queue
	newReqBatchIterator := k.NewRequestBatchIterator(ctx, ctx.BlockHeight())
	defer newReqBatchIterator.Close()

	for ; newReqBatchIterator.Valid(); newReqBatchIterator.Next() {
		var reqContextID []byte
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(newReqBatchIterator.Value(), &reqContextID)

		reqContext, _ := k.GetRequestContext(ctx, reqContextID)

		if reqContext.State == RUNNING {
			providers, totalPrices := k.FilterServiceProviders(ctx, reqContext.ServiceName, reqContext.Providers, reqContext.ServiceFeeCap)
			if len(reqContext.ModuleName) == 0 || len(providers) >= int(reqContext.ResponseThreshold) {
				if err := k.DeductServiceFees(ctx, reqContext.Consumer, totalPrices); err != nil {
					reqContext.State = PAUSED
				}

				reqContext.BatchCounter++
				k.SetRequestContext(ctx, reqContextID, reqContext)

				k.InitiateRequests(ctx, reqContextID, providers)
				k.AddRequestBatchExpiration(ctx, reqContextID, ctx.BlockHeight()+reqContext.Timeout)
			}
		}

		k.DeleteNewRequestBatch(ctx, reqContextID, ctx.BlockHeight())
	}

	return
}
