package types

import (
	sdk "github.com/irisnet/irishub/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	QueryDefinition       = "definition"       // QueryDefinition
	QueryBinding          = "binding"          // QueryBinding
	QueryBindings         = "bindings"         // QueryBindings
	QueryWithdrawAddress  = "withdraw_address" // QueryWithdrawAddress
	QueryRequest          = "request"          // QueryRequest
	QueryRequests         = "requests"         // QueryRequests
	QueryResponse         = "response"         // QueryResponse
	QueryRequestContext   = "context"          // QueryRequestContext
	QueryRequestsByReqCtx = "requests_by_ctx"  // QueryRequestsByReqCtx
	QueryResponses        = "responses"        // QueryResponses
	QueryEarnedFees       = "fees"             // QueryEarnedFees
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

// QueryBindingsParams defines the params to query all bindings of a service definition
type QueryBindingsParams struct {
	ServiceName string
}

// QueryWithdrawAddressParams defines the params to query the withdrawal address of a provider
type QueryWithdrawAddressParams struct {
	Provider sdk.AccAddress
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
	RequestID cmn.HexBytes
}

// QueryRequestContextParams defines the params to query the request context
type QueryRequestContextParams struct {
	RequestContextID cmn.HexBytes
}

// QueryRequestsByReqCtxParams defines the params to query active requests by the request context ID
type QueryRequestsByReqCtxParams struct {
	RequestContextID cmn.HexBytes
	BatchCounter     uint64
}

// QueryResponsesParams defines the params to query active responses by the request context ID
type QueryResponsesParams struct {
	RequestContextID cmn.HexBytes
	BatchCounter     uint64
}

// QueryEarnedFeesParams defines the params to query the earned fees for a provider
type QueryEarnedFeesParams struct {
	Provider sdk.AccAddress
}
