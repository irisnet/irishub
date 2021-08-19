package types

const (
	// QueryPool is the liquidity query endpoint supported by the coinswap querier
	QueryPool = "pool"
	// QueryPools is the liquidity query endpoint supported by the coinswap querier
	QueryPools = "pools"
)

// QueryPoolParams is the query parameters for 'custom/swap/liquidity'
type QueryPoolParams struct {
	LptDenom string `json:"lpt-denom" yaml:"lpt-denom"` // same as uniDenom
}
