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

	id, err := hex.DecodeString(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid htlc id %s", request.Id)
	}

	htlc, found := k.GetHTLC(ctx, id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "HTLC %s not found", request.Id)
	}

	return &types.QueryHTLCResponse{Htlc: &htlc}, nil
}

func (k Keeper) AssetSupply(c context.Context, request *types.QueryAssetSupplyRequest) (*types.QueryAssetSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	assetSupply, found := k.GetAssetSupply(ctx, request.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, string(request.Denom))
	}

	return &types.QueryAssetSupplyResponse{AssetSupply: &assetSupply}, nil
}

func (k Keeper) AssetSupplies(c context.Context, request *types.QueryAssetSuppliesRequest) (*types.QueryAssetSuppliesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	assets := k.GetAllAssetSupplies(ctx)
	if assets == nil {
		assets = []types.AssetSupply{}
	}

	return &types.QueryAssetSuppliesResponse{AssetSupplies: assets}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
