package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryToken  = "token"
	QueryTokens = "tokens"
	QueryFees   = "fees"
	QueryParams = "params"
)

// QueryTokenParams is the query parameters for 'custom/token/token'
type QueryTokenParams struct {
	Denom string
}

// QueryTokensParams is the query parameters for 'custom/token/tokens'
type QueryTokensParams struct {
	Owner sdk.AccAddress
}

// QueryTokenFeesParams is the query parameters for 'custom/token/fees'
type QueryTokenFeesParams struct {
	Symbol string
}
