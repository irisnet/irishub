package types

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	QueryDefinition       = "definition"      // QueryDefinition
	QueryBinding          = "binding"         // QueryBinding
	QueryBindings         = "bindings"        // QueryBindings
	QueryRequest          = "request"         // QueryRequest
	QueryRequests         = "requests"        // QueryRequests
	QueryResponse         = "response"        // QueryResponse
	QueryRequestContext   = "context"         // QueryRequestContext
	QueryRequestsByReqCtx = "requests_by_ctx" // QueryRequestsByReqCtx
	QueryResponses        = "responses"       // QueryResponses
	QueryFees             = "fees"            // QueryFees
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

// QueryRequestParams defines the params to query the request by ID
type QueryRequestParams struct {
	RequestID string
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

type QueryRequestContextParams struct {
	RequestContextID []byte
}

type QueryRequestsByReqCtxParams struct {
	RequestContextID []byte
	BatchCounter     uint64
}

type QueryResponsesParams struct {
	RequestContextID []byte
	BatchCounter     uint64
}

// QueryFeesParams defines the params to query the earned fees for a provider
type QueryFeesParams struct {
	Address sdk.AccAddress
}
