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
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockHash := []byte(ctx.BlockHeader().LastBlockId.Hash)

	// get pending random number requests for lastBlockHeight
	rqIterator := k.IterateRandomRequestQueueByHeight(ctx, lastBlockHeight)
	defer rqIterator.Close()

	handledNormalRandReqNum := 0
	requestedOracleRandNum := 0

	for ; rqIterator.Valid(); rqIterator.Next() {
		var request types.Request
		k.GetCdc().MustUnmarshalBinaryBare(rqIterator.Value(), &request)

		if request.Oracle {
			// get the request id
			reqID := types.GenerateRequestID(request)

			if err := k.StartRequestContext(ctx, request.ServiceContextID, request.Consumer); err == nil {
				k.SetOracleRandRequest(ctx, request.ServiceContextID, request)
				requestedOracleRandNum++

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeGenerateRandom,
						sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(reqID)),
						sdk.NewAttribute(types.AttributeKeyRequestContextID, request.ServiceContextID.String()),
					),
				)
			} else {
				ctx.Logger().Info(fmt.Sprintf("start service error : %s", err.Error()))
			}

			k.DequeueRandomRequest(ctx, lastBlockHeight, reqID)
		} else {
			// get the request id
			reqID := types.GenerateRequestID(request)

			// generate a random number
			rand := types.MakePRNG(lastBlockHash, currentTimestamp, request.Consumer, nil, false).GetRand()
			k.SetRandom(ctx, reqID, types.NewRandom(request.TxHash, lastBlockHeight, rand.String()))

			// remove the request
			k.DequeueRandomRequest(ctx, lastBlockHeight, reqID)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeGenerateRandom,
					sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(reqID)),
					sdk.NewAttribute(types.AttributeKeyRandom, rand.String()),
				),
			)
			handledNormalRandReqNum++
		}
	}

	k.Logger(ctx).Info(fmt.Sprintf("%d rand requests are handled", handledNormalRandReqNum))
	return
}
