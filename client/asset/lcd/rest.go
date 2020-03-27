package lcd

import (
	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

const (
	RestParamSymbol = "symbol"
	RestParamOwner  = "owner"
)

// RegisterRoutes registers asset-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc)
}
