package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Get token by id
	r.HandleFunc(
		"/asset/tokens/{id}",
		queryTokenHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
	// Search tokens
	r.HandleFunc(
		"/asset/tokens",
		queryTokensHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
	// Get the gateway from a moniker
	r.HandleFunc(
		"/asset/gateways/{moniker}",
		monikerGatewayHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get all gateways with an optional owner
	r.HandleFunc(
		"/asset/gateways",
		gatewaysHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get gateway creation fee
	r.HandleFunc(
		"/asset/fees/gateways/{moniker}",
		gatewayFeeHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get token fees
	r.HandleFunc(
		"/asset/fees/tokens/{id}",
		tokenFeesHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
}

// queryTokenHandlerFn performs token information query
func queryTokenHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryToken(cliCtx, queryRoute)
}

// queryTokenHandlerFn performs token information query
func queryTokensHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryTokens(cliCtx, queryRoute)
}

// monikerGatewayHandlerFn is the HTTP request handler to query a gateway of the given moniker
func monikerGatewayHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryGateway(cliCtx, queryRoute)
}

// gatewaysHandlerFn is the HTTP request handler to query a set of gateways
func gatewaysHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryGateways(cliCtx, queryRoute)
}

// gatewayFeeHandlerFn is the HTTP request handler to query gateway creation fee
func gatewayFeeHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryGatewayFee(cliCtx, queryRoute)
}

// tokenFeesHandlerFn is the HTTP request handler to query token fees
func tokenFeesHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryTokenFees(cliCtx, queryRoute)
}
