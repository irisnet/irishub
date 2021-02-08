package types

const (
	// QueryLiquidity is the liquidity query endpoint supported by the coinswap querier
	QueryLiquidity = "liquidity"
)

// QueryLiquidityParams is the query parameters for 'custom/swap/liquidity'
type QueryLiquidityParams struct {
	Denom string `json:"denom" yaml:"denom"` // same as uniDenom
}
