package rand

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, k Keeper) {
	currentTimestamp := ctx.BlockHeader().Time.Unix()
	preBlockHeight := ctx.BlockHeight() - 1
	preBlockHash := ctx.BlockHeader().LastBlockId.Hash

	// get pending random number requests for lastBlockHeight
	iterator := k.IterateRandRequestQueueByHeight(ctx, preBlockHeight)
	defer iterator.Close()

	handledRandReqNum := 0
	for ; iterator.Valid(); iterator.Next() {
		var request Request
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		// get the request id
		reqID := GenerateRequestID(request)

		// generate a random number
		rand := MakePRNG(preBlockHash, currentTimestamp, request.Consumer).GetRand()
		k.SetRand(ctx, reqID, NewRand(request.TxHash, preBlockHeight, rand.FloatString(RandPrec)))

		// remove the request
		k.DequeueRandRequest(ctx, preBlockHeight, reqID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeGenerateRand,
				sdk.NewAttribute(AttributeKeyRequestID, hex.EncodeToString(reqID)),
				sdk.NewAttribute(AttributeKeyRand, rand.FloatString(RandPrec)),
			),
		)

		handledRandReqNum++
	}

	k.Logger(ctx).Info(fmt.Sprintf("%d rand requests are handled", handledRandReqNum))
	return
}
