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

// QueryRequestsParams defines the params to query all requests for a service binding
type QueryRequestsParams struct {
	ServiceName string
	Provider    sdk.AccAddress
}

// QueryResponseParams defines the params to query the response for a request
type QueryResponseParams struct {
	RequestID string
}

// QueryFeesParams defines the params to query the earned fees for a provider
type QueryFeesParams struct {
	Address sdk.AccAddress
}