package rand

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/rand/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/rand"))

	currentTimestamp := ctx.BlockHeader().Time.Unix()
	preBlockHeight := ctx.BlockHeight() - 1
	preBlockHash := []byte(ctx.BlockHeader().LastBlockId.Hash)

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
		k.SetRand(ctx, reqID, NewRand(request.TxHash, preBlockHeight, rand))

		// remove the request
		k.DequeueRandRequest(ctx, preBlockHeight, reqID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeGenerateRand,
				sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(reqID)),
				sdk.NewAttribute(types.AttributeKeyRand, rand.FloatString(RandPrec)),
			),
		)

		handledRandReqNum++
	}

	ctx.Logger().Info(fmt.Sprintf("%d rand requests are handled", handledRandReqNum))
	return
}
