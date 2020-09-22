package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) HTLC(c context.Context, request *types.QueryHTLCRequest) (*types.QueryHTLCResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	htlc, found := k.GetHTLC(ctx, request.HashLock)
	if !found {
		return nil, status.Errorf(codes.NotFound, "HTLC %s not found", request.HashLock.String())
	}

	return &types.QueryHTLCResponse{Htlc: &htlc}, nil
}
