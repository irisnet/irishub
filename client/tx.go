package client

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app"
	"github.com/tendermint/tendermint/crypto"
	"io/ioutil"
	"net/http"
)

func RegisterRoutes(ctx app.Context, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/tx/send", SendTxRequestHandlerFn(cdc, kb, ctx)).Methods("POST")
}

type sendTx struct {
	Msgs       []string       `json:"msgs"`
	Fee        auth.StdFee    `json:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

type StdSignature struct {
	PubKey        []byte `json:"pub_key"` // optional
	Signature     []byte `json:"signature"`
	AccountNumber int64  `json:"account_number"`
	Sequence      int64  `json:"sequence"`
}

//send traction(sign with rainbow) to irishub
func SendTxRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx app.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tx sendTx
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if err = json.Unmarshal(body, &tx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		var sig = make([]auth.StdSignature, len(tx.Signatures))
		for index, s := range tx.Signatures {
			var pubkey crypto.PubKey
			if err := cdc.UnmarshalBinaryBare(s.PubKey, &pubkey); err != nil {
				panic(err)
			}
			sig[index].PubKey = pubkey

			var signature crypto.Signature
			if err := cdc.UnmarshalBinaryBare(s.Signature, &signature); err != nil {
				panic(err)
			}
			sig[index].Signature = signature
			sig[index].AccountNumber = s.AccountNumber
			sig[index].Sequence = s.Sequence
		}

		var msgs = make([]sdk.Msg, len(tx.Msgs))
		for index, msgS := range tx.Msgs {
			var data = []byte(msgS)
			var msg sdk.Msg
			if err := cdc.UnmarshalJSON(data, &msg); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			msgs[index] = msg
		}

		var stdTx = auth.StdTx{
			Msgs:       msgs,
			Fee:        tx.Fee,
			Signatures: sig,
			Memo:       tx.Memo,
		}
		txByte, _ := cdc.MarshalBinary(stdTx)
		// send
		res, err := ctx.BroadcastTxSync(txByte)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
