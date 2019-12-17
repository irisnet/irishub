package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	"github.com/irisnet/irishub/modules/service"
)

// Rest variable names
// nolint
const (
	RestDefChainId  = "def-chain-id"
	RestServiceName = "service-name"
	RestBindChainId = "bind-chain-id"
	RestReqChainId  = "req-chain-id"
	RestReqId       = "request-id"
	RestProvider    = "provider"
	RestConsumer    = "consumer"
	RestAddress     = "address"
)

// RegisterRoutes defines routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type ServiceDefinitionReq struct {
	BaseReq            rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
	ServiceName        string       `json:"service_name" yaml:"service_name"`
	ServiceDescription string       `json:"service_description" yaml:"service_description"`
	AuthorDescription  string       `json:"author_description" yaml:"author_description"`
	Tags               []string     `json:"tags" yaml:"tags"`
	IdlContent         string       `json:"idl_content" yaml:"idl_content"`
	AuthorAddr         string       `json:"author_addr" yaml:"author_addr"`
}

type ServiceBindingReq struct {
	BaseReq     rest.BaseReq  `json:"base_req" yaml:"base_req"` // base req
	ServiceName string        `json:"service_name" yaml:"service_name"`
	DefChainId  string        `json:"def_chain_id" yaml:"def_chain_id"`
	BindingType string        `json:"binding_type" yaml:"binding_type"`
	Deposit     string        `json:"deposit" yaml:"deposit"`
	Prices      []string      `json:"prices" yaml:"prices"`
	Level       service.Level `json:"level" yaml:"level"`
	Provider    string        `json:"provider" yaml:"provider"`
}

type ServiceBindingUpdateReq struct {
	BaseReq     rest.BaseReq  `json:"base_req" yaml:"base_req"` // base req
	BindingType string        `json:"binding_type" yaml:"binding_type"`
	Deposit     string        `json:"deposit" yaml:"deposit"`
	Prices      []string      `json:"prices" yaml:"prices"`
	Level       service.Level `json:"level" yaml:"level"`
}

type ServiceBindingEnableReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
	Deposit string       `json:"deposit" yaml:"deposit"`
}

type ServiceRequest struct {
	ServiceName string `json:"service_name" yaml:"service_name"`
	BindChainId string `json:"bind_chain_id" yaml:"bind_chain_id"`
	DefChainId  string `json:"def_chain_id" yaml:"def_chain_id"`
	MethodId    int16  `json:"method_id" yaml:"method_id"`
	Provider    string `json:"provider" yaml:"provider"`
	Consumer    string `json:"consumer" yaml:"consumer"`
	ServiceFee  string `json:"service_fee" yaml:"service_fee"`
	Data        string `json:"data" yaml:"data"`
	Profiling   bool   `json:"profiling" yaml:"profiling"`
}

type ServiceRequestReq struct {
	BaseReq  rest.BaseReq     `json:"base_req" yaml:"base_req"` // base req
	Requests []ServiceRequest `json:"requests" yaml:"requests"`
}

type ServiceResponseReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
	ReqChainId string       `json:"req_chain_id" yaml:"req_chain_id"`
	RequestId  string       `json:"request_id" yaml:"request_id"`
	Data       string       `json:"data" yaml:"data"`
	Provider   string       `json:"provider" yaml:"provider"`
	ErrorMsg   string       `json:"error_msg" yaml:"error_msg"`
}

type BasicReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"` // base req
}
