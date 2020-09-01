package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestRequestID = "request-id"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// RequestRandomReq defines the properties of a request random request's body
type RequestRandomReq struct {
	BaseReq       rest.BaseReq   `json:"base_req" yaml:"base_req"` // base req
	Consumer      sdk.AccAddress `json:"consumer"`                 // request address
	BlockInterval uint64         `json:"block_interval"`           // block interval
	Oracle        bool           `json:"oracle"`                   // oracle method
	ServiceFeeCap sdk.Coins      `json:"service_fee_cap"`          // service fee cap
}
