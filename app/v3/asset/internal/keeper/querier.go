package keeper

import (
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, k)
		case types.QueryTokens:
			return queryTokens(ctx, req, k)
		case types.QueryFees:
			return queryFees(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	if err := types.CheckSymbol(params.Symbol); err != nil {
		return nil, err
	}

	token, err := queryTokenBySymbol(ctx, keeper, strings.ToLower(params.Symbol))
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, token)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokensParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	tokens, err := queryTokensByOwner(ctx, keeper, params.Owner)
	if err != nil {
		return nil, err
	}

	bz, er := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if er != nil {
		return nil, sdk.MarshalResultErr(er)
	}

	return bz, nil
}

func queryFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenFeesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	if err := types.CheckSymbol(params.Symbol); err != nil {
		return nil, err
	}

	symbol := strings.ToLower(params.Symbol)
	issueFee := keeper.getTokenIssueFee(ctx, symbol)
	mintFee := keeper.getTokenMintFee(ctx, symbol)

	fees := types.TokenFeesOutput{
		Exist:    keeper.HasToken(ctx, symbol),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fees)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryTokenBySymbol(ctx sdk.Context, keeper Keeper, symbol string) (types.TokenOutput, sdk.Error) {
	if symbol == sdk.Iris {
		return types.NewTokenOutputFrom(getIrisToken()), nil
	}

	token, err := keeper.getToken(ctx, symbol)
	if err != nil {
		return types.TokenOutput{}, err
	}

	return types.NewTokenOutputFrom(token), nil
}

func queryTokensByOwner(ctx sdk.Context, keeper Keeper, owner sdk.AccAddress) (tokens types.TokensOutput, err sdk.Error) {
	if len(owner) == 0 {
		keeper.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
			tokens = append(tokens, types.NewTokenOutputFrom(token))
			return false
		})

		tokens = append(tokens, types.NewTokenOutputFrom(getIrisToken()))
		return
	}

	keeper.iterateTokensWithOwner(ctx, owner, func(token types.FungibleToken) (stop bool) {
		tokens = append(tokens, types.NewTokenOutputFrom(token))
		return false
	})

	return
}

func getIrisToken() types.FungibleToken {
	initSupply := uint64(sdk.InitialIssue.Int64())
	maxSupply := types.MaximumAssetMaxSupply

	return types.NewFungibleToken(sdk.Iris, sdk.IrisCoinType.Desc, sdk.IrisAtto, sdk.AttoScale, initSupply, maxSupply, true, sdk.AccAddress{})
}
