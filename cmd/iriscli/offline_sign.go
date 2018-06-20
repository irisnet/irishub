package main

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"encoding/json"
	"github.com/tendermint/go-crypto/keys"
	"github.com/cosmos/cosmos-sdk/client/context"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
)


func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/tx/send", SendTxRequestHandlerFn(cdc, kb, ctx)).Methods("POST")
}

// /accounts/{address}
type sendTxReq struct {
	tx []byte
}

//send traction(sign with rainbow) to irishub
func SendTxRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sendReq sendTxReq
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err = json.Unmarshal(body, &sendReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// send
		res, err := ctx.BroadcastTx(sendReq.tx)
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