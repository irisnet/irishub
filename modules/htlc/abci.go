package htlc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker handles block beginning logic for HTLC
func BeginBlocker(ctx sdk.Context, k Keeper) {
	currentBlockHeight := uint64(ctx.BlockHeight())
	iterator := k.IterateHTLCExpireQueueByHeight(ctx, currentBlockHeight)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// get the hash lock
		var hashLock HTLCHashLock
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &hashLock)

		htlc, _ := k.GetHTLC(ctx, hashLock)

		// update the state
		htlc.State = EXPIRED
		k.SetHTLC(ctx, htlc, hashLock)

		// delete from the expiration queue
		k.DeleteHTLCFromExpireQueue(ctx, currentBlockHeight, hashLock)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeExpiredHTLC,
				sdk.NewAttribute(AttributeValueHashLock, hashLock.String()),
			),
		)

		k.Logger(ctx).Info(fmt.Sprintf("HTLC [%s] is expired", hashLock.String()))
	}

	return
}
