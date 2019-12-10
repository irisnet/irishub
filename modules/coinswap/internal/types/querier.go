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
	Id string `json:"id" yaml:"id"` // same as uniDenom
}

// QueryLiquidityResponse is the query response for 'custom/swap/liquidity'
type QueryLiquidityResponse struct {
	Standard  sdk.Coin `json:"standard" yaml:"standard"`
	Token     sdk.Coin `json:"token" yaml:"token"`
	Liquidity sdk.Coin `json:"liquidity" yaml:"liquidity"`
	Fee       string   `json:"fee" yaml:"fee"`
}
