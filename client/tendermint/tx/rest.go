package tx

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/txs", SearchTxRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs/send", SendTxRequestHandlerFn(cliCtx, cdc)).Methods("POST")
}
