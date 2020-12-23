package bank

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"

	"github.com/irisnet/irishub/modules/legacy/types"
)

// RegisterRoutes registers all transaction routes on the provided router.
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc(fmt.Sprintf("/%s/bank/accounts/{address}", types.ModuleName), QueryAccountRequestHandlerFn(clientCtx)).Methods("GET")
}
