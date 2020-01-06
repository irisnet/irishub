package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// Rest variable names
// nolint
const (
	RestDefChainID  = "def-chain-id"
	RestServiceName = "service-name"
	RestBindChainID = "bind-chain-id"
	RestReqChainID  = "req-chain-id"
	RestReqID       = "request-id"
	RestProvider    = "provider"
	RestConsumer    = "consumer"
	RestAddress     = "address"
)

// RegisterRoutes defines routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// DefineServiceReq defines the properties of a service definition request's body
type DefineServiceReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
	Name              string       `json:"name" yaml:"name"`                             // Name of the service
	Description       string       `json:"description" yaml:"description"`               // Description of the service
	Tags              []string     `json:"tags" yaml:"tags"`                             // Tags of the service
	Author            string       `json:"author" yaml:"author"`                         // Author of the service
	AuthorDescription string       `json:"author_description" yaml:"author_description"` // Description of the author of the service
	Schemas           string       `json:"schemas" yaml:"schemas"`                       // Interface JSON Schemas of the service
}

// ServiceBindingReq defines the properties of a service binding request's body
type ServiceBindingReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
	ServiceName string       `json:"service_name" yaml:"service_name"`
	DefChainID  string       `json:"def_chain_id" yaml:"def_chain_id"`
	BindingType string       `json:"binding_type" yaml:"binding_type"`
	Deposit     string       `json:"deposit" yaml:"deposit"`
	Prices      []string     `json:"prices" yaml:"prices"`
	Level       types.Level  `json:"level" yaml:"level"`
	Provider    string       `json:"provider" yaml:"provider"`
}

// ServiceBindingUpdateReq defines the properties of a service binding update request's body
type ServiceBindingUpdateReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
	BindingType string       `json:"binding_type" yaml:"binding_type"`
	Deposit     string       `json:"deposit" yaml:"deposit"`
	Prices      []string     `json:"prices" yaml:"prices"`
	Level       types.Level  `json:"level" yaml:"level"`
}

// ServiceBindingEnableReq defines the properties of a service binding enable request's body
type ServiceBindingEnableReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
	Deposit string       `json:"deposit" yaml:"deposit"`
}

// ServiceRequest defines the properties of ServiceRequest
type ServiceRequest struct {
	ServiceName string `json:"service_name" yaml:"service_name"`
	BindChainID string `json:"bind_chain_id" yaml:"bind_chain_id"`
	DefChainID  string `json:"def_chain_id" yaml:"def_chain_id"`
	MethodID    int16  `json:"method_id" yaml:"method_id"`
	Provider    string `json:"provider" yaml:"provider"`
	Consumer    string `json:"consumer" yaml:"consumer"`
	ServiceFee  string `json:"service_fee" yaml:"service_fee"`
	Data        string `json:"data" yaml:"data"`
	Profiling   bool   `json:"profiling" yaml:"profiling"`
}

// ServiceRequestReq defines the properties of a service request's body
type ServiceRequestReq struct {
	BaseReq  rest.BaseReq     `json:"base_req" yaml:"base_req"` // base req
	Requests []ServiceRequest `json:"requests" yaml:"requests"`
}

// ServiceResponseReq defines the properties of a service response request's body
type ServiceResponseReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
	ReqChainID string       `json:"req_chain_id" yaml:"req_chain_id"`
	RequestID  string       `json:"request_id" yaml:"request_id"`
	Data       string       `json:"data" yaml:"data"`
	Provider   string       `json:"provider" yaml:"provider"`
	ErrorMsg   string       `json:"error_msg" yaml:"error_msg"`
}

// BasicReq defines the properties of a basic request's body
type BasicReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
}
