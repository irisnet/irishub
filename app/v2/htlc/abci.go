package htlc

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// BeginBlocker handles block beginning logic
func BeginBlocker(ctx sdk.Context, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/htlc"))

	currentBlockHeight := uint64(ctx.BlockHeight())
	iterator := k.IterateHTLCExpireQueueByHeight(ctx, currentBlockHeight)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// get the hash lock
		var hashLock []byte
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &hashLock)

		htlc, _ := k.GetHTLC(ctx, hashLock)

		// update the state
		htlc.State = EXPIRED
		k.SetHTLC(ctx, htlc, hashLock)

		// delete from the expiration queue
		k.DeleteHTLCFromExpireQueue(ctx, currentBlockHeight, hashLock)

		// add tags
		tags = tags.AppendTags(sdk.NewTags(
			TagHashLock, []byte(hex.EncodeToString(hashLock)),
		))

		ctx.Logger().Info(fmt.Sprintf("HTLC [%s] is expired", hex.EncodeToString(hashLock)))
	}

	return
}
