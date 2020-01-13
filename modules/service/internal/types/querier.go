package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryDefinition = "definition" // QueryDefinition
	QueryBinding    = "binding"    // QueryBinding
	QueryBindings   = "bindings"   // QueryBindings
	QueryRequests   = "requests"   // QueryRequests
	QueryResponse   = "response"   // QueryResponse
	QueryFees       = "fees"       // QueryFees
)

// QueryDefinitionParams
type QueryDefinitionParams struct {
	ServiceName string
}

// QueryBindingParams
type QueryBindingParams struct {
	DefChainID  string
	ServiceName string
	BindChainID string
	Provider    sdk.AccAddress
}

// QueryBindingsParams
type QueryBindingsParams struct {
	DefChainID  string
	ServiceName string
}

// QueryRequestsParams
type QueryRequestsParams struct {
	DefChainID  string
	ServiceName string
	BindChainID string
	Provider    sdk.AccAddress
}

// QueryResponseParams
type QueryResponseParams struct {
	ReqChainID string
	RequestID  string
}

// QueryFeesParams
type QueryFeesParams struct {
	Address sdk.AccAddress
}

// FeesOutput
type FeesOutput struct {
	ReturnedFee sdk.Coins `json:"returned_fee" yaml:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee" yaml:"incoming_fee"`
}
