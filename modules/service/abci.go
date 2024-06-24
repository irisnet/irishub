package service

import (
	"encoding/json"
	"strings"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/service/keeper"
	"mods.irisnet.org/service/types"
)

// BeginBlocker handles block beginning logic for service
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "irismod/service"))
	k.SetInternalIndex(ctx, 0)
}

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "irismod/service"))

	// handler for the active request on expired
	expiredRequestHandler := func(requestID tmbytes.HexBytes, request types.Request) {
		_ = k.Slash(ctx, requestID)
		consumer, _ := sdk.AccAddressFromBech32(request.Consumer)
		_ = k.RefundServiceFee(ctx, consumer, request.ServiceFee)

		provider, _ := sdk.AccAddressFromBech32(request.Provider)
		k.DeleteActiveRequest(
			ctx,
			request.ServiceName,
			provider,
			request.ExpirationHeight,
			requestID,
		)
	}

	// handler for the expired request batch
	expiredRequestBatchHandler := func(requestContextID tmbytes.HexBytes, requestContext types.RequestContext) {
		if requestContext.BatchState != types.BATCHCOMPLETED {
			k.IterateActiveRequests(
				ctx,
				requestContextID,
				requestContext.BatchCounter,
				expiredRequestHandler,
			)
			resContext := k.CompleteBatch(ctx, requestContext, requestContextID)
			requestContext = resContext
		}

		k.DeleteRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight())
		k.SetRequestContext(ctx, requestContextID, requestContext)

		if requestContext.State == types.COMPLETED {
			k.CompleteServiceContext(ctx, requestContext, requestContextID)
		}

		if requestContext.State == types.RUNNING {
			if requestContext.Repeated &&
				(requestContext.RepeatedTotal < 0 || int64(requestContext.BatchCounter) < requestContext.RepeatedTotal) {
				k.AddNewRequestBatch(
					ctx,
					requestContextID,
					ctx.BlockHeight()-requestContext.Timeout+int64(
						requestContext.RepeatedFrequency,
					),
				)
			} else {
				k.CompleteServiceContext(ctx, requestContext, requestContextID)
			}
		}

		k.CleanBatch(ctx, requestContext, requestContextID)
	}

	providerRequests := make(map[string][]string)

	// handler for the new request batch
	newRequestBatchHandler := func(requestContextID tmbytes.HexBytes, requestContext *types.RequestContext) {
		consumer, _ := sdk.AccAddressFromBech32(requestContext.Consumer)
		providers := make([]sdk.AccAddress, len(requestContext.Providers))
		for i, provider := range requestContext.Providers {
			pd, _ := sdk.AccAddressFromBech32(provider)
			providers[i] = pd
		}

		if requestContext.State == types.RUNNING {
			providers, totalPrices, rawDenom, err := k.FilterServiceProviders(
				ctx,
				requestContext.ServiceName,
				providers,
				requestContext.Timeout,
				requestContext.ServiceFeeCap,
				consumer,
			)
			if err != nil {
				ctx.EventManager().EmitEvents(sdk.Events{
					sdk.NewEvent(
						types.EventTypeNoExchangeRate,
						sdk.NewAttribute(types.AttributeKeyPriceDenom, rawDenom),
						sdk.NewAttribute(
							types.AttributeKeyRequestContextID,
							requestContextID.String(),
						),
						sdk.NewAttribute(types.AttributeKeyServiceName, requestContext.ServiceName),
						sdk.NewAttribute(types.AttributeKeyConsumer, requestContext.Consumer),
					),
				})
				return
			}

			if len(providers) > 0 && len(providers) >= int(requestContext.ResponseThreshold) {
				if err := k.DeductServiceFees(ctx, consumer, totalPrices); err != nil {
					k.OnRequestContextPaused(
						ctx,
						requestContext,
						requestContextID,
						"insufficient balances",
					)
				}

				if requestContext.State == types.RUNNING {
					_ = k.InitiateRequests(ctx, requestContextID, providers, providerRequests)
					k.AddRequestBatchExpiration(
						ctx,
						requestContextID,
						ctx.BlockHeight()+requestContext.Timeout,
					)
				}
			} else {
				k.SkipCurrentRequestBatch(ctx, requestContextID, *requestContext)
			}

			requestContext, _ := k.GetRequestContext(ctx, requestContextID)
			batchState := types.BatchState{
				BatchCounter:           requestContext.BatchCounter,
				State:                  requestContext.BatchState,
				BatchResponseThreshold: requestContext.BatchResponseThreshold,
				BatchRequestCount:      requestContext.BatchRequestCount,
				BatchResponseCount:     requestContext.BatchResponseCount,
			}
			stateJSON, _ := json.Marshal(batchState)

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeNewBatch,
					sdk.NewAttribute(types.AttributeKeyRequestContextID, requestContextID.String()),
					sdk.NewAttribute(types.AttributeKeyRequestContextState, string(stateJSON)),
				),
			})
		}

		k.DeleteNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	// handle the expired request batch queue
	k.IterateExpiredRequestBatch(ctx, ctx.BlockHeight(), expiredRequestBatchHandler)

	// handle the new request batch queue
	k.IterateNewRequestBatch(ctx, ctx.BlockHeight(), newRequestBatchHandler)

	for provider, requests := range providerRequests {
		requestsJSON, _ := json.Marshal(requests)
		str := strings.Split(provider, ".")
		if len(str) != 2 {
			continue
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeNewBatchRequestProvider,
				sdk.NewAttribute(types.AttributeKeyServiceName, str[0]),
				sdk.NewAttribute(types.AttributeKeyProvider, str[1]),
				sdk.NewAttribute(types.AttributeKeyRequests, string(requestsJSON)),
			),
		})
	}
}
