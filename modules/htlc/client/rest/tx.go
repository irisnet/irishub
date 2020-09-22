package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/htlc/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// create an HTLC
	r.HandleFunc("/htlc/htlcs", createHTLCHandlerFn(cliCtx)).Methods("POST")
	// claim an HTLC
	r.HandleFunc(fmt.Sprintf("/htlc/htlcs/{%s}/claim", RestHashLock), claimHTLCHandlerFn(cliCtx)).Methods("POST")
	// refund an HTLC
	r.HandleFunc(fmt.Sprintf("/htlc/htlcs/{%s}/refund", RestHashLock), refundHTLCHandlerFn(cliCtx)).Methods("POST")
}

// CreateHTLCReq defines the properties of an HTLC creation request's body.
type CreateHTLCReq struct {
	BaseReq              rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Sender               sdk.AccAddress `json:"sender" yaml:"sender"`
	To                   sdk.AccAddress `json:"to" yaml:"to"`
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain" yaml:"receiver_on_other_chain"`
	Amount               sdk.Coins      `json:"amount" yaml:"amount"`
	HashLock             string         `json:"hash_lock" yaml:"hash_lock"`
	TimeLock             uint64         `json:"time_lock" yaml:"time_lock"`
	Timestamp            uint64         `json:"timestamp" yaml:"timestamp"`
}

// ClaimHTLCReq defines the properties of an HTLC claim request's body.
type ClaimHTLCReq struct {
	BaseReq rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Sender  sdk.AccAddress `json:"sender" yaml:"sender"`
	Secret  string         `json:"secret" yaml:"secret"`
}

// RefundHTLCReq defines the properties of an HTLC refund request's body.
type RefundHTLCReq struct {
	BaseReq rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Sender  sdk.AccAddress `json:"sender" yaml:"sender"`
}

func createHTLCHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateHTLCReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		hashLock, err := hex.DecodeString(req.HashLock)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateHTLC(
			req.Sender, req.To, req.ReceiverOnOtherChain,
			req.Amount, hashLock, req.Timestamp, req.TimeLock,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func claimHTLCHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hashLock, err := hex.DecodeString(vars[RestHashLock])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req ClaimHTLCReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		secret, err := hex.DecodeString(req.Secret)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgClaimHTLC(req.Sender, hashLock, secret)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}

func refundHTLCHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hashLock, err := hex.DecodeString(vars[RestHashLock])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RefundHTLCReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgRefundHTLC(req.Sender, hashLock)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}
