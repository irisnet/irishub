package asset

import (
	"fmt"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryToken    = "token"
	QueryTokens   = "tokens"
	QueryGateway  = "gateway"
	QueryGateways = "gateways"
	QueryFees     = "fees"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryToken:
			return queryToken(ctx, req, k)
		case QueryTokens:
			return queryTokens(ctx, req, k)
		case QueryGateway:
			return queryGateway(ctx, req, k)
		case QueryGateways:
			return queryGateways(ctx, req, k)
		case QueryFees:
			return queryFees(ctx, path[1:], req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown asset query endpoint")
		}
	}
}

// QueryTokenParams is the query parameters for 'custom/asset/tokens/{id}'
type QueryTokenParams struct {
	TokenId string
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTokenParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var token FungibleToken
	if params.TokenId == sdk.NativeTokenName {
		initSupply, err := sdk.IRIS.ConvertToMinCoin(sdk.NewCoin(sdk.NativeTokenName, sdk.InitialIssue).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		maxSupply, err := sdk.IRIS.ConvertToMinCoin(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(int64(MaximumAssetMaxSupply))).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		token = NewFungibleToken(NATIVE, "", sdk.IRIS.GetMainUnit().Denom, sdk.IRIS.Desc, uint8(sdk.IRIS.GetMinUnit().Decimal), "", sdk.IRIS.GetMinUnit().Denom, initSupply.Amount, maxSupply.Amount, true, sdk.AccAddress{})
	} else {
		var found bool
		token, found = keeper.getToken(ctx, params.TokenId)
		if !found {
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("token %s does not exist", params.TokenId))
		}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, token)

	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// QueryTokensParams is the query parameters for 'custom/asset/tokens'
type QueryTokensParams struct {
	Source  string
	Gateway string
	Owner   string
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTokensParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	source := NATIVE
	gateway := ""
	owner := sdk.AccAddress{}
	nonSymbolTokenId := ""

	if len(params.Source) > 0 { // if source is specified
		source, err = AssetSourceFromString(params.Source)
		if err != nil {
			return nil, sdk.ParseParamsErr(err)
		}
	} else if len(params.Gateway) > 0 { // if source is not specified, and gateway is specified
		source = GATEWAY
	}

	if source == GATEWAY { // ignore gateway moniker if source != GATEWAY
		gateway = params.Gateway
		if len(gateway) == 0 {
			return nil, sdk.ErrUnknownRequest("gateway moniker is required for querying gateway tokens")
		}
	}

	if len(params.Owner) > 0 && source != EXTERNAL { // ignore owner if source == EXTERNAL
		owner, err = sdk.AccAddressFromBech32(params.Owner)
		if err != nil {
			return nil, sdk.ParseParamsErr(err)
		}
	}

	if len(params.Source) > 0 || len(params.Gateway) > 0 {
		nonSymbolTokenId, err = GetTokenID(source, "", gateway)
	}

	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var tokens []FungibleToken

	// Add iris to the list
	if source == NATIVE && owner.Empty() {
		initSupply, err := sdk.IRIS.ConvertToMinCoin(sdk.NewCoin(sdk.NativeTokenName, sdk.InitialIssue).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		maxSupply, err := sdk.IRIS.ConvertToMinCoin(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(int64(MaximumAssetMaxSupply))).String())
		if err != nil {
			return nil, sdk.MarshalResultErr(err)
		}
		token := NewFungibleToken(NATIVE, "", sdk.IRIS.GetMainUnit().Denom, sdk.IRIS.Desc, uint8(sdk.IRIS.GetMinUnit().Decimal), "", sdk.IRIS.GetMinUnit().Denom, initSupply.Amount, maxSupply.Amount, true, sdk.AccAddress{})
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
		tokens = append(tokens, token)
	}

	if len(tokens) == 0 {
		tokens = make([]FungibleToken, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tokens)

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

	if err := ValidateMoniker(params.Moniker); err != nil {
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
	var gateways = make([]Gateway, 0)

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
	var gateways = make([]Gateway, 0)

	gatewaysIterator := keeper.GetAllGateways(ctx)
	defer gatewaysIterator.Close()

	for ; gatewaysIterator.Valid(); gatewaysIterator.Next() {
		var gateway Gateway
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(gatewaysIterator.Value(), &gateway)

		gateways = append(gateways, gateway)
	}

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

// QueryGatewayFeeParams is the query parameters for 'custom/asset/fees/gateways'
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
	if err := ValidateMoniker(moniker); err != nil {
		return nil, err
	}

	fee := GatewayFeeOutput{
		Exist: keeper.HasGateway(ctx, moniker),
		Fee:   getGatewayCreateFee(ctx, keeper, moniker),
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fee)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

// QueryTokenFeesParams is the query parameters for 'custom/asset/fees/tokens'
type QueryTokenFeesParams struct {
	ID string
}

func queryTokenFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryTokenFeesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	id := params.ID
	if err := CheckAssetID(id); err != nil {
		return nil, err
	}

	source, symbol := ParseAssetID(id)

	var (
		issueFee sdk.Coin
		mintFee  sdk.Coin
	)

	if source == "" || source == "x" {
		issueFee = getTokenIssueFee(ctx, keeper, symbol)
		mintFee = getTokenMintFee(ctx, keeper, symbol)
	} else {
		issueFee = getGatewayTokenIssueFee(ctx, keeper, symbol)
		mintFee = getGatewayTokenMintFee(ctx, keeper, symbol)
	}

	fees := TokenFeesOutput{
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
