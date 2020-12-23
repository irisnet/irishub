package tx

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"

	"github.com/irisnet/irishub/modules/legacy/types"
)

// RegisterTxRoutes registers all transaction routes on the provided router.
func RegisterTxRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc(fmt.Sprintf("/%s/txs/{hash}", types.ModuleName), QueryTxRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/txs", types.ModuleName), QueryTxsRequestHandlerFn(clientCtx)).Methods("GET")
}
