package asset

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryAsset   = "asset"
	QueryGateway = "gateway"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAsset:
			return queryAsset(ctx, req, k)
		case QueryGateway:
			return queryGateway(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

func queryAsset(context sdk.Context, query abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// TODO
	return nil, nil
}

// QueryGatewaysParams is the query parameters for 'custom/asset/gateways'
type QueryGatewaysParams struct {
	Owner sdk.AccAddress
}

func queryGateways(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryGatewaysParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var gateways []Gateway

	gatewaysIterator := keeper.GetGateways(ctx, params.Owner)
	defer gatewaysIterator.Close()

	for ; gatewaysIterator.Valid(); gatewaysIterator.Next() {
		gateway := Gateway{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(gatewaysIterator.Value(), &gateway)
		gateways = append(gateways, gateway)
	}

	if len(gateways) == 0 {
		return nil, nil
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, gateways)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// QueryGatewayParams is the query parameters for 'custom/asset/gateway'
type QueryGatewayParams struct {
	Moniker string
}

func queryGateway(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryGatewayParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	gateway, err := keeper.GetGatewayByMoniker(ctx, params.Moniker)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, gateway)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}
