package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// create an HTLC
	r.HandleFunc("/htlc/htlcs", createHTLCHandlerFn(cliCtx)).Methods("POST")
	// claim an HTLC
	r.HandleFunc(fmt.Sprintf("/htlc/htlcs/{%s}/claim", RestHashLock), claimHTLCHandlerFn(cliCtx)).Methods("POST")
	// refund an HTLC
	r.HandleFunc(fmt.Sprintf("/htlc/htlcs/{%s}/refund", RestHashLock), refundHTLCHandlerFn(cliCtx)).Methods("POST")
}

// HTTP request handler to create HTLC.
func createHTLCHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateHTLCReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		hashLock, err := hex.DecodeString(req.HashLock)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the NewMsgCreateHTLC message
		msg := types.NewMsgCreateHTLC(
			req.Sender,
			req.To,
			req.ReceiverOnOtherChain,
			req.Amount,
			hashLock,
			req.Timestamp,
			req.TimeLock,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

// HTTP request handler to claim HTLC.
func claimHTLCHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hashLockStr := vars[RestHashLock]
		hashLock, err := hex.DecodeString(hashLockStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req ClaimHTLCReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the NewMsgClaimHTLC message
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

// HTTP request handler to refund HTLC.
func refundHTLCHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hashLockStr := vars[RestHashLock]
		hashLock, err := hex.DecodeString(hashLockStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RefundHTLCReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the NewMsgRefundHTLC message
		msg := types.NewMsgRefundHTLC(req.Sender, hashLock)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
