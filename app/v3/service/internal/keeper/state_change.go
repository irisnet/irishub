package keeper

import (
	"encoding/json"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// CompleteBatch completes a running batch
func (k Keeper) CompleteBatch(ctx sdk.Context, requestContext types.RequestContext, requestContextID cmn.HexBytes,
) (types.RequestContext, sdk.Tags) {
	tags := sdk.NewTags()
	requestContext.BatchState = types.BATCHCOMPLETED

	if len(requestContext.ModuleName) != 0 {
		tags = tags.AppendTags(k.Callback(ctx, requestContextID))
	}

	batchState := types.BatchState{
		BatchCounter:       requestContext.BatchCounter,
		State:              types.BATCHCOMPLETED,
		ResponseThreshold:  requestContext.ResponseThreshold,
		BatchRequestCount:  requestContext.BatchRequestCount,
		BatchResponseCount: requestContext.BatchResponseCount,
	}
	stateJson, _ := json.Marshal(batchState)

	tags = tags.AppendTags(sdk.NewTags(
		sdk.ActionTag(types.ActionCompleteBatch, types.TagRequestContextID), []byte(requestContextID.String()),
		sdk.ActionTag(types.ActionCompleteBatch, requestContextID.String()), stateJson,
	))

	return requestContext, tags
}

// CleanBatch cleans up all requests and responses related to the batch
func (k Keeper) CleanBatch(ctx sdk.Context, requestContext types.RequestContext, requestContextID cmn.HexBytes) {
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
func (k Keeper) CompleteServiceContext(ctx sdk.Context, context types.RequestContext, requestContextID cmn.HexBytes) sdk.Tags {
	tags := sdk.NewTags()
	k.DeleteRequestContext(ctx, requestContextID)

	tags = tags.AppendTags(sdk.NewTags(
		sdk.ActionTag(types.ActionCompleteContext, types.TagRequestContextID), []byte(requestContextID.String()),
	))
	return tags
}

// OnRequestContextPaused handles the event where the specified request context is paused due to certain cause
func (k Keeper) OnRequestContextPaused(
	ctx sdk.Context,
	requestContext types.RequestContext,
	requestContextID cmn.HexBytes,
	cause string,
) sdk.Tags {
	tags := sdk.NewTags()

	requestContext.BatchState = types.BATCHCOMPLETED
	requestContext.State = types.PAUSED

	k.SetRequestContext(ctx, requestContextID, requestContext)

	if len(requestContext.ModuleName) > 0 {
		stateCallback, _ := k.GetStateCallback(requestContext.ModuleName)
		tags = tags.AppendTags(stateCallback(ctx, requestContextID, cause))
	} else {
		tags = tags.AppendTags(sdk.NewTags(
			sdk.ActionTag(types.ActionPauseContext, types.TagRequestContextID), []byte(requestContextID.String()),
		))
	}

	return tags
}
