package rand

import (
	"encoding/hex"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/irisnet/irishub/types"
)

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/rand"))

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

			if requestContextID, err := k.RequestService(ctx, reqID, request.Consumer); err == nil {
				request.ReqCtxID = requestContextID
				k.EnqueueOracleTimeoutRandRequest(
					ctx,
					lastBlockHeight+k.GetMaxServiceRequestTimeout(ctx),
					reqID,
					request,
				)
				requestedOracleRandNum++
			} else {
				ctx.Logger().Info(fmt.Sprintf("request service error : %s", err.Error()))
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
			tags = tags.AppendTags(sdk.NewTags(
				TagReqID, []byte(hex.EncodeToString(reqID)),
				TagRand, []byte(rand.Rat.FloatString(RandPrec)),
			))

			handledNormalRandReqNum++
		}
	}

	ctx.Logger().Info(fmt.Sprintf("%d normal rand requests are handled", handledNormalRandReqNum))
	ctx.Logger().Info(fmt.Sprintf("%d oracle rand requests are pending", requestedOracleRandNum))

	// ----------------------------------------------------------------------------------------

	// get pending random number requests for lastBlockHeight
	orqIterator := k.IterateRandRequestOracleTimeoutQueueByHeight(ctx, lastBlockHeight)
	defer orqIterator.Close()

	expiredOracleRandReqNum := 0
	for ; orqIterator.Valid(); orqIterator.Next() {
		var request Request
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(orqIterator.Value(), &request)

		// get the request id
		reqID := GenerateRequestID(request)

		// remove the request
		k.DequeueRandRequest(ctx, lastBlockHeight, reqID)

		expiredOracleRandReqNum++
	}

	ctx.Logger().Info(fmt.Sprintf("%d oracle rand requests are expired", expiredOracleRandReqNum))

	return
}
