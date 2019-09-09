package htlc

import (
	"encoding/hex"
	"fmt"

	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// EndBlocker handles block ending logic
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	// check htlc expire and set state from Open to Expired
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "EndBlock").With("module", "iris/htlc"))

	currentBlockHeight := uint64(ctx.BlockHeight())
	iterator := k.IterateHTLCExpireQueueByHeight(ctx, currentBlockHeight)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {

		secretHashLock := iterator.Key()
		htlc, _ := k.GetHTLC(ctx, secretHashLock)

		htlc.State = types.StateExpired
		k.SetHTLC(ctx, htlc, secretHashLock)
		k.DeleteHTLCFromExpireQueue(ctx, currentBlockHeight, secretHashLock)

		ctx.Logger().Info(fmt.Sprintf("HTLC [%s] is expired", hex.EncodeToString(secretHashLock)))
	}

	// TODO: alternative
	// check expire => refund => delete HTLC from expire queue

	return nil
}
