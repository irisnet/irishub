package service

import (
	"github.com/irisnet/irishub/app/v3/service/internal/types"
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
	logger := ctx.Logger()

	params := k.GetParamSet(ctx)
	slashFraction := params.SlashFraction

	expiredReqBatchIterator := k.ExpiredRequestBatchIterator(ctx, ctx.BlockHeight())
	defer expiredReqBatchIterator.Close()

	for ; expiredReqBatchIterator.Valid(); expiredReqBatchIterator.Next() {
		var reqContextID []byte
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(expiredReqBatchIterator.Value(), &reqContextID)

		reqContext, _ := k.GetRequestContext(ctx, reqContextID)
		if reqContext.BatchState == 0x01 {
			continue
		}

		reqIterator := k.RequestIterator(ctx, reqContextID, reqContext.BatchCounter)
		defer reqIterator.Close()

		respCount := uint16(0)
		var respOutputs []string

		for ; reqIterator.Valid(); reqIterator.Next() {
			requestID := reqIterator.Key[2:]

			var request CompactRequest
			k.GetCdc().MustUnmarshalBinaryLengthPrefixed(reqIterator.Value(), &request)

			resp, found := k.GetResponse(ctx, requestID)
			if !found {
				slashCoins := sdk.NewCoins()
				if !request.Profiling {
					binding, found := k.GetServiceBinding(ctx, request.ServiceName, request.Provider)
					if found {
						for _, coin := range binding.Deposit {
							taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(slashFraction).TruncateInt()
							slashCoins.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxAmount)))
						}

						if err := k.Slash(ctx, binding, slashCoins); err != nil {
							panic(err)
						}
					}
				}
			} else {
				respOutputs = append(respOutputs, resp.Output)

				k.DeleteActiveRequest(ctx, reqContext.ServiceName, request.Provider, ctx.BlockHeight(), requestID)
				k.GetMetrics().ActiveRequests.Add(-1)

				respCount++
			}
		}

		if len(reqContext.ModuleName) != 0 {
			respCallback, _ := k.GetResponseCallback(reqContext.ModuleName)

			if respCount >= reqContext.ResponseThreshold {
				respCallback(ctx, reqContextID, respOutputs)
			} else {
				respCallback(ctx, reqContextID, []string{})
			}
		}

		k.DeleteRequestBatchExpiration(ctx, reqContextID, ctx.BlockHeight())
		k.AddNewRequestBatch(ctx, reqContextID, ctx.BlockHeight()-reqContext.Timeout+int64(reqContext.RepeatedFrequency))
	}

	newReqBatchIterator := k.NewRequestBatchIterator(ctx, ctx.BlockHeight())
	defer newReqBatchIterator.Close()

	for ; newReqBatchIterator.Valid(); newReqBatchIterator.Next() {
		var reqContextID []byte
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(newReqBatchIterator.Value(), &reqContextID)

		reqContext, _ := k.GetRequestContext(ctx, reqContextID)

		providers, totalPrice := k.FilterServiceProviders(ctx, reqContext.ServiceName, reqContext.Providers, reqContext.ServiceFeeCap)
		if len(reqContext.ModuleName) == 0 || len(providers) >= int(reqContext.ResponseThreshold) {
			err := k.DeductServiceFees(ctx, reqContext.Consumer, totalPrice)
			if err != nil {
				reqContext.State = types.RequestContextState(0x01)
			}

			reqContext.BatchCounter++
			k.InitiateRequests(ctx, reqContextID, providers)
		}

		k.DeleteNewRequestBatch(ctx, reqContextID, ctx.BlockHeight())

		if reqContext.State == types.RequestContextState(0x00) {
			if reqContext.Repeated && (reqContext.RepeatedTotal < 0 || int64(reqContext.BatchCounter) < reqContext.RepeatedTotal) {
				k.AddRequestBatchExpiration(ctx, reqContextID, ctx.BlockHeight()+reqContext.Timeout)
			}
		}

		k.SetRequestContext(ctx, reqContextID, reqContext)
	}

	return
}
