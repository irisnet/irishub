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
	ID string `json:"id" yaml:"id"` // same as uniDenom
}

// QueryLiquidityResponse is the query response for 'custom/swap/liquidity'
type QueryLiquidityResponse struct {
	Standard  sdk.Coin `json:"standard" yaml:"standard"`   // standard token
	Token     sdk.Coin `json:"token" yaml:"token"`         // the other token in swap pool
	Liquidity sdk.Coin `json:"liquidity" yaml:"liquidity"` // liquidity of swap pool
	Fee       string   `json:"fee" yaml:"fee"`             // fee of swap pool
}
