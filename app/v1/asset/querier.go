package asset

import (
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryAsset = "asset"
	QueryGateway    = "gateway"
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

func queryGateway(context sdk.Context, query abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// TODO
	return nil, nil
}


