package rand

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/rand"))

	currentTimestamp := ctx.BlockHeader().Time.Unix()
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockHash := []byte(ctx.BlockHeader().LastBlockId.Hash)

	// get pending random number requests for lastBlockHeight
	iterator := k.IterateRandRequestQueueByHeight(ctx, lastBlockHeight)
	defer iterator.Close()

	handledRandReqNum := 0
	for ; iterator.Valid(); iterator.Next() {
		var request Request
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		// get the request id
		reqID := GenerateRequestID(request)

		// generate a random number
		rand := MakePRNG(lastBlockHash, currentTimestamp, request.Consumer).GetRand()
		k.SetRand(ctx, reqID, NewRand(request.TxHash, lastBlockHeight, rand))

		// remove the request
		k.DequeueRandRequest(ctx, lastBlockHeight, reqID)

		// add tags
		tags = tags.AppendTags(sdk.NewTags(
			TagReqID, []byte(hex.EncodeToString(reqID)),
			TagRand, []byte(rand.Rat.FloatString(RandPrec)),
		))

		handledRandReqNum++
	}

	ctx.Logger().Info(fmt.Sprintf("%d rand requests are handled", handledRandReqNum))
	return
}
