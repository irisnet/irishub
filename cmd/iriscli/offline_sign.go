package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app"
	"github.com/tendermint/tendermint/crypto"
	"io/ioutil"
	"net/http"
	"fmt"
)

func RegisterRoutes(ctx app.Context, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/tx/send", SendTxRequestHandlerFn(cdc, kb, ctx)).Methods("POST")
}

type sendTx struct {
	Msgs       []string       `json:"msgs"`
	MsgType    string         `json:"type"`
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

type Msgs = []sdk.Msg

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

		msgs, err := convertMsg(tx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var stdTx = auth.StdTx{
			Msgs:       msgs,
			Fee:        tx.Fee,
			Signatures: sig,
			Memo:       tx.Memo,
		}
		txByte, _ := cdc.MarshalBinary(stdTx)
		fmt.Println(txByte)
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

func convertMsg(tx sendTx) (Msgs, error) {
	var msgs Msgs
	for _, msgS := range tx.Msgs {

		switch tx.MsgType {
		case "transfer":
			{
				var msg bank.MsgSend
				var data = []byte(msgS)
				if err := json.Unmarshal(data, &msg); err != nil {
					return nil, err
				}
				msgs = append(msgs, msg)
			}
		case "delegate":
			//var msg stake.MsgDelegate
			//if err := json.Unmarshal(data, &msg); err != nil {
			//	return nil, err
			//}
			//return msg, nil
			//case "unbond":
			//	var msg stake.MsgUnbond
			//	if err := json.Unmarshal(data, &msg); err != nil {
			//		return nil, err
			//	}
			//	return msg, nil
		}
	}
	return msgs, nil
}
