package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	QueryDefinition = "definition" // QueryDefinition
	QueryBinding    = "binding"    // QueryBinding
	QueryBindings   = "bindings"   // QueryBindings
	QueryRequests   = "requests"   // QueryRequests
	QueryResponse   = "response"   // QueryResponse
	QueryFees       = "fees"       // QueryFees
)

// QueryDefinitionParams defines the params to query a service definition
type QueryDefinitionParams struct {
	ServiceName string
}

// QueryBindingParams defines the params to query a service binding
type QueryBindingParams struct {
	ServiceName string
	Provider    sdk.AccAddress
}

// QueryBindingsParams defines the params to query all service bindings of a service
type QueryBindingsParams struct {
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
	ReturnedFee sdk.Coins `json:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee"`
}
