package random

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/random/keeper"
	"github.com/irisnet/irismod/modules/random/types"
)

// BeginBlocker handles block beginning logic for random
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	currentTimestamp := ctx.BlockHeader().Time.Unix()
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockHash := ctx.BlockHeader().LastBlockId.Hash

	// get pending random number requests for lastBlockHeight
	rqIterator := k.IterateRandomRequestQueueByHeight(ctx, lastBlockHeight)
	defer rqIterator.Close()

	handledNormalRandReqNum := 0
	requestedOracleRandNum := 0

	for ; rqIterator.Valid(); rqIterator.Next() {
		var request types.Request
		k.GetCdc().MustUnmarshalBinaryBare(rqIterator.Value(), &request)

		consumer, _ := sdk.AccAddressFromBech32(request.Consumer)
		serviceContextID, _ := hex.DecodeString(request.ServiceContextID)

		if request.Oracle {
			// get the request id
			reqID := types.GenerateRequestID(request)

			if err := k.StartRequestContext(ctx, serviceContextID, consumer); err == nil {
				k.SetOracleRandRequest(ctx, serviceContextID, request)
				requestedOracleRandNum++

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeRequestService,
						sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(reqID)),
						sdk.NewAttribute(types.AttributeKeyRequestContextID, request.ServiceContextID),
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
			random := types.MakePRNG(lastBlockHash, currentTimestamp, consumer, nil, false).GetRand()
			k.SetRandom(ctx, reqID, types.NewRandom(request.TxHash, lastBlockHeight, random.FloatString(types.RandPrec)))

			// remove the request
			k.DequeueRandomRequest(ctx, lastBlockHeight, reqID)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeGenerateRandom,
					sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(reqID)),
					sdk.NewAttribute(types.AttributeKeyRandom, random.String()),
				),
			)
			handledNormalRandReqNum++
		}
	}

	k.Logger(ctx).Info(fmt.Sprintf("%d random requests are handled", handledNormalRandReqNum))
}
