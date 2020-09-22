package htlc

import (
	"fmt"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
)

// BeginBlocker handles block beginning logic for HTLC
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "irismod/htlc"))

	currentBlockHeight := uint64(ctx.BlockHeight())

	k.IterateHTLCExpiredQueueByHeight(
		ctx,
		currentBlockHeight,
		func(hlock tmbytes.HexBytes, h types.HTLC) (stop bool) {
			// update the state
			h.State = types.Expired
			k.SetHTLC(ctx, h, hlock)

			// delete from the expiration queue
			k.DeleteHTLCFromExpiredQueue(ctx, currentBlockHeight, hlock)

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeHTLCExpired,
					sdk.NewAttribute(types.AttributeKeyHashLock, hlock.String()),
				),
			})

			ctx.Logger().Info(fmt.Sprintf("HTLC [%s] is expired", hlock.String()))

			return false
		},
	)
}
