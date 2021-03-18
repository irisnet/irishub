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
		ctx, currentBlockHeight,
		func(id tmbytes.HexBytes, h types.HTLC) (stop bool) {
			// refund HTLC
			_ = k.RefundHTLC(ctx, h, id)
			// delete from the expiration queue
			k.DeleteHTLCFromExpiredQueue(ctx, currentBlockHeight, id)

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeRefundHTLC,
					sdk.NewAttribute(types.AttributeKeyID, id.String()),
				),
			})

			ctx.Logger().Info(fmt.Sprintf("HTLC [%s] is refunded", id.String()))

			return false
		},
	)

	k.UpdateTimeBasedSupplyLimits(ctx)
}
