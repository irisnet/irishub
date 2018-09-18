package tx

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/txs", SearchTxRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs/broadcast", BroadcastTxRequestHandlerFn(cliCtx, cdc)).Methods("POST")
}
