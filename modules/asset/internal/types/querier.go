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

// QueryTokenParams is the query parameters for 'custom/asset/tokens/{id}'
type QueryTokenParams struct {
	TokenID string `json:"token_id" yaml:"token_id"`
}

// QueryTokensParams is the query parameters for 'custom/asset/tokens'
type QueryTokensParams struct {
	Source  string `json:"source" yaml:"source"`
	Gateway string `json:"gateway" yaml:"gateway"`
	Owner   string `json:"owner" yaml:"owner"`
}

// QueryGatewayParams is the query parameters for 'custom/asset/gateway'
type QueryGatewayParams struct {
	Moniker string `json:"moniker" yaml:"moniker"`
}

// QueryGatewaysParams is the query parameters for 'custom/asset/gateways'
type QueryGatewaysParams struct {
	Owner sdk.AccAddress `json:"owner" yaml:"owner"`
}

// QueryGatewayFeeParams is the query parameters for 'custom/asset/fees/gateways'
type QueryGatewayFeeParams struct {
	Moniker string `json:"moniker" yaml:"moniker"`
}

// QueryTokenFeesParams is the query parameters for 'custom/asset/fees/tokens'
type QueryTokenFeesParams struct {
	ID string `json:"id" yaml:"id"`
}

// GatewayFeeOutput is for the gateway fee query output
type GatewayFeeOutput struct {
	Exist bool     `json:"exist" yaml:"exist"` // indicate if the gateway has existed
	Fee   sdk.Coin `json:"fee" yaml:"fee"`     // creation fee
}

// TokenFeesOutput is for the token fees query output
type TokenFeesOutput struct {
	Exist    bool     `json:"exist" yaml:"exist"`         // indicate if the token has existed
	IssueFee sdk.Coin `json:"issue_fee" yaml:"issue_fee"` // issue fee
	MintFee  sdk.Coin `json:"mint_fee" yaml:"mint_fee"`   // mint fee
}
