package service

import (
	sdk "github.com/irisnet/irishub/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithCoinFlowTrigger(sdk.ServiceEndBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))

	k.SetIntraTxCounter(ctx, 0)

	// handler for the active request on expired
	expiredRequestHandler := func(requestID cmn.HexBytes, request Request) {
		if !request.SuperMode {
			slashTags, err := k.Slash(ctx, requestID)
			if err != nil {
				panic(err)
			}

			if err := k.RefundServiceFee(ctx, request.Consumer, request.ServiceFee); err != nil {
				panic(err)
			}

			tags = tags.AppendTags(slashTags)
		}

		k.DeleteActiveRequest(ctx, request.ServiceName, request.Provider, request.ExpirationHeight, requestID)
	}

	// handler for the expired request batch
	expiredRequestBatchHandler := func(requestContextID cmn.HexBytes, requestContext RequestContext) {
		if requestContext.BatchState != BATCHCOMPLETED {
			k.IterateActiveRequests(ctx, requestContextID, requestContext.BatchCounter, expiredRequestHandler)

			if len(requestContext.ModuleName) != 0 {
				k.Callback(ctx, requestContextID)
			}

			requestContext.BatchState = BATCHCOMPLETED
		}

		k.DeleteRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight())

		if requestContext.State == RUNNING {
			if requestContext.Repeated && (requestContext.RepeatedTotal < 0 || int64(requestContext.BatchCounter) < requestContext.RepeatedTotal) {
				k.AddNewRequestBatch(ctx, requestContextID, ctx.BlockHeight()-requestContext.Timeout+int64(requestContext.RepeatedFrequency))
			} else {
				requestContext.State = COMPLETED
			}
		}

		k.SetRequestContext(ctx, requestContextID, requestContext)
	}

	// handler for the new request batch
	newRequestBatchHandler := func(requestContextID cmn.HexBytes, requestContext RequestContext) {
		if requestContext.State == RUNNING {
			providers, totalPrices := k.FilterServiceProviders(ctx, requestContext.ServiceName, requestContext.Providers, requestContext.ServiceFeeCap)

			if len(requestContext.ModuleName) == 0 || len(providers) >= int(requestContext.ResponseThreshold) {
				if !requestContext.SuperMode {
					if err := k.DeductServiceFees(ctx, requestContext.Consumer, totalPrices); err != nil {
						requestContext.State = PAUSED
						k.SetRequestContext(ctx, requestContextID, requestContext)
					}
				}

				if requestContext.State == RUNNING {
					requestContext.BatchCounter++
					requestContext.BatchResponseCount = 0
					k.SetRequestContext(ctx, requestContextID, requestContext)

					requestTags := k.InitiateRequests(ctx, requestContextID, providers)
					k.AddRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight()+requestContext.Timeout)

					tags = tags.AppendTags(requestTags)
				}
			}
		}

		k.DeleteNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	// handle the expired request batch queue
	k.IterateExpiredRequestBatch(ctx, ctx.BlockHeight(), expiredRequestBatchHandler)

	// handle the new request batch queue
	k.IterateNewRequestBatch(ctx, ctx.BlockHeight(), newRequestBatchHandler)

	return
}
