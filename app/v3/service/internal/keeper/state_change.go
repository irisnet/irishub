package keeper

import (
	"encoding/json"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// CompleteBatch completes a running bath
func (k Keeper) CompleteBatch(ctx sdk.Context, requestContext types.RequestContext, requestContextID cmn.HexBytes,
) (types.RequestContext, sdk.Tags) {
	tags := sdk.NewTags()
	requestContext.BatchState = types.BATCHCOMPLETED

	if len(requestContext.ModuleName) != 0 {
		k.Callback(ctx, requestContextID)
	}

	// remove all requests and responses of this batch
	iterator := k.RequestsIteratorByReqCtx(ctx, requestContextID, requestContext.BatchCounter)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]
		k.DeleteCompactRequest(ctx, requestID)
		k.DeleteResponse(ctx, requestID)
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

func (k Keeper) CompleteServiceContext(ctx sdk.Context, context types.RequestContext, requestContextID cmn.HexBytes) {
	k.DeleteRequestContext(ctx, requestContextID)
}
