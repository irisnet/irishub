package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/record/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Record(c context.Context, req *types.QueryRecordRequest) (*types.QueryRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	record, _ := k.GetRecord(ctx, req.RecordId)
	return &types.QueryRecordResponse{Record: &record}, nil
}
