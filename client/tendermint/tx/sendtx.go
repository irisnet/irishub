package tx

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/tendermint/tendermint/crypto"
	"io/ioutil"
	"net/http"
)

type sendTx struct {
	Msgs       []string       `json:"msgs"`
	Fee        auth.StdFee    `json:"fee"`
	Signatures []stdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

type stdSignature struct {
	PubKey        []byte `json:"pub_key"` // optional
	Signature     []byte `json:"signature"`
	AccountNumber int64  `json:"account_number"`
	Sequence      int64  `json:"sequence"`
}

func SendTxRequestHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
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
		cliCtx.Async = utils.AsyncOnlyArg(r)

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
		txBytes, err := cdc.MarshalBinary(stdTx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var res interface{}
		if cliCtx.Async {
			res, err = cliCtx.BroadcastTxAsync(txBytes)
		} else {
			res, err = cliCtx.BroadcastTx(txBytes)
		}

		output, err := cdc.MarshalJSONIndent(res, "", "  ")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
