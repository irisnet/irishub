package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/auth"
)

type sendBody struct {
	Amount string       `json:"amount"`
	Sender string       `json:"sender"`
	BaseTx utils.BaseTx `json:"base_tx"`
}

type burnBody struct {
	Amount string       `json:"amount"`
	Owner  string       `json:"owner"`
	BaseTx utils.BaseTx `json:"base_tx"`
}

// SendRequestHandlerFn - http request handler to send coins to a address
// nolint: gocyclo
func SendRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Init context and read request parameters
		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		to, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var m sendBody
		err = utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}
		baseReq := m.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}
		// Build message
		amount, err := cliCtx.ParseCoins(m.Amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		sender, err := sdk.AccAddressFromBech32(m.Sender)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode delegator. Error: %s", err.Error())))
			return
		}
		msg := bank.BuildBankSendMsg(sender, to, amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Broadcast or return unsigned transaction
		utils.SendOrReturnUnsignedTx(w, cliCtx, m.BaseTx, []sdk.Msg{msg})
	}
}

// BurnRequestHandlerFn - http request handler to burn coins
// nolint: gocyclo
func BurnRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Init context and read request parameters
		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		var m burnBody
		err := utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}
		baseReq := m.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}
		// Build message
		amount, err := cliCtx.ParseCoins(m.Amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		owner, err := sdk.AccAddressFromBech32(m.Owner)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode delegator. Error: %s", err.Error())))
			return
		}
		msg := bank.BuildBankBurnMsg(owner, amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Broadcast or return unsigned transaction
		utils.SendOrReturnUnsignedTx(w, cliCtx, m.BaseTx, []sdk.Msg{msg})
	}
}

type broadcastBody struct {
	Tx auth.StdTx `json:"tx"`
}

// BroadcastTxRequestHandlerFn returns the broadcast tx REST handler
func BroadcastTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		var m broadcastBody
		if err := utils.ReadPostBody(w, r, cliCtx.Codec, &m); err != nil {
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(m.Tx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if cliCtx.DryRun {
			rawRes, err := cliCtx.Query("/app/simulate", txBytes)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			var simulationResult sdk.Result
			if err := cdc.UnmarshalBinaryLengthPrefixed(rawRes, &simulationResult); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			utils.WriteSimulationResponse(w, cliCtx, simulationResult.GasUsed, simulationResult)
			return
		}
		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
