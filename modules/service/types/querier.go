package types

import (
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryDefinition       = "definition"       // query definition
	QueryBinding          = "binding"          // query binding
	QueryBindings         = "bindings"         // query bindings
	QueryWithdrawAddress  = "withdraw_address" // query withdrawal address
	QueryRequest          = "request"          // query request
	QueryRequests         = "requests"         // query requests
	QueryResponse         = "response"         // query response
	QueryRequestContext   = "context"          // query request context
	QueryRequestsByReqCtx = "requests_by_ctx"  // query requests by the request context
	QueryResponses        = "responses"        // query responses
	QueryEarnedFees       = "fees"             // query earned fees
	QuerySchema           = "schema"           // query schema
	QueryParameters       = "parameters"       // query parameters
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

// QueryBindingsParams defines the params to query all bindings of a service definition with an optional owner
type QueryBindingsParams struct {
	ServiceName string
	Owner       sdk.AccAddress
}

// QueryWithdrawAddressParams defines the params to query the withdrawal address of an owner
type QueryWithdrawAddressParams struct {
	Owner sdk.AccAddress
}

// QueryRequestParams defines the params to query the request by ID
type QueryRequestParams struct {
	RequestID []byte
}

// QueryRequestsParams defines the params to query active requests for a service binding
type QueryRequestsParams struct {
	ServiceName string
	Provider    sdk.AccAddress
}

// QueryResponseParams defines the params to query the response to a request
type QueryResponseParams struct {
	RequestID tmbytes.HexBytes
}

// QueryRequestContextParams defines the params to query the request context
type QueryRequestContextParams struct {
	RequestContextID tmbytes.HexBytes
}

// QueryRequestsByReqCtxParams defines the params to query active requests by the request context ID
type QueryRequestsByReqCtxParams struct {
	RequestContextID tmbytes.HexBytes
	BatchCounter     uint64
}

// QueryResponsesParams defines the params to query active responses by the request context ID
type QueryResponsesParams struct {
	RequestContextID tmbytes.HexBytes
	BatchCounter     uint64
}

// QueryEarnedFeesParams defines the params to query the earned fees for a provider
type QueryEarnedFeesParams struct {
	Provider sdk.AccAddress
}

// QuerySchemaParams defines the params to query the system schemas by the schema name
type QuerySchemaParams struct {
	SchemaName string
}
