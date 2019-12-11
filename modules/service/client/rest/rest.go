package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
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
	BaseReq            rest.BaseReq `json:"base_req"` // base req
	ServiceName        string       `json:"service_name"`
	ServiceDescription string       `json:"service_description"`
	AuthorDescription  string       `json:"author_description"`
	Tags               []string     `json:"tags"`
	IdlContent         string       `json:"idl_content"`
	AuthorAddr         string       `json:"author_addr"`
}

type ServiceBindingReq struct {
	BaseReq     rest.BaseReq  `json:"base_req"` // base req
	ServiceName string        `json:"service_name"`
	DefChainId  string        `json:"def_chain_id"`
	BindingType string        `json:"binding_type"`
	Deposit     string        `json:"deposit"`
	Prices      []string      `json:"prices"`
	Level       service.Level `json:"level"`
	Provider    string        `json:"provider"`
}

type ServiceBindingUpdateReq struct {
	BaseReq     rest.BaseReq  `json:"base_req"` // base req
	BindingType string        `json:"binding_type"`
	Deposit     string        `json:"deposit"`
	Prices      []string      `json:"prices"`
	Level       service.Level `json:"level"`
}

type ServiceBindingEnableReq struct {
	BaseReq rest.BaseReq `json:"base_req"` // base req
	Deposit string       `json:"deposit"`
}

type ServiceRequest struct {
	ServiceName string `json:"service_name"`
	BindChainId string `json:"bind_chain_id"`
	DefChainId  string `json:"def_chain_id"`
	MethodId    int16  `json:"method_id"`
	Provider    string `json:"provider"`
	Consumer    string `json:"consumer"`
	ServiceFee  string `json:"service_fee"`
	Data        string `json:"data"`
	Profiling   bool   `json:"profiling"`
}

type ServiceRequestReq struct {
	BaseReq  rest.BaseReq     `json:"base_req"` // base req
	Requests []ServiceRequest `json:"requests"`
}

type ServiceResponseReq struct {
	BaseReq    rest.BaseReq `json:"base_req"` // base req
	ReqChainId string       `json:"req_chain_id"`
	RequestId  string       `json:"request_id"`
	Data       string       `json:"data"`
	Provider   string       `json:"provider"`
	ErrorMsg   string       `json:"error_msg"`
}

type BasicReq struct {
	BaseReq rest.BaseReq `json:"base_req"` // base req
}
