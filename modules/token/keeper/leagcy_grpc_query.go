package keeper

import (
	"context"

	v1 "github.com/irisnet/irismod/modules/token/types/v1"
	"github.com/irisnet/irismod/modules/token/types/v1beta1"
)

var _ v1beta1.QueryServer = leagcyQueryServer{}

type leagcyQueryServer struct {
	server v1.QueryServer
}

// NewLeagcyQueryServer returns an implementation of the token QueryServer interface
// for the provided Keeper.
func NewLeagcyQueryServer(server v1.QueryServer) v1beta1.QueryServer {
	return &leagcyQueryServer{
		server: server,
	}
}

func (q leagcyQueryServer) Token(c context.Context, req *v1beta1.QueryTokenRequest) (*v1beta1.QueryTokenResponse, error) {
	res, err := q.server.Token(c, &v1.QueryTokenRequest{
		Denom: req.Denom,
	})
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryTokenResponse{Token: res.Token}, nil
}

func (q leagcyQueryServer) Tokens(c context.Context, req *v1beta1.QueryTokensRequest) (*v1beta1.QueryTokensResponse, error) {
	res, err := q.server.Tokens(c, &v1.QueryTokensRequest{
		Owner:      req.Owner,
		Pagination: req.Pagination,
	})
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryTokensResponse{Tokens: res.Tokens, Pagination: res.Pagination}, nil
}

func (q leagcyQueryServer) Fees(c context.Context, req *v1beta1.QueryFeesRequest) (*v1beta1.QueryFeesResponse, error) {
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
func (q leagcyQueryServer) Params(c context.Context, req *v1beta1.QueryParamsRequest) (*v1beta1.QueryParamsResponse, error) {
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
func (q leagcyQueryServer) TotalBurn(c context.Context, req *v1beta1.QueryTotalBurnRequest) (*v1beta1.QueryTotalBurnResponse, error) {
	res, err := q.server.TotalBurn(c, &v1.QueryTotalBurnRequest{})
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryTotalBurnResponse{
		BurnedCoins: res.BurnedCoins,
	}, nil
}
