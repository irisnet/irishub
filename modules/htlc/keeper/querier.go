package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/htlc/types"
)

// NewQuerier creates a new HTLC Querier instance
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryHTLC:
			return queryHTLC(ctx, req, k, legacyQuerierCdc)
		case types.QueryAssetSupply:
			return queryAssetSupply(ctx, req, k, legacyQuerierCdc)
		case types.QueryAssetSupplies:
			return queryAssetSupplies(ctx, k, legacyQuerierCdc)
		case types.QueryParameters:
			return queryParams(ctx, k, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query path: %s", types.ModuleName, path[0])
		}
	}
}

func queryHTLC(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryHTLCParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	htlc, found := k.GetHTLC(ctx, params.ID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownHTLC, params.ID.String())
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, htlc)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAssetSupply(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var requestParams types.QueryAssetSupplyParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &requestParams); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	assetSupply, found := k.GetAssetSupply(ctx, requestParams.Denom)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrAssetSupplyNotFound, string(requestParams.Denom))
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, assetSupply)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAssetSupplies(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	assets := k.GetAllAssetSupplies(ctx)
	if assets == nil {
		assets = []types.AssetSupply{}
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, assets)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
