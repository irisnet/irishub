package lcd

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v2/htlc"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// create an HTLC
	r.HandleFunc(
		"/htlc/htlcs",
		createHTLCHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// claim an HTLC
	r.HandleFunc(
		"/htlc/htlcs/{hash-lock}/claim",
		claimHTLCHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// refund an HTLC
	r.HandleFunc(
		"/htlc/htlcs/{hash-lock}/refund",
		refundHTLCHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type createHTLCReq struct {
	BaseTx               utils.BaseTx   `json:"base_tx"`
	Sender               sdk.AccAddress `json:"sender"`
	To                   sdk.AccAddress `json:"to"`
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain"`
	Amount               sdk.Coins      `json:"amount"`
	HashLock             string         `json:"hash_lock"`
	TimeLock             uint64         `json:"time_lock"`
	Timestamp            uint64         `json:"timestamp"`
}

func createHTLCHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createHTLCReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		hashLock, err := hex.DecodeString(req.HashLock)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the NewMsgCreateHTLC message
		msg := htlc.NewMsgCreateHTLC(
			req.Sender, req.To, req.ReceiverOnOtherChain, req.Amount,
			hashLock, req.Timestamp, req.TimeLock)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

type claimHTLCReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Sender sdk.AccAddress `json:"sender"`
	Secret string         `json:"secret"`
}

func claimHTLCHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hashLockStr := vars["hash-lock"]
		hashLock, err := hex.DecodeString(hashLockStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req claimHTLCReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the NewMsgClaimHTLC message
		secret, err := hex.DecodeString(req.Secret)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := htlc.NewMsgClaimHTLC(
			req.Sender, hashLock, secret)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

type RefundHTLCReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Sender sdk.AccAddress `json:"sender"`
}

func refundHTLCHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hashLockStr := vars["hash-lock"]
		hashLock, err := hex.DecodeString(hashLockStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RefundHTLCReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the NewMsgRefundHTLC message
		msg := htlc.NewMsgRefundHTLC(
			req.Sender, hashLock)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
