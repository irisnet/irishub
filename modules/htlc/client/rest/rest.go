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
	RestID    = "id"
	RestDenom = "denom"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// CreateHTLCReq defines the properties of an HTLC creation request's body.
type CreateHTLCReq struct {
	BaseReq              rest.BaseReq `json:"base_req" yaml:"base_req"`
	Sender               string       `json:"sender" yaml:"sender"`
	To                   string       `json:"to" yaml:"to"`
	ReceiverOnOtherChain string       `json:"receiver_on_other_chain" yaml:"receiver_on_other_chain"`
	SenderOnOtherChain   string       `json:"sender_on_other_chain" yaml:"sender_on_other_chain"`
	Amount               sdk.Coins    `json:"amount" yaml:"amount"`
	HashLock             string       `json:"hash_lock" yaml:"hash_lock"`
	TimeLock             uint64       `json:"time_lock" yaml:"time_lock"`
	Timestamp            uint64       `json:"timestamp" yaml:"timestamp"`
	Transfer             bool         `json:"transfer" yaml:"transfer"`
}

// ClaimHTLCReq defines the properties of an HTLC claim request's body.
type ClaimHTLCReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Sender  string       `json:"sender" yaml:"sender"`
	Secret  string       `json:"secret" yaml:"secret"`
}
