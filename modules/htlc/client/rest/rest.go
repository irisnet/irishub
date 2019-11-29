package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes registers HTLC-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type createHTLCReq struct {
	BaseTx               rest.BaseReq   `json:"base_tx"`
	Sender               sdk.AccAddress `json:"sender"`
	To                   sdk.AccAddress `json:"to"`
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain"`
	Amount               sdk.Coins      `json:"amount"`
	HashLock             string         `json:"hash_lock"`
	TimeLock             uint64         `json:"time_lock"`
	Timestamp            uint64         `json:"timestamp"`
}

type claimHTLCReq struct {
	BaseTx rest.BaseReq   `json:"base_tx"`
	Sender sdk.AccAddress `json:"sender"`
	Secret string         `json:"secret"`
}

type RefundHTLCReq struct {
	BaseTx rest.BaseReq   `json:"base_tx"`
	Sender sdk.AccAddress `json:"sender"`
}
