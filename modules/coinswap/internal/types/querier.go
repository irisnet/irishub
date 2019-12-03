package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// QueryLiquidity liquidity query endpoint supported by the coinswap querier
	QueryLiquidity = "liquidity"
)

// QueryLiquidityParams is the query parameters for 'custom/swap/liquidity'
type QueryLiquidityParams struct {
	Id string
}

// QueryLiquidityResponse is the query response for 'custom/swap/liquidity'
type QueryLiquidityResponse struct {
	Iris      sdk.Coin `json:"iris" yaml:"iris"`
	Token     sdk.Coin `json:"token" yaml:"token"`
	Liquidity sdk.Coin `json:"liquidity" yaml:"liquidity"`
	Fee       string   `json:"fee" yaml:"fee"`
}
