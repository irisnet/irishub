package rand

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewHandler handles all "rand" messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRequestRand:
			return handleMsgRequestRand(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parsed in rand module").Result()
		}

		return sdk.ErrTxDecode("invalid message parsed in rand module").Result()
	}
}

// handleMsgRequestRand handles MsgRequestRand
func handleMsgRequestRand(ctx sdk.Context, k Keeper, msg MsgRequestRand) sdk.Result {
	tags, err := k.RequestRand(ctx, msg.Consumer, msg.BlockInterval)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/rand"))

	// get data of the last block
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockTimestamp := ctx.BlockHeader().Time.Unix()
	lastBlockHash := []byte(ctx.BlockHeader().LastBlockId.Hash)

	// get pending random number requests for lastBlockHeight
	iterator := k.IterateRandRequestQueueByHeight(ctx, lastBlockHeight)
	defer iterator.Close()

	handledRandReqNum := 0
	for ; iterator.Valid(); iterator.Next() {
		var reqID string
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &reqID)

		var request Request
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		// generate a random number
		rand := MakePRNG(lastBlockHash, lastBlockTimestamp, request.Consumer).GetRand()
		k.SetRand(ctx, reqID, NewRand(request.TxHash, lastBlockHeight, rand))

		// remove the request
		k.DequeueRandRequest(ctx, lastBlockHeight, reqID)

		// add tags
		tags.AppendTags(sdk.NewTags(
			TagReqID, []byte(reqID),
			TagRand, []byte(rand.String()),
		))

		handledRandReqNum++
	}

	ctx.Logger().Info(fmt.Sprintf("%d rand requests are handled", handledRandReqNum))
	return
}
