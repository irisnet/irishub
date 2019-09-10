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
	// create a htlc
	r.HandleFunc(
		"/htlc/htlcs",
		createHtlcHandlerFn(cdc, cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/htlc/htlcs/{hash-lock}/claim",
		claimHtlcHandlerFn(cdc, cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/htlc/htlcs/{hash-lock}/refund",
		refundHtlcHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type createHtlcReq struct {
	BaseTx               utils.BaseTx   `json:"base_tx"`
	Sender               sdk.AccAddress `json:"sender"`
	Receiver             sdk.AccAddress `json:"receiver"`
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain"`
	HashLock             string         `json:"hash_lock"`
	InAmount             uint64         `json:"in_amount"`
	Amount               sdk.Coin       `json:"amount"`
	TimeLock             uint64         `json:"time_lock"`
	Timestamp            uint64         `json:"timestamp"`
}

func createHtlcHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createHtlcReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		receiverOnOtherChain, err := hex.DecodeString(req.ReceiverOnOtherChain)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// create the NewMsgCreateHTLC message
		msg := htlc.NewMsgCreateHTLC(
			req.Sender, req.Receiver, receiverOnOtherChain, req.Amount, uint64(req.InAmount),
			req.HashLock, req.Timestamp, req.TimeLock)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

type claimHtlcReq struct {
	BaseTx  utils.BaseTx   `json:"base_tx"`
	Claimer sdk.AccAddress `json:"claimer"`
	Secret  string         `json:"secret"`
}

func claimHtlcHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashLock := vars["hash-lock"]

		var req claimHtlcReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the NewMsgClaimHTLC message
		msg := htlc.NewMsgClaimHTLC(
			req.Claimer, hashLock, req.Secret)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

type RefundHtlcReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Sender sdk.AccAddress `json:"sender"`
}

func refundHtlcHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashLock := vars["hash-lock"]

		var req RefundHtlcReq
		err := utils.ReadPostBody(w, r, cdc, &req)
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
