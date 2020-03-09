package service

import (
	"encoding/hex"
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	for reqContextIDStr, requestContext := range data.RequestContexts {
		requestContextID, _ := hex.DecodeString(reqContextIDStr)
		k.SetRequestContext(ctx, requestContextID, requestContext)
	}

	for requestIDStr, request := range data.Requests {
		requestID, _ := ConvertRequestID(requestIDStr)
		k.SetCompactRequest(ctx, requestID, request)
	}

	for requestIDStr, response := range data.Responses {
		requestID, _ := ConvertRequestID(requestIDStr)
		k.SetResponse(ctx, requestID, response)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	requestContexts := make(map[string]RequestContext)
	requests := make(map[string]CompactRequest)
	responses := make(map[string]Response)

	k.IterateRequestContexts(ctx, func(requestContextID cmn.HexBytes, requestContext RequestContext) bool {
		requestContexts[requestContextID.String()] = requestContext
		return false
	})

	k.IterateRequests(ctx, func(requestID cmn.HexBytes, request CompactRequest) bool {
		requests[requestID.String()] = request
		return false
	})

	k.IterateResponses(ctx, func(requestID cmn.HexBytes, response Response) bool {
		responses[requestID.String()] = response
		return false
	})

	return NewGenesisState(
		k.GetParamSet(ctx),
		requestContexts,
		requests,
		responses,
	)
}

// PrepForZeroHeightGenesis refunds the deposits, service fees and earned fees
func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	// refund deposits from all binding services
	if err := k.RefundDeposits(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the deposits: %s", err))
	}

	// refund service fees from all active requests
	if err := k.RefundServiceFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the service fees: %s", err))
	}

	// refund all the earned fees
	if err := k.RefundEarnedFees(ctx); err != nil {
		panic(fmt.Sprintf("failed to refund the earned fees: %s", err))
	}
}
