package keeper

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/record/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Record(c context.Context, req *types.QueryRecordRequest) (*types.QueryRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	recordId, err := hex.DecodeString(req.RecordId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid record ID %s", req.RecordId)
	}
	record, _ := k.GetRecord(ctx, recordId)
	return &types.QueryRecordResponse{Record: &record}, nil
}
