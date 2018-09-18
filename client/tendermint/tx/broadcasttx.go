package tx

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"io/ioutil"
	"net/http"
)

// Tx Broadcast Body
type BroadcastTxBody struct {
	TxBytes string `json:"tx"`
}

// BroadcastTx REST Handler
func BroadcastTxRequestHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m BroadcastTxBody

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = cdc.UnmarshalJSON(body, &m)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var res interface{}
		if utils.AsyncOnlyArg(r) {
			res, err = cliCtx.BroadcastTxAsync([]byte(m.TxBytes))
		} else {
			res, err = cliCtx.BroadcastTx([]byte(m.TxBytes))
		}

		output, err := cdc.MarshalJSONIndent(res, "", "  ")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
