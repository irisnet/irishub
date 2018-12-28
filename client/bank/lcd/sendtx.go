package lcd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

type sendBody struct {
	Amount string       `json:"amount"`
	Sender string       `json:"sender"`
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
		msg := bank.BuildMsg(sender, to, amount)
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

type sendTx struct {
	Msgs       []string       `json:"msgs"`
	Fee        auth.StdFee    `json:"fee"`
	Signatures []stdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

type stdSignature struct {
	PubKey        []byte `json:"pub_key"` // optional
	Signature     []byte `json:"signature"`
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
}

func SendTxRequestHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendTxBody sendTx
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = json.Unmarshal(body, &sendTxBody)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx = utils.InitReqCliCtx(cliCtx, r)

		var sig = make([]auth.StdSignature, len(sendTxBody.Signatures))
		for index, s := range sendTxBody.Signatures {
			var pubkey crypto.PubKey
			if err := cdc.UnmarshalBinaryBare(s.PubKey, &pubkey); err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			sig[index].PubKey = pubkey
			sig[index].Signature = s.Signature
			sig[index].AccountNumber = s.AccountNumber
			sig[index].Sequence = s.Sequence
		}

		var msgs = make([]sdk.Msg, len(sendTxBody.Msgs))
		for index, msgS := range sendTxBody.Msgs {
			var data = []byte(msgS)
			var msg sdk.Msg
			if err := cdc.UnmarshalJSON(data, &msg); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			msgs[index] = msg
		}

		var stdTx = auth.StdTx{
			Msgs:       msgs,
			Fee:        sendTxBody.Fee,
			Signatures: sig,
			Memo:       sendTxBody.Memo,
		}
		txBytes, err := cdc.MarshalBinaryLengthPrefixed(stdTx)
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
