package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k Keeper) {
	// Reset the intra-transaction counter.
	k.SetIntraTxCounter(ctx, 0)

	params := k.GetParams(ctx)
	slashFraction := params.SlashFraction

	activeIterator := k.ActiveRequestQueueIterator(ctx, ctx.BlockHeight())
	defer activeIterator.Close()

	for ; activeIterator.Valid(); activeIterator.Next() {
		var req SvcRequest
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &req)

		// if not Profiling mode,should slash provider
		slashCoins := sdk.NewCoins()
		if !req.Profiling {
			binding, found := k.GetServiceBinding(ctx, req.DefChainID, req.DefName, req.BindChainID, req.Provider)
			if found {
				for _, coin := range binding.Deposit {
					taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(slashFraction).TruncateInt()
					slashCoins.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxAmount))...)
				}
			}

			if err := k.Slash(ctx, binding, slashCoins); err != nil {
				panic(err)
			}
		}

		k.AddReturnFee(ctx, req.Consumer, req.ServiceFee)

		k.DeleteActiveRequest(ctx, req)
		k.DeleteRequestExpiration(ctx, req)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeSvcCallTimeout,
				sdk.NewAttribute(AttributeKeyRequestID, req.RequestID()),
				sdk.NewAttribute(AttributeKeyProvider, req.Provider.String()),
				sdk.NewAttribute(AttributeKeySlashCoins, slashCoins.String()),
			),
		)

		k.Logger(ctx).Info("Remove timeout request", "request_id", req.RequestID(), "consumer", req.Consumer.String())
	}

	return
}
