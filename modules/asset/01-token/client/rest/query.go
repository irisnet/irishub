package rest

import (
	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Get token by symbol
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}", RestParamSymbol),
		queryTokenHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
	// Search tokens
	r.HandleFunc(
		"/asset/tokens",
		queryTokensHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get token fees
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}/fees", RestParamSymbol),
		tokenFeesHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get the current asset parameter values
	r.HandleFunc(
		"/asset/tokens/parameters",
		paramsHandlerFn(cliCtx, queryRoute),
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

// tokenFeesHandlerFn is the HTTP request handler to query token fees
func tokenFeesHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryTokenFees(cliCtx, queryRoute)
}

// HTTP request handler to query the staking params values
func paramsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/parameters", queryRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
