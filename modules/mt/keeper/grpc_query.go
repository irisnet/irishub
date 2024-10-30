package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mods.irisnet.org/modules/mt/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Supply(context.Context, *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	panic("implement me")
}

func (k Keeper) Denoms(c context.Context, request *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var denoms []types.Denom
	store := ctx.KVStore(k.storeKey)
	denomStore := prefix.NewStore(store, types.KeyDenom(""))
	pageRes, err := query.Paginate(denomStore, request.Pagination, func(key, value []byte) error {
		var denom types.Denom
		k.cdc.MustUnmarshal(value, &denom)
		denoms = append(denoms, denom)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryDenomsResponse{
		Denoms:     denoms,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Denom(c context.Context, request *types.QueryDenomRequest) (*types.QueryDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denomID := strings.TrimSpace(request.DenomId)
	if len(denomID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "denom id is required")
	}

	denom, found := k.GetDenom(ctx, denomID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "denom not found (%s)", request.DenomId)
	}

	return &types.QueryDenomResponse{Denom: &denom}, nil
}

func (k Keeper) MTSupply(c context.Context, request *types.QueryMTSupplyRequest) (*types.QueryMTSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denomID := strings.TrimSpace(request.DenomId)
	mtID := strings.TrimSpace(request.MtId)

	if len(denomID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "denom id is required")
	}

	if len(mtID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "mt id is required")
	}

	return &types.QueryMTSupplyResponse{Amount: k.GetMTSupply(ctx, denomID, mtID)}, nil
}

func (k Keeper) MTs(c context.Context, request *types.QueryMTsRequest) (*types.QueryMTsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denomID := strings.TrimSpace(request.DenomId)
	if len(denomID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "denom id is required")
	}

	var mts []types.MT
	store := ctx.KVStore(k.storeKey)
	mtStore := prefix.NewStore(store, types.KeyMT(denomID, ""))
	pageRes, err := query.Paginate(mtStore, request.Pagination, func(key, value []byte) error {
		var mt types.MT
		k.cdc.MustUnmarshal(value, &mt)

		mt.Supply = k.GetMTSupply(ctx, denomID, mt.GetID())
		mts = append(mts, mt)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryMTsResponse{
		Mts:        mts,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) MT(c context.Context, request *types.QueryMTRequest) (*types.QueryMTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denomID := strings.TrimSpace(request.DenomId)
	mtID := strings.TrimSpace(request.MtId)

	if len(denomID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "denom id is required")
	}

	if len(mtID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "mt id is required")
	}

	mt, err := k.GetMT(ctx, denomID, mtID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "mt not found (%s)", mtID)
	}

	m := types.NewMT(mt.GetID(), mt.GetSupply(), mt.GetData())
	return &types.QueryMTResponse{Mt: &m}, nil
}

func (k Keeper) Balances(c context.Context, request *types.QueryBalancesRequest) (*types.QueryBalancesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	denomID := strings.TrimSpace(request.DenomId)
	owner := strings.TrimSpace(request.Owner)

	if len(denomID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "denom id is required")
	}

	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sender address (%s)", err)
	}

	var bals []types.Balance
	store := ctx.KVStore(k.storeKey)
	mtStore := prefix.NewStore(store, types.KeyBalance(addr, denomID, ""))
	pageRes, err := query.Paginate(mtStore, request.Pagination, func(key, value []byte) error {
		bal := types.Balance{
			MtId:   string(key),
			Amount: types.MustUnMarshalAmount(k.cdc, value),
		}
		bals = append(bals, bal)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryBalancesResponse{
		Balance:    bals,
		Pagination: pageRes,
	}, nil
}
