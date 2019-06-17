package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get the gateway from a moniker
	r.HandleFunc(
		"/asset/gateways/{moniker}",
		monikerGatewayHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all gateways with an optional owner
	r.HandleFunc(
		"/asset/gateways",
		gatewaysHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get gateway creation fee
	r.HandleFunc(
		"/asset/fees/gateways/{moniker}",
		gatewayFeeHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get token fees
	r.HandleFunc(
		"/asset/fees/tokens/{id}",
		tokenFeesHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// monikerGatewayHandlerFn is the HTTP request handler to query a gateway of the given moniker
func monikerGatewayHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryGateway(cliCtx, cdc, "custom/asset/gateway")
}

// gatewaysHandlerFn is the HTTP request handler to query a set of gateways
func gatewaysHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryGateways(cliCtx, cdc, "custom/asset/gateways")
}

// gatewayFeeHandlerFn is the HTTP request handler to query gateway creation fee
func gatewayFeeHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryGatewayFee(cliCtx, cdc, "custom/asset/fees/gateways")
}

// tokenFeesHandlerFn is the HTTP request handler to query token fees
func tokenFeesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryTokenFees(cliCtx, cdc, "custom/asset/fees/tokens")
}
