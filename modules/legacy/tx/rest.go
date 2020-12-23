package tx

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
)

// RegisterTxRoutes registers all transaction routes on the provided router.
func RegisterTxRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc("/v0.16.3/txs/{hash}", QueryTxRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/v0.16.3/txs", QueryTxsRequestHandlerFn(clientCtx)).Methods("GET")
}
