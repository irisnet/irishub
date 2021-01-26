package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/token/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, k, legacyQuerierCdc)
		case types.QueryTokens:
			return queryTokens(ctx, req, k, legacyQuerierCdc)
		case types.QueryFees:
			return queryFees(ctx, req, k, legacyQuerierCdc)
		case types.QueryParams:
			return queryParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryTotalBurn:
			return queryTotalBurn(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown token query endpoint")
		}
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTokenParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, err
	}

	token, err := keeper.GetToken(ctx, params.Denom)
	if err != nil {
		return nil, err
	}

	return codec.MarshalJSONIndent(legacyQuerierCdc, token)
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTokensParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, err
	}
	tokens := keeper.GetTokens(ctx, params.Owner)
	return codec.MarshalJSONIndent(legacyQuerierCdc, tokens)
}

func queryFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTokenFeesParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, err
	}

	if err := types.ValidateSymbol(params.Symbol); err != nil {
		return nil, err
	}

	issueFee := keeper.GetTokenIssueFee(ctx, params.Symbol)
	mintFee := keeper.GetTokenMintFee(ctx, params.Symbol)

	fees := types.QueryFeesResponse{
		Exist:    keeper.HasToken(ctx, params.Symbol),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}

	return codec.MarshalJSONIndent(legacyQuerierCdc, fees)
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := keeper.GetParamSet(ctx)
	return codec.MarshalJSONIndent(legacyQuerierCdc, params)
}

func queryTotalBurn(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	resp, err := keeper.TotalBurn(sdk.WrapSDKContext(ctx), &types.QueryTotalBurnRequest{})
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(legacyQuerierCdc, resp)
}
