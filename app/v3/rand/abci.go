package rand

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/irisnet/irishub/types"
)

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/rand"))
	tags = sdk.NewTags()

	currentTimestamp := ctx.BlockHeader().Time.Unix()
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockHash := []byte(ctx.BlockHeader().LastBlockId.Hash)

	// get pending random number requests for lastBlockHeight
	rqIterator := k.IterateRandRequestQueueByHeight(ctx, lastBlockHeight)
	defer rqIterator.Close()

	handledNormalRandReqNum := 0
	requestedOracleRandNum := 0
	for ; rqIterator.Valid(); rqIterator.Next() {
		var request Request
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(rqIterator.Value(), &request)

		if request.Oracle {
			// get the request id
			reqID := GenerateRequestID(request)

			if err := k.StartRequestContext(ctx, request.ServiceContextID, request.Consumer); err == nil {
				k.SetOracleRandRequest(ctx, request.ServiceContextID, request)
				requestedOracleRandNum++

				// add tags
				tags = tags.AppendTags(
					sdk.NewTags(
						TagReqID, []byte(reqID.String()),
						TagRequestContextID, []byte(request.ServiceContextID.String()),
					),
				)
			} else {
				ctx.Logger().Info(fmt.Sprintf("start service error : %s", err.Error()))
			}

			k.DequeueRandRequest(ctx, lastBlockHeight, reqID)
		} else {
			// get the request id
			reqID := GenerateRequestID(request)

			// generate a random number
			rand := MakePRNG(lastBlockHash, currentTimestamp, request.Consumer, nil, false).GetRand()
			k.SetRand(ctx, reqID, NewRand(request.TxHash, lastBlockHeight, rand))

			// remove the request
			k.DequeueRandRequest(ctx, lastBlockHeight, reqID)

			// add tags
			tags = tags.AppendTags(
				sdk.NewTags(
					TagReqID, []byte(reqID.String()),
					TagRand(reqID.String()), []byte(rand.Rat.FloatString(RandPrec)),
				),
			)

			handledNormalRandReqNum++
		}
	}

	ctx.Logger().Info(fmt.Sprintf("%d normal rand requests are handled", handledNormalRandReqNum))
	ctx.Logger().Info(fmt.Sprintf("%d oracle rand requests are pending", requestedOracleRandNum))

	return
}
