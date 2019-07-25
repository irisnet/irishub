package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset/internal/types"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, k)
		case types.QueryTokens:
			return queryTokens(ctx, req, k)
		case types.QueryGateway:
			return queryGateway(ctx, req, k)
		case types.QueryGateways:
			return queryGateways(ctx, req, k)
		case types.QueryFees:
			return queryFees(ctx, path[1:], req, k)
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

	var token types.FungibleToken
	if params.TokenId == sdk.Iris {
		initSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.InitialIssue).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		maxSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.NewInt(int64(types.MaximumAssetMaxSupply))).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		token = types.NewFungibleToken(types.NATIVE, "", sdk.Iris, sdk.IrisCoinType.Desc, sdk.AttoScale, "", sdk.IrisAtto, initSupply.Amount, maxSupply.Amount, true, sdk.AccAddress{})
	} else {
		var found bool
		token, found = keeper.getToken(ctx, params.TokenId)
		if !found {
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("token %s does not exist", params.TokenId))
		}

		if token.Source == types.GATEWAY {
			gateway, _ := keeper.GetGateway(ctx, token.Gateway)
			token.Owner = gateway.Owner
		}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, token)

	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokensParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	source := types.NATIVE
	gateway := ""
	owner := sdk.AccAddress{}
	nonSymbolTokenId := ""

	if len(params.Source) > 0 { // if source is specified
		source, err = types.AssetSourceFromString(params.Source)
		if err != nil {
			return nil, sdk.ParseParamsErr(err)
		}
	} else if len(params.Gateway) > 0 { // if source is not specified, and gateway is specified
		source = types.GATEWAY
	}

	if source == types.GATEWAY { // ignore gateway moniker if source != GATEWAY
		gateway = params.Gateway
		if len(gateway) == 0 {
			return nil, sdk.ErrUnknownRequest("gateway moniker is required for querying gateway tokens")
		}
	}

	if len(params.Owner) > 0 && source == types.NATIVE { // ignore owner if source != NATIVE
		owner, err = sdk.AccAddressFromBech32(params.Owner)
		if err != nil {
			return nil, sdk.ParseParamsErr(err)
		}
	}

	if len(params.Source) > 0 || len(params.Gateway) > 0 {
		nonSymbolTokenId, err = types.GetTokenID(source, "", gateway)
	}

	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var tokens types.Tokens

	// Add iris to the list
	if source == types.NATIVE && owner.Empty() {
		initSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.InitialIssue).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		maxSupply, err := sdk.IrisCoinType.ConvertToMinDenomCoin(sdk.NewCoin(sdk.Iris, sdk.NewInt(int64(types.MaximumAssetMaxSupply))).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		token := types.NewFungibleToken(types.NATIVE, "", sdk.Iris, sdk.IrisCoinType.Desc, sdk.AttoScale, "", sdk.IrisAtto, initSupply.Amount, maxSupply.Amount, true, sdk.AccAddress{})
		tokens = append(tokens, token)
	}

	// Query from db
	iter := keeper.getTokens(ctx, owner, nonSymbolTokenId)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var tokenId string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &tokenId)
		token, found := keeper.getToken(ctx, tokenId)
		if !found {
			continue
		}

		if token.Source == types.GATEWAY {
			gateway, _ := keeper.GetGateway(ctx, token.Gateway)
			token.Owner = gateway.Owner
		}

		tokens = append(tokens, token)
	}

	if len(tokens) == 0 {
		tokens = make([]types.FungibleToken, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tokens)

	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func queryGateway(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryGatewayParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	if err := types.ValidateMoniker(params.Moniker); err != nil {
		return nil, err
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

func queryGateways(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryGatewaysParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var gateways []types.Gateway

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

func queryGatewaysByOwner(ctx sdk.Context, owner sdk.AccAddress, keeper Keeper) []types.Gateway {
	var gateways = make([]types.Gateway, 0)

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

func queryAllGateways(ctx sdk.Context, keeper Keeper) []types.Gateway {
	var gateways = make([]types.Gateway, 0)

	keeper.IterateGateways(ctx, func(gw types.Gateway) (stop bool) {
		gateways = append(gateways, gw)
		return false
	})

	return gateways
}

func queryFees(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	switch path[0] {
	case "gateways":
		return queryGatewayFee(ctx, req, keeper)
	case "tokens":
		return queryTokenFees(ctx, req, keeper)
	default:
		return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
	}
}

func queryGatewayFee(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryGatewayFeeParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	moniker := params.Moniker
	if err := types.ValidateMoniker(moniker); err != nil {
		return nil, err
	}

	fee := types.GatewayFeeOutput{
		Exist: keeper.HasGateway(ctx, moniker),
		Fee:   GetGatewayCreateFee(ctx, keeper, moniker),
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fee)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryTokenFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenFeesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
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
	} else {
		issueFee = GetGatewayTokenIssueFee(ctx, keeper, symbol)
		mintFee = GetGatewayTokenMintFee(ctx, keeper, symbol)
	}

	fees := types.TokenFeesOutput{
		Exist:    keeper.HasToken(ctx, id),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fees)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}
