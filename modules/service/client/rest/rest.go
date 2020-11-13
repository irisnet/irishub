package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestServiceName      = "service-name"
	RestRequestID        = "request-id"
	RestOwner            = "owner"
	RestProvider         = "provider"
	RestConsumer         = "consumer"
	RestRequestContextID = "request-context-id"
	RestBatchCounter     = "batch-counter"
	RestArg1             = "arg1"
	RestArg2             = "arg2"
	RestSchemaName       = "schema-name"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// DefineServiceReq defines the properties of a define service request's body.
type DefineServiceReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
	Name              string       `json:"name" yaml:"name"`
	Description       string       `json:"description" yaml:"description"`
	Tags              []string     `json:"tags" yaml:"tags"`
	Author            string       `json:"author" yaml:"author"`
	AuthorDescription string       `json:"author_description" yaml:"author_description"`
	Schemas           string       `json:"schemas" yaml:"schemas"`
}

// BindServiceReq defines the properties of a bind service request's body.
type BindServiceReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	ServiceName string       `json:"service_name" yaml:"service_name"`
	Provider    string       `json:"provider" yaml:"provider"`
	Deposit     string       `json:"deposit" yaml:"deposit"`
	Pricing     string       `json:"pricing" yaml:"pricing"`
	QoS         uint64       `json:"qos" yaml:"qos"`
	Options     string       `json:"options" yaml:"options"`
	Owner       string       `json:"owner" yaml:"owner"`
}

// UpdateServiceBindingReq defines the properties of an update service binding request's body.
type UpdateServiceBindingReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Deposit string       `json:"deposit" yaml:"deposit"`
	Pricing string       `json:"pricing" yaml:"pricing"`
	QoS     uint64       `json:"qos" yaml:"qos"`
	Options string       `json:"options" yaml:"options"`
	Owner   string       `json:"owner" yaml:"owner"`
}

// SetWithdrawAddrReq defines the properties of a set withdraw address request's body.
type SetWithdrawAddrReq struct {
	BaseReq         rest.BaseReq `json:"base_req" yaml:"base_req"`
	WithdrawAddress string       `json:"withdraw_address" yaml:"withdraw_address"`
}

// DisableServiceBindingReq defines the properties of a disable service binding request's body.
type DisableServiceBindingReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Owner   string       `json:"owner" yaml:"owner"`
}

// EnableServiceBindingReq defines the properties of an enable service binding request's body.
type EnableServiceBindingReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Deposit string       `json:"deposit" yaml:"deposit"`
	Owner   string       `json:"owner" yaml:"owner"`
}

// RefundServiceDepositReq defines the properties of a refund service deposit request's body.
type RefundServiceDepositReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Owner   string       `json:"owner" yaml:"owner"`
}

type callServiceReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Consumer          string       `json:"consumer"`
	Input             string       `json:"input"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	SuperMode         bool         `json:"super_mode"`
	Repeated          bool         `json:"repeated"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
}

type respondServiceReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	RequestID string       `json:"request_id"`
	Provider  string       `json:"provider"`
	Result    string       `json:"result"`
	Output    string       `json:"output"`
}

type pauseRequestContextReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Consumer string       `json:"consumer"`
}

type startRequestContextReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Consumer string       `json:"consumer"`
}

type killRequestContextReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Consumer string       `json:"consumer"`
}

type updateRequestContextReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	Providers         []string     `json:"providers"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	Consumer          string       `json:"consumer"`
}

type withdrawEarnedFeesReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Owner   string       `json:"owner" yaml:"owner"`
}
