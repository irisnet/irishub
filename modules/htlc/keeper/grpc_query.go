package keeper

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) HTLC(c context.Context, request *types.QueryHTLCRequest) (*types.QueryHTLCResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	hashLock, err := hex.DecodeString(request.HashLock)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid hash lock %s", request.HashLock)
	}

	htlc, found := k.GetHTLC(ctx, hashLock)
	if !found {
		return nil, status.Errorf(codes.NotFound, "HTLC %s not found", request.HashLock)
	}

	return &types.QueryHTLCResponse{Htlc: &htlc}, nil
}
