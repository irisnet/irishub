package main

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"encoding/json"
	"github.com/tendermint/go-crypto/keys"
	"github.com/cosmos/cosmos-sdk/client/context"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/pkg/errors"
)


func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/tx/send", SendTxRequestHandlerFn(cdc, kb, ctx)).Methods("POST")
}


type sendTx struct {
	Msg        string    	   `json:"msg"`
	MsgType    string          `json:"type"`
	Fee        auth.StdFee     `json:"fee"`
	Signatures []StdSignature  `json:"signatures"`
}

type StdSignature struct {
	PubKey    		 crypto.PubKeyEd25519	`json:"pub_key"` // optional
	Signature 		 crypto.SignatureEd25519	`json:"signature"`
	AccountNumber    int64 		`json:"account_number"`
	Sequence         int64 		`json:"sequence"`
}

//send traction(sign with rainbow) to irishub
func SendTxRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext) http.HandlerFunc {
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
		var sig = make([]auth.StdSignature,len(tx.Signatures))
		for index,s := range tx.Signatures {
			sig[index].PubKey = s.PubKey
			sig[index].Signature = s.Signature
			sig[index].AccountNumber =s.AccountNumber
			sig[index].Sequence = s.Sequence
		}

		msg,err := convertMsg(tx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var stdTx = auth.StdTx{
			Msg:msg,
			Fee:tx.Fee,
			Signatures:sig,
		}
		txByte,_ := cdc.MarshalBinary(stdTx)
		// send
		res, err := ctx.BroadcastTx(txByte)
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

func convertMsg(tx sendTx) (sdk.Msg,error){
	data := []byte(tx.Msg)
	switch tx.MsgType {
	case "transfer":{
		var msg bank.MsgSend
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil,err
		}
		return msg,nil
	}
	case "delegate":
		var msg stake.MsgDelegate
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil,err
		}
		return msg,nil
	case "unbond":
		var msg stake.MsgUnbond
		if err := json.Unmarshal(data, &msg); err != nil {
			return nil,err
		}
		return msg,nil
	}

	return nil,errors.New("invalid message type")
}