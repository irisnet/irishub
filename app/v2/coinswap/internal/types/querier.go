package types

import "github.com/irisnet/irishub/types"

const (
	// QueryLiquidities liquidity query endpoint supported by the coinswap querier
	QueryLiquidities = "liquidities"
)

// QueryLiquidityParams is the query parameters for 'custom/swap/liquidity'
type QueryLiquidityParams struct {
	Id string
}

// QueryLiquidityResponse is the query response for 'custom/swap/liquidity'
type QueryLiquidityResponse struct {
	Iris      types.Coin `json:"iris"`
	Token     types.Coin `json:"token"`
	Liquidity types.Coin `json:"liquidity"`
	Fee       string     `json:"fee"`
}
