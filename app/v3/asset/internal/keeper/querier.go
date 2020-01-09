package keeper

import (
	"github.com/irisnet/irishub/app/v3/asset/internal/types"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return querierToken(ctx, req, k)
		case types.QueryTokens:
			return querierTokens(ctx, req, k)
		case types.QueryFees:
			return queryFees(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

func querierToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}
	token, err := queryToken(ctx, keeper, params.TokenId)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, token)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func querierTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (bz []byte, err sdk.Error) {
	var params types.QueryTokensParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}
	var tokens []types.TokenOutput
	if len(params.TokenID) > 0 {
		token, err := queryToken(ctx, keeper, params.TokenID)
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		tokens = append(tokens, token)
	} else {
		tokens, err = queryTokens(ctx, keeper, params.Owner)
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
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

	symbol := params.Symbol
	issueFee := GetTokenIssueFee(ctx, keeper, symbol)
	mintFee := GetTokenMintFee(ctx, keeper, symbol)

	tokenID := types.GetTokenID(symbol)
	fees := types.TokenFeesOutput{
		Exist:    keeper.HasToken(ctx, tokenID),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fees)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryToken(ctx sdk.Context, keeper Keeper, tokenID string) (types.TokenOutput, sdk.Error) {
	if tokenID == sdk.Iris {
		return types.NewTokenOutputFrom(getIrisToken()), nil
	}
	token, err := keeper.getToken(ctx, tokenID)
	if err != nil {
		return types.TokenOutput{}, err
	}
	return types.NewTokenOutputFrom(token), nil
}

func queryTokens(ctx sdk.Context, keeper Keeper, owner string) (tokens types.TokensOutput, error sdk.Error) {
	if len(owner) == 0 {
		keeper.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
			tokens = append(tokens, types.NewTokenOutputFrom(token))
			return false
		})
		tokens = append(tokens, types.NewTokenOutputFrom(getIrisToken()))
		return
	}

	ownerAcc, er := sdk.AccAddressFromBech32(owner)
	if er != nil {
		return nil, sdk.ParseParamsErr(er)
	}
	keeper.IterateTokensWithOwner(ctx, ownerAcc, func(token types.FungibleToken) (stop bool) {
		tokens = append(tokens, types.NewTokenOutputFrom(token))
		return false
	})
	return
}

func getIrisToken() types.FungibleToken {
	initSupply, _ := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.InitialIssue).String())
	maxSupply, _ := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.NewInt(int64(types.MaximumAssetMaxSupply))).String())
	return types.NewFungibleToken(sdk.Iris, sdk.IrisCoinType.Desc, sdk.AttoScale, initSupply.Amount, maxSupply.Amount, true, sdk.AccAddress{})
}
