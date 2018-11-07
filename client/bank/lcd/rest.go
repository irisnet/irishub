package lcd

import (
	"github.com/cosmos/cosmos-sdk/codec"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/bank/{address}/send", SendRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/bank/accounts/{address}",
		QueryAccountRequestHandlerFn("acc", cdc, authcmd.GetAccountDecoder(cdc), cliCtx)).Methods("GET")
	r.HandleFunc("/bank/coin/{coin-type}",
		QueryCoinTypeRequestHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc("/tx/sign", SignTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/tx/broadcast", BroadcastTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc("/txs/send", SendTxRequestHandlerFn(cliCtx, cdc)).Methods("POST")
}
