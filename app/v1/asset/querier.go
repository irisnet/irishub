package asset

import (
	"fmt"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryAsset      = "asset"
	QueryGateway    = "gateway"
	QueryGateways   = "gateways"
	QueryGatewayFee = "gatewayFee"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryAsset:
			return queryAsset(ctx, req, k)
		case QueryGateway:
			return queryGateway(ctx, req, k)
		case QueryGateways:
			return queryGateways(ctx, req, k)
		case QueryGatewayFee:
			return queryGatewayFee(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

// QueryAssetParams is the query parameters for 'custom/asset/asset'
type QueryAssetParams struct {
	Asset string
}

func queryAsset(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryAssetParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	asset, found := keeper.getAsset(ctx, GetKeyIDFromUniqueID(params.Asset))
	if !found {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("asset %s does not exist", params.Asset))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, asset)
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

	gateway, err2 := keeper.GetGateway(ctx, params.Moniker)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, gateway)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
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

	if len(params.Owner) != 0 {
		// if the owner provided
		gateways = queryGatewaysByOwner(ctx, params.Owner, keeper)
	} else {
		// if the owner not given
		gateways = queryAllGateways(ctx, keeper)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, gateways)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func queryGatewaysByOwner(ctx sdk.Context, owner sdk.AccAddress, keeper Keeper) []Gateway {
	var gateways []Gateway

	gatewaysIterator := keeper.GetGateways(ctx, owner)
	defer gatewaysIterator.Close()

	for ; gatewaysIterator.Valid(); gatewaysIterator.Next() {
		var moniker string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(gatewaysIterator.Value(), &moniker)

		gateway, err := keeper.GetGateway(ctx, moniker)
		if err != nil {
			continue
		}

		gateways = append(gateways, gateway)
	}

	return gateways
}

func queryAllGateways(ctx sdk.Context, keeper Keeper) []Gateway {
	var gateways []Gateway

	gatewaysIterator := keeper.GetAllGateways(ctx)
	defer gatewaysIterator.Close()

	for ; gatewaysIterator.Valid(); gatewaysIterator.Next() {
		var gateway Gateway
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(gatewaysIterator.Value(), &gateway)

		gateways = append(gateways, gateway)
	}

	return gateways
}

// QueryGatewayFeeParams is the query parameters for 'custom/asset/gatewayFee'
type QueryGatewayFeeParams struct {
	Moniker string
}

func queryGatewayFee(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryGatewayFeeParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	moniker := params.Moniker
	if len(moniker) < MinimumGatewayMonikerSize || len(moniker) > MaximumGatewayMonikerSize {
		return nil, ErrInvalidMoniker(keeper.Codespace(), fmt.Sprintf("the length of the moniker must be between [%d,%d]", MinimumGatewayMonikerSize, MaximumGatewayMonikerSize))
	}

	assetParams := keeper.GetParamSet(ctx)
	gatewayBaseFee := assetParams.CreateGatewayBaseFee
	fee := sdk.NewDec(int64(gatewayBaseFee)).Quo(calcGatewayFeeFactor(moniker))

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fee)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}
