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
		createHtlcHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// claim an HTLC
	r.HandleFunc(
		"/htlc/htlcs/{hash-lock}/claim",
		claimHtlcHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// refund an HTLC
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
	Amount               sdk.Coin       `json:"amount"`
	HashLock             string         `json:"hash_lock"`
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

		hashLock, err := hex.DecodeString(req.HashLock)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the NewMsgCreateHTLC message
		msg := htlc.NewMsgCreateHTLC(
			req.Sender, req.Receiver, receiverOnOtherChain, req.Amount,
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

type claimHtlcReq struct {
	BaseTx utils.BaseTx   `json:"base_tx"`
	Sender sdk.AccAddress `json:"sender"`
	Secret string         `json:"secret"`
}

func claimHtlcHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		hashLockStr := vars["hash-lock"]
		hashLock, err := hex.DecodeString(hashLockStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req claimHtlcReq
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
			req.Sender, secret, hashLock)
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

		hashLockStr := vars["hash-lock"]
		hashLock, err := hex.DecodeString(hashLockStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RefundHtlcReq
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
