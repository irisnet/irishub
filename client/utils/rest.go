package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client/context"
	"net/http"
	"net/url"
	"github.com/cosmos/cosmos-sdk/wire"
	"io/ioutil"
)

const (
	Async        = "async"
	GenerateOnly = "generate-only"
)

// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func ReadPostBody(w http.ResponseWriter, r *http.Request, cdc *wire.Codec, req interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return err
	}
	err = cdc.UnmarshalJSON(body, req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}

func InitRequestClictx(cliCtx context.CLIContext, r *http.Request, name string, signerAddress string) context.CLIContext {
	cliCtx.GenerateOnly = GenerateOnlyArg(r)
	cliCtx.Async = AsyncOnlyArg(r)
	cliCtx.FromAddressName = name
	cliCtx.Signer = signerAddress
	return cliCtx
}

func SendOrReturnUnsignedTx(w http.ResponseWriter, cliCtx context.CLIContext, txCtx context.TxContext, baseTx context.BaseTx, msgs []sdk.Msg) {

	if cliCtx.GenerateOnly {
		WriteGenerateStdTxResponse(w, txCtx, msgs)
		return
	}

	txBytes, err := txCtx.BuildAndSign(baseTx.LocalAccountName, baseTx.Password, msgs)
	if err != nil {
		WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	var res interface{}
	if cliCtx.Async {
		res, err = cliCtx.BroadcastTxAsync(txBytes)
	} else {
		res, err = cliCtx.BroadcastTx(txBytes)
	}

	output, err := txCtx.Codec.MarshalJSONIndent(res, "", "  ")
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(output)
}

// WriteGenerateStdTxResponse writes response for the generate_only mode.
func WriteGenerateStdTxResponse(w http.ResponseWriter, txCtx context.TxContext, msgs []sdk.Msg) {
	stdMsg, err := txCtx.Build(msgs)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	output, err := txCtx.Codec.MarshalJSON(auth.NewStdTx(stdMsg.Msgs, stdMsg.Fee, nil, stdMsg.Memo))
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Write(output)
	return
}

func AsyncOnlyArg(r *http.Request) bool { return urlQueryHasArg(r.URL, Async) }

func GenerateOnlyArg(r *http.Request) bool { return urlQueryHasArg(r.URL, GenerateOnly) }

func urlQueryHasArg(url *url.URL, arg string) bool { return url.Query().Get(arg) == "true" }
