package types

const (
	// QueryLiquidity liquidity query endpoint supported by the coinswap querier
	QueryLiquidity = "liquidity"
	// QueryParameters parameters query endpoint supported by the coinswap querier
	QueryParameters = "parameters"
	// ParamFee fee query endpoint supported by the coinswap querier
	ParamFee = "fee"
)

// QueryLiquidityParams is the query parameters for 'custom/swap/liquidity/{id}'
type QueryLiquidityParams struct {
	TokenId string
}