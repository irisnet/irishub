package random

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/random/keeper"
	"github.com/irisnet/irishub/modules/random/types"
)

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	currentTimestamp := ctx.BlockHeader().Time.Unix()
	preBlockHeight := ctx.BlockHeight() - 1
	preBlockHash := ctx.BlockHeader().LastBlockId.Hash

	// get pending random number requests for lastBlockHeight
	iterator := k.IterateRandomRequestQueueByHeight(ctx, preBlockHeight)
	defer iterator.Close()

	handledRandomReqNum := 0
	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		k.GetCdc().MustUnmarshalBinaryBare(iterator.Value(), &request)

		// get the request id
		reqID := types.GenerateRequestID(request)

		// generate a random number
		rand := types.MakePRNG(preBlockHash, currentTimestamp, request.Consumer).GetRandom()
		k.SetRandom(ctx, reqID, types.NewRandom(request.TxHash, preBlockHeight, rand.FloatString(types.RandomPrec)))

		// remove the request
		k.DequeueRandomRequest(ctx, preBlockHeight, reqID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeGenerateRandom,
				sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(reqID)),
				sdk.NewAttribute(types.AttributeKeyRandom, rand.FloatString(types.RandomPrec)),
			),
		)

		handledRandomReqNum++
	}

	k.Logger(ctx).Info(fmt.Sprintf("%d rand requests are handled", handledRandomReqNum))
	return
}
