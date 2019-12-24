package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestHashLock = "hash-lock"
)

// RegisterRoutes registers HTLC-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// CreateHTLCReq defines the properties of a create HTLC request's body
type CreateHTLCReq struct {
	BaseTx               rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Sender               sdk.AccAddress `json:"sender" yaml:"sender"`
	To                   sdk.AccAddress `json:"to" yaml:"to"`
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain" yaml:"receiver_on_other_chain"`
	Amount               sdk.Coins      `json:"amount" yaml:"amount"`
	HashLock             string         `json:"hash_lock" yaml:"hash_lock"`
	TimeLock             uint64         `json:"time_lock" yaml:"time_lock"`
	Timestamp            uint64         `json:"timestamp" yaml:"timestamp"`
}

// ClaimHTLCReq defines the properties of a claim HTLC request's body
type ClaimHTLCReq struct {
	BaseTx rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	Secret string         `json:"secret" yaml:"secret"`
}

// RefundHTLCReq defines the properties of a refund HTLC request's body
type RefundHTLCReq struct {
	BaseTx rest.BaseReq   `json:"base_tx" yaml:"base_tx"`
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
}
