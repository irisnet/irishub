package service

import (
	"encoding/json"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub/types"
)

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k Keeper) (tags sdk.Tags) {
	tags = sdk.NewTags()
	ctx = ctx.WithCoinFlowTrigger(sdk.ServiceEndBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))

	// handler for the active request on expired
	expiredRequestHandler := func(requestID cmn.HexBytes, request Request) {
		if !request.SuperMode {
			slashTags, _ := k.Slash(ctx, requestID)
			_ = k.RefundServiceFee(ctx, request.Consumer, request.ServiceFee)

			tags = tags.AppendTags(slashTags)
		}

		k.DeleteActiveRequest(ctx, request.ServiceName, request.Provider, request.ExpirationHeight, requestID)
	}

	// handler for the expired request batch
	expiredRequestBatchHandler := func(requestContextID cmn.HexBytes, requestContext RequestContext) {
		if requestContext.BatchState != BATCHCOMPLETED {
			k.IterateActiveRequests(ctx, requestContextID, requestContext.BatchCounter, expiredRequestHandler)
			k.CompleteBatch(ctx, requestContext, requestContextID)
		}

		k.DeleteRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight())

		if requestContext.State == RUNNING {
			if requestContext.Repeated && (requestContext.RepeatedTotal < 0 || int64(requestContext.BatchCounter) < requestContext.RepeatedTotal) {
				k.AddNewRequestBatch(ctx, requestContextID, ctx.BlockHeight()-requestContext.Timeout+int64(requestContext.RepeatedFrequency))
			} else {
				k.CompleteServiceContext(ctx, requestContext, requestContextID)
			}
		}

		k.SetRequestContext(ctx, requestContextID, requestContext)
	}

	providerRequests := make(map[string][]string)

	// handler for the new request batch
	newRequestBatchHandler := func(requestContextID cmn.HexBytes, requestContext RequestContext) {
		if requestContext.State == RUNNING {
			providers, totalPrices := k.FilterServiceProviders(ctx, requestContext.ServiceName, requestContext.Providers, requestContext.ServiceFeeCap)

			if len(providers) > 0 && len(providers) >= int(requestContext.ResponseThreshold) {
				if !requestContext.SuperMode {
					if err := k.DeductServiceFees(ctx, requestContext.Consumer, totalPrices); err != nil {
						requestContext.BatchState = BATCHCOMPLETED
						requestContext.State = PAUSED

						k.SetRequestContext(ctx, requestContextID, requestContext)
					}
				}

				if requestContext.State == RUNNING {
					requestContext.BatchCounter++
					requestContext.BatchResponseCount = 0
					k.SetRequestContext(ctx, requestContextID, requestContext)

					requestTags := k.InitiateRequests(ctx, requestContextID, providers, providerRequests)
					k.AddRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight()+requestContext.Timeout)

					tags = tags.AppendTags(requestTags)
				}
			} else {
				k.SkipCurrentRequestBatch(ctx, requestContextID, requestContext)
			}
		}

		k.DeleteNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	// handle the expired request batch queue
	k.IterateExpiredRequestBatch(ctx, ctx.BlockHeight(), expiredRequestBatchHandler)

	// handle the new request batch queue
	k.IterateNewRequestBatch(ctx, ctx.BlockHeight(), newRequestBatchHandler)

	for provider, requests := range providerRequests {
		requestsJson, _ := json.Marshal(requests)
		tags = tags.AppendTags(sdk.NewTags(
			sdk.ActionTag(types.ActionNewBatchRequest, types.TagProvider), []byte(provider),
			sdk.ActionTag(types.ActionNewBatchRequest, provider), requestsJson,
		))
	}

	return
}
