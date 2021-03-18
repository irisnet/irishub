package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/htlc/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// create an HTLC
	r.HandleFunc("/htlc/htlcs", createHTLCHandlerFn(cliCtx)).Methods("POST")
	// claim an HTLC
	r.HandleFunc(fmt.Sprintf("/htlc/htlcs/{%s}/claim", RestID), claimHTLCHandlerFn(cliCtx)).Methods("POST")
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

		if _, err := hex.DecodeString(req.HashLock); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgCreateHTLC(
			req.Sender, req.To, req.ReceiverOnOtherChain, req.SenderOnOtherChain,
			req.Amount, req.HashLock, req.Timestamp, req.TimeLock, req.Transfer,
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

		if _, err := hex.DecodeString(vars[RestID]); err != nil {
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

		if _, err := hex.DecodeString(req.Secret); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgClaimHTLC(req.Sender, vars[RestID], req.Secret)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, &msg)
	}
}
