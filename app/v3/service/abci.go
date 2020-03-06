package service

import (
	sdk "github.com/irisnet/irishub/types"
)

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))

	k.SetIntraTxCounter(ctx, 0)

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
				var requestID []byte
				k.GetCdc().MustUnmarshalBinaryLengthPrefixed(activeReqIterator.Value(), &requestID)

				request, _ := k.GetRequest(ctx, requestID)

				if !request.SuperMode {
					binding, found := k.GetServiceBinding(ctx, request.ServiceName, request.Provider)
					if found {
						slashedCoins := sdk.NewCoins()

						for _, coin := range binding.Deposit {
							taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(slashFraction).TruncateInt()
							slashedCoins = slashedCoins.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxAmount)))
						}

						if err := k.Slash(ctx, binding, slashedCoins); err != nil {
							panic(err)
						}

						if err := k.RefundServiceFee(ctx, request.Consumer, request.ServiceFee); err != nil {
							panic(err)
						}

						tags = tags.AppendTags(sdk.NewTags(
							TagRequestID, []byte(RequestIDToString(requestID)),
							TagProvider, []byte(request.Provider.String()),
							TagConsumer, []byte(request.Consumer.String()),
							TagSlashedCoins, []byte(slashedCoins.String()),
						))
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
		}

		k.DeleteRequestBatchExpiration(ctx, reqContextID, ctx.BlockHeight())

		if reqContext.State == RUNNING {
			if reqContext.Repeated && (reqContext.RepeatedTotal < 0 || int64(reqContext.BatchCounter) < reqContext.RepeatedTotal) {
				k.AddNewRequestBatch(ctx, reqContextID, ctx.BlockHeight()-reqContext.Timeout+int64(reqContext.RepeatedFrequency))
			} else {
				reqContext.State = COMPLETED
			}
		}

		k.SetRequestContext(ctx, reqContextID, reqContext)
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
				if !reqContext.SuperMode {
					if err := k.DeductServiceFees(ctx, reqContext.Consumer, totalPrices); err != nil {
						reqContext.State = PAUSED
						k.SetRequestContext(ctx, reqContextID, reqContext)
					}
				}

				if reqContext.State == RUNNING {
					reqContext.BatchCounter++
					k.SetRequestContext(ctx, reqContextID, reqContext)

					requestTags := k.InitiateRequests(ctx, reqContextID, providers)
					k.AddRequestBatchExpiration(ctx, reqContextID, ctx.BlockHeight()+reqContext.Timeout)

					tags = tags.AppendTags(requestTags)
				}
			}
		}

		k.DeleteNewRequestBatch(ctx, reqContextID, ctx.BlockHeight())
	}

	return
}
