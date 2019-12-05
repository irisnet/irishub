package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/asset/internal/types"
	iristypes "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, k)
		case types.QueryTokens:
			return queryTokens(ctx, req, k)
		case types.QueryFees:
			return queryFees(ctx, path[1:], req, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, iristypes.ParseParamsErr(err)
	}

	var token types.FungibleToken

	var found bool
	token, found = keeper.GetToken(ctx, params.TokenID)
	if !found {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("token %s does not exist", params.TokenID))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, token)

	if err != nil {
		return nil, iristypes.MarshalResultErr(err)
	}
	return bz, nil
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokensParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, iristypes.ParseParamsErr(err)
	}

	source := types.NATIVE
	owner := sdk.AccAddress{}
	nonSymbolTokenId := ""

	if len(params.Source) > 0 { // if source is specified
		source, err = types.AssetSourceFromString(params.Source)
		if err != nil {
			return nil, iristypes.ParseParamsErr(err)
		}
	}

	if len(params.Owner) > 0 && source == types.NATIVE { // ignore owner if source != NATIVE
		owner, err = sdk.AccAddressFromBech32(params.Owner)
		if err != nil {
			return nil, iristypes.ParseParamsErr(err)
		}
	}

	if len(params.Source) > 0 || len(params.Gateway) > 0 {
		nonSymbolTokenId, err = types.GetTokenID(source, "")
	}

	if err != nil {
		return nil, iristypes.ParseParamsErr(err)
	}

	var tokens types.Tokens

	// Add iris to the list
	//if source == types.NATIVE && owner.Empty() {
	//	initSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.InitialIssue).String())
	//	if err != nil {
	//		return nil, iristypes.MarshalResultErr(err)
	//	}
	//	maxSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.NewInt(int64(types.MaximumAssetMaxSupply))).String())
	//	if err != nil {
	//		return nil, iristypes.MarshalResultErr(err)
	//	}
	//	token := types.NewFungibleToken(types.NATIVE, "", sdk.Iris, sdk.IrisCoinType.Desc, sdk.AttoScale, "", sdk.IrisAtto, initSupply.Amount, maxSupply.Amount, true, sdk.AccAddress{})
	//	tokens = append(tokens, token)
	//}

	// Query from db
	iter := keeper.GetTokens(ctx, owner, nonSymbolTokenId)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var tokenId string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &tokenId)
		token, found := keeper.GetToken(ctx, tokenId)
		if !found {
			continue
		}

		tokens = append(tokens, token)
	}

	if len(tokens) == 0 {
		tokens = make([]types.FungibleToken, 0)
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, tokens)

	if err != nil {
		return nil, iristypes.MarshalResultErr(err)
	}
	return bz, nil
}

func queryFees(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	switch path[0] {
	case "tokens":
		return queryTokenFees(ctx, req, keeper)
	default:
		return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
	}
}

func queryTokenFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenFeesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, iristypes.ParseParamsErr(err)
	}

	id := params.ID
	if err := types.CheckTokenID(id); err != nil {
		return nil, err
	}

	prefix, symbol := types.GetTokenIDParts(id)

	var (
		issueFee sdk.Coin
		mintFee  sdk.Coin
	)

	if prefix == "x" {
		return nil, sdk.ErrUnknownRequest("unsupported token source: external")
	} else if prefix == "" || prefix == "i" {
		issueFee = GetTokenIssueFee(ctx, keeper, symbol)
		mintFee = GetTokenMintFee(ctx, keeper, symbol)
	}

	fees := types.TokenFeesOutput{
		Exist:    keeper.HasToken(ctx, id),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, fees)
	if err != nil {
		return nil, iristypes.MarshalResultErr(err)
	}

	return bz, nil
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParamSet(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
