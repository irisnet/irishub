package rest

import (
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
	"net/http"
	sdk "github.com/cosmos/cosmos-sdk"
	"encoding/hex"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/pkg/errors"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/commands/query"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/go-wire"
	"io"
)


type ServiceByteTx struct {
	manager keys.Manager
}

type RequestTx struct {
	Tx sdk.Tx `json:"tx" validate:"required"`
}

func NewServiceByteTx(manager keys.Manager) *ServiceByteTx {
	return &ServiceByteTx{
		manager: manager, // XXX keycmd.GetKeyManager()
	}
}

func (s *ServiceByteTx) RegisterByteTx(r *mux.Router) error {
	r.HandleFunc("/byteTx", s.ByteTx).Methods("POST")
	return nil
}


func (s *ServiceByteTx) RegisterqueryTx(r *mux.Router) error {
	r.HandleFunc("/tx/{hash}", s.queryTx).Methods("GET")
	return nil
}

func (s *ServiceByteTx) ByteTx(w http.ResponseWriter, r *http.Request) {
	req := new(RequestTx)
	if err := sdk.ParseRequestAndValidateJSON(r, req); err != nil {
		sdk.WriteError(w, err)
		return
	}

	tx := req.Tx

	if sign, ok := tx.Unwrap().(keys.Signable); ok {
		sdk.WriteSuccess(w, hex.EncodeToString(sign.SignBytes()))
		return
	}
	sdk.WriteSuccess(w, "")
}

func (s *ServiceByteTx) queryTx(w http.ResponseWriter, r *http.Request){
	args := mux.Vars(r)
	hash := args["hash"]

	if hash == "" {
		sdk.WriteError(w, errors.Errorf("[%s] argument must be non-empty ", "hash"))
		return
	}
	// with tx, we always just parse key as hex and use to lookup
	hashByte, err := hex.DecodeString(cmn.StripHex(hash))

	// get the proof -> this will be used by all prover commands
	node := commands.GetNode()
	prove := !viper.GetBool(commands.FlagTrustNode)
	res, err := node.Tx(hashByte, prove)
	if err != nil {
		sdk.WriteError(w, err)
		return
	}

	// no checks if we don't get a proof
	if !prove {
		sdk.WriteSuccess(w,showTx(w,res.Height, res.Tx))
		return
	}

	cert, err := commands.GetCertifier()
	if err != nil {
		sdk.WriteError(w, err)
		return
	}

	check, err := client.GetCertifiedCommit(res.Height, node, cert)
	if err != nil {
		sdk.WriteError(w, err)
		return
	}
	err = res.Proof.Validate(check.Header.DataHash)
	if err != nil {
		sdk.WriteError(w, err)
		return
	}

	// note that we return res.Proof.Data, not res.Tx,
	// as res.Proof.Validate only verifies res.Proof.Data
	sdk.WriteSuccess(w,showTx(w,res.Height, res.Proof.Data))
}

// showTx parses anything that was previously registered as sdk.Tx
func showTx(w io.Writer ,h int64, tx types.Tx) error {
	var info sdk.Tx
	err := wire.ReadBinaryBytes(tx, &info)
	if err != nil {
		return err
	}
	return query.FoutputProof(w,info,h)
}