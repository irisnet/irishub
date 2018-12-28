package lcd

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/auth/accounts/{address}",
		QueryAccountRequestHandlerFn("acc", cdc, utils.GetAccountDecoder(cdc), cliCtx)).Methods("GET")

	r.HandleFunc("/bank/balances/{address}",
		QueryBalancesRequestHandlerFn("acc", cdc, utils.GetAccountDecoder(cdc), cliCtx)).Methods("GET")
	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/bank/coin/{coin-type}",
		QueryCoinTypeRequestHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc("/tx/sign", SignTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/tx/broadcast", BroadcastTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc("/txs/send", SendTxRequestHandlerFn(cliCtx, cdc)).Methods("POST")
}
