package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Query liquidity
	r.HandleFunc(
		"/swap/liquidities",
		queryLiquidityHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// queryLiquidityHandlerFn performs liquidity information query
func queryLiquidityHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryLiquidity(cliCtx, cdc, "custom/swap/liquidities")
}