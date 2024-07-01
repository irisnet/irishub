package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	v1 "mods.irisnet.org/modules/token/types/v1"
	"mods.irisnet.org/modules/token/types/v1beta1"
)

var _ v1beta1.QueryServer = legacyQueryServer{}

type legacyQueryServer struct {
	server v1.QueryServer
	cdc    codec.Codec
}

// NewLegacyQueryServer returns an implementation of the token QueryServer interface
// for the provided Keeper.
func NewLegacyQueryServer(server v1.QueryServer, cdc codec.Codec) v1beta1.QueryServer {
	return &legacyQueryServer{
		server: server,
		cdc:    cdc,
	}
}

func (q legacyQueryServer) Token(c context.Context, req *v1beta1.QueryTokenRequest) (*v1beta1.QueryTokenResponse, error) {
	res, err := q.server.Token(c, &v1.QueryTokenRequest{
		Denom: req.Denom,
	})
	if err != nil {
		return nil, err
	}

	v1beta1Token, err := v1TokenToV1beta1(q.cdc, res.Token)
	if err != nil {
		return nil, err
	}

	return &v1beta1.QueryTokenResponse{Token: v1beta1Token}, nil
}

func (q legacyQueryServer) Tokens(c context.Context, req *v1beta1.QueryTokensRequest) (*v1beta1.QueryTokensResponse, error) {
	res, err := q.server.Tokens(c, &v1.QueryTokensRequest{
		Owner:      req.Owner,
		Pagination: req.Pagination,
	})
	if err != nil {
		return nil, err
	}

	var tokens []*codectypes.Any
	for _, token := range res.Tokens {
		v1beta1Token, err := v1TokenToV1beta1(q.cdc, token)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, v1beta1Token)
	}
	return &v1beta1.QueryTokensResponse{Tokens: tokens, Pagination: res.Pagination}, nil
}

func (q legacyQueryServer) Fees(c context.Context, req *v1beta1.QueryFeesRequest) (*v1beta1.QueryFeesResponse, error) {
	res, err := q.server.Fees(c, &v1.QueryFeesRequest{
		Symbol: req.Symbol,
	})
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryFeesResponse{
		Exist:    res.Exist,
		IssueFee: res.IssueFee,
		MintFee:  res.MintFee,
	}, nil
}

// Params return the all the parameter in tonken module
func (q legacyQueryServer) Params(c context.Context, req *v1beta1.QueryParamsRequest) (*v1beta1.QueryParamsResponse, error) {
	res, err := q.server.Params(c, &v1.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryParamsResponse{
		Params: v1beta1.Params{
			TokenTaxRate:      res.Params.TokenTaxRate,
			IssueTokenBaseFee: res.Params.IssueTokenBaseFee,
			MintTokenFeeRatio: res.Params.MintTokenFeeRatio,
		},
		Res: res.Res,
	}, nil
}

// TotalBurn return the all burn coin
func (q legacyQueryServer) TotalBurn(c context.Context, req *v1beta1.QueryTotalBurnRequest) (*v1beta1.QueryTotalBurnResponse, error) {
	res, err := q.server.TotalBurn(c, &v1.QueryTotalBurnRequest{})
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryTotalBurnResponse{
		BurnedCoins: res.BurnedCoins,
	}, nil
}

func v1TokenToV1beta1(cdc codec.Codec, v1token *codectypes.Any) (*codectypes.Any, error) {
	var v1beta1Token v1beta1.Token
	if err := cdc.Unmarshal(v1token.GetValue(), &v1beta1Token); err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(&v1beta1Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return any, nil
}
