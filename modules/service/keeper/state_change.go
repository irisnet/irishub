package keeper

import (
	"encoding/json"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/service/types"
)

// CompleteBatch completes a running batch
func (k Keeper) CompleteBatch(ctx sdk.Context, requestContext types.RequestContext, requestContextID tmbytes.HexBytes,
) types.RequestContext {
	requestContext.BatchState = types.BATCHCOMPLETED

	if len(requestContext.ModuleName) != 0 {
		k.Callback(ctx, requestContextID)
	}

	batchState := types.BatchState{
		BatchCounter:           requestContext.BatchCounter,
		State:                  types.BATCHCOMPLETED,
		BatchResponseThreshold: requestContext.BatchResponseThreshold,
		BatchRequestCount:      requestContext.BatchRequestCount,
		BatchResponseCount:     requestContext.BatchResponseCount,
	}
	stateJSON, _ := json.Marshal(batchState)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCompleteBatch,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, requestContextID.String()),
			sdk.NewAttribute(types.AttributeKeyRequestContextState, string(stateJSON)),
		),
	})
	return requestContext
}

// CleanBatch cleans up all requests and responses related to the batch
func (k Keeper) CleanBatch(ctx sdk.Context, requestContext types.RequestContext, requestContextID tmbytes.HexBytes) {
	// remove all requests and responses of this batch
	iterator := k.RequestsIteratorByReqCtx(ctx, requestContextID, requestContext.BatchCounter)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]
		k.DeleteCompactRequest(ctx, requestID)
		k.DeleteResponse(ctx, requestID)
	}
}

// CompleteServiceContext completes a running or paused context
func (k Keeper) CompleteServiceContext(ctx sdk.Context, context types.RequestContext, requestContextID tmbytes.HexBytes) {
	k.DeleteRequestContext(ctx, requestContextID)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCompleteContext,
			sdk.NewAttribute(types.AttributeKeyRequestContextID, requestContextID.String()),
		),
	})
}

// OnRequestContextPaused handles the event where the specified request context is paused due to certain cause
func (k Keeper) OnRequestContextPaused(
	ctx sdk.Context,
	requestContext types.RequestContext,
	requestContextID tmbytes.HexBytes,
	cause string,
) {

	requestContext.BatchState = types.BATCHCOMPLETED
	requestContext.State = types.PAUSED

	k.SetRequestContext(ctx, requestContextID, requestContext)

	if len(requestContext.ModuleName) > 0 {
		stateCallback, _ := k.GetStateCallback(requestContext.ModuleName)
		stateCallback(ctx, requestContextID, cause)
	} else {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypePauseContext,
				sdk.NewAttribute(types.AttributeKeyRequestContextID, requestContextID.String()),
			),
		})
	}
}
