package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// QueryToken - query token path
	QueryToken = "token"
	// QueryTokens - query tokens path
	QueryTokens = "tokens"
	// QueryFees - query fees path
	QueryFees = "fees"
	// QueryParameters - query parameters path
	QueryParameters = "parameters"
)

// QueryTokenParams is the query parameters for 'custom/asset/tokens/{symbol}'
type QueryTokenParams struct {
	Symbol string `json:"symbol" yaml:"symbol"` //
}

// QueryTokensParams is the query parameters for 'custom/asset/tokens'
type QueryTokensParams struct {
	Symbol string `json:"symbol" yaml:"symbol"` //
	Owner  string `json:"owner" yaml:"owner"`   //
}

// QueryTokenFeesParams is the query parameters for 'custom/asset/fees/tokens'
type QueryTokenFeesParams struct {
	Symbol string `json:"symbol" yaml:"symbol"` //
}

// TokenFeesOutput is for the token fees query output
type TokenFeesOutput struct {
	Exist    bool     `json:"exist" yaml:"exist"`         // indicate if the token has existed
	IssueFee sdk.Coin `json:"issue_fee" yaml:"issue_fee"` // issue fee
	MintFee  sdk.Coin `json:"mint_fee" yaml:"mint_fee"`   // mint fee
}
