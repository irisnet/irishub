package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// RegisterRoutes registers asset-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	token.RegisterRESTRoutes(cliCtx, r, queryRoute)
}
