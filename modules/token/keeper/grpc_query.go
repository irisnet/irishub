package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	gogotypes "github.com/cosmos/gogoproto/types"

	"irismod.io/token/types"
	v1 "irismod.io/token/types/v1"
)

var _ v1.QueryServer = Keeper{}

// Token queries a token by denomination.
//
// Parameters:
// - c: Context object
// - req: QueryTokenRequest object
//
// Returns:
// - QueryTokenResponse object containing token
// - Error if any
func (k Keeper) Token(c context.Context, req *v1.QueryTokenRequest) (*v1.QueryTokenResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	token, err := k.GetToken(ctx, req.Denom)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "token %s not found", req.Denom)
	}
	msg, ok := token.(proto.Message)
	if !ok {
		return nil, status.Errorf(codes.Internal, "can't protomarshal %T", token)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.QueryTokenResponse{Token: any}, nil
}

// Tokens queries a list of tokens based on the given request parameters.
//
// Parameters:
// - c: Context object
// - req: QueryTokensRequest object
//
// Returns:
// - QueryTokensResponse object containing all tokens own by the owner
// - Error if any
func (k Keeper) Tokens(c context.Context, req *v1.QueryTokensRequest) (*v1.QueryTokensResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var owner sdk.AccAddress
	var err error
	if len(req.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("invalid owner address (%s)", err),
			)
		}
	}

	var tokens []v1.TokenI
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	if owner == nil {
		tokenStore := prefix.NewStore(store, types.PrefixTokenForSymbol)
		pageRes, err = query.Paginate(
			tokenStore,
			shapePageRequest(req.Pagination),
			func(_ []byte, value []byte) error {
				var token v1.Token
				k.cdc.MustUnmarshal(value, &token)
				tokens = append(tokens, &token)
				return nil
			},
		)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	} else {
		tokenStore := prefix.NewStore(store, types.KeyTokens(owner, ""))
		pageRes, err = query.Paginate(tokenStore, shapePageRequest(req.Pagination), func(_ []byte, value []byte) error {
			var symbol gogotypes.StringValue
			k.cdc.MustUnmarshal(value, &symbol)
			token, err := k.GetToken(ctx, symbol.Value)
			if err == nil {
				tokens = append(tokens, token)
			}
			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	}
	result := make([]*codectypes.Any, len(tokens))
	for i, token := range tokens {
		msg, ok := token.(proto.Message)
		if !ok {
			return nil, status.Errorf(codes.Internal, "%T does not implement proto.Message", token)
		}

		var err error
		if result[i], err = codectypes.NewAnyWithValue(msg); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &v1.QueryTokensResponse{Tokens: result, Pagination: pageRes}, nil
}

// Fees retrieves the issue fee and mint fee for a specific token symbol.
//
// Parameters:
// - c: Context object
// - req: QueryFeesRequest object containing the token symbol
//
// Returns:
// - QueryFeesResponse object containing issue fee, mint fee, and token existence status
// - Error if any
func (k Keeper) Fees(c context.Context, req *v1.QueryFeesRequest) (*v1.QueryFeesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := types.ValidateSymbol(req.Symbol); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	issueFee, err := k.GetTokenIssueFee(ctx, req.Symbol)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	mintFee, err := k.GetTokenMintFee(ctx, req.Symbol)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.QueryFeesResponse{
		Exist:    k.HasToken(ctx, req.Symbol),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}, nil
}

// Params returns all the parameters in the token module.
//
// Parameters:
// - c: Context object
// - req: QueryParamsRequest object
//
// Returns:
// - QueryParamsResponse object containing token params
// - Error if any
func (k Keeper) Params(c context.Context, req *v1.QueryParamsRequest) (*v1.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &v1.QueryParamsResponse{Params: params}, nil
}

// TotalBurn return the all burn coin
//
// Parameters:
// - c: Context object
// - req: QueryFeesRequest object
//
// Returns:
// - QueryTotalBurnResponse object containing token params
// - Error if any
func (k Keeper) TotalBurn(c context.Context, req *v1.QueryTotalBurnRequest) (*v1.QueryTotalBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &v1.QueryTotalBurnResponse{
		BurnedCoins: k.GetAllBurnCoin(ctx),
	}, nil
}

// Balances retrieves the balances of a given address for a specific token.
//
// Parameters:
// - c: the context.Context object.
// - req: the v1.QueryBalancesRequest object containing the address and token denomination.
//
// Returns:
// - *v1.QueryBalancesResponse: the response containing the balances of the address for the specified token.
// - error: an error if the request is empty, the address is invalid, or the token is not found.
func (k Keeper) Balances(c context.Context, req *v1.QueryBalancesRequest) (*v1.QueryBalancesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address (%s)", err)
	}

	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasToken(ctx, req.Denom) {
		balance := k.bankKeeper.GetBalance(ctx, addr, req.Denom)
		balances := sdk.NewCoins(balance)
		return &v1.QueryBalancesResponse{Balances: balances}, nil
	}

	token, err := k.GetToken(ctx, req.Denom)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "token %s not found", req.Denom)
	}

	balance := k.bankKeeper.GetBalance(ctx, addr, token.GetMinUnit())
	balances := sdk.NewCoins(balance)

	if len(token.GetContract()) > 0 {
		contract := common.HexToAddress(token.GetContract())
		account := common.BytesToAddress(addr.Bytes())

		erc20Balance, err := k.BalanceOf(ctx, contract, account)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		balances = balances.Add(sdk.NewCoin("erc20/"+token.GetContract(), sdkmath.NewIntFromBigInt(erc20Balance)))
	}
	return &v1.QueryBalancesResponse{Balances: balances}, nil
}
