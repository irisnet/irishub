package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/asset/internal/types"
)

// NewQuerier creates a new asset Querier instance
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
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
	}

	token, found := keeper.GetToken(ctx, params.TokenID)
	if !found {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("token %s does not exist", params.TokenID))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, token)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokensParams
	var err error
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
	}

	source := types.NATIVE
	owner := sdk.AccAddress{}
	nonSymbolTokenID := ""

	if len(params.Source) > 0 { // if source is specified
		source, err = types.AssetSourceFromString(params.Source)
		if err != nil {
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
		}
	}

	if len(params.Owner) > 0 && source == types.NATIVE { // ignore owner if source != NATIVE
		owner, err = sdk.AccAddressFromBech32(params.Owner)
		if err != nil {
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
		}
	}

	if len(params.Source) > 0 || len(params.Gateway) > 0 {
		nonSymbolTokenID, err = types.GetTokenID(source, "")
		if err != nil {
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
		}
	}

	var tokens types.Tokens

	// Add iris to the list
	//if source == types.NATIVE && owner.Empty() {
	//	initSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.InitialIssue).String())
	//	if err != nil {
	//		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	//	}
	//	maxSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.NewInt(int64(types.MaximumAssetMaxSupply))).String())
	//	if err != nil {
	//		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	//	}
	//	token := types.NewFungibleToken(types.NATIVE, "", sdk.Iris, sdk.IrisCoinType.Desc, sdk.AttoScale, "", sdk.IrisAtto, initSupply.Amount, maxSupply.Amount, true, sdk.AccAddress{})
	//	tokens = append(tokens, token)
	//}

	// Query from db
	iter := keeper.GetTokens(ctx, owner, nonSymbolTokenID)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var tokenID string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &tokenID)
		token, found := keeper.GetToken(ctx, tokenID)
		if !found {
			continue
		}

		tokens = append(tokens, token)
	}

	if len(tokens) == 0 {
		tokens = make([]types.FungibleToken, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tokens)

	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
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
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
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

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fees)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryParameters(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	params := keeper.GetParamSet(ctx)

	res, err := codec.MarshalJSONIndent(keeper.cdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
