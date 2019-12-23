package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// get token by id
	r.HandleFunc(fmt.Sprintf("/asset/tokens/{%s}", RestTokenID), queryTokenHandlerFn(cliCtx, queryRoute)).Methods("GET")
	// search tokens
	r.HandleFunc("/asset/tokens", queryTokensHandlerFn(cliCtx, queryRoute)).Methods("GET")
	// get token fees
	r.HandleFunc(fmt.Sprintf("/asset/fees/tokens/{%s}", RestTokenID), tokenFeesHandlerFn(cliCtx, queryRoute)).Methods("GET")
	// get the current asset parameter values
	r.HandleFunc("/asset/parameters", paramsHandlerFn(cliCtx, queryRoute)).Methods("GET")
}

// HTTP request handler to query token information.
func queryTokenHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryToken(cliCtx, queryRoute)
}

// HTTP request handler to query token information.
func queryTokensHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryTokens(cliCtx, queryRoute)
}

// HTTP request handler to query token fees.
func tokenFeesHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return queryTokenFees(cliCtx, queryRoute)
}

// HTTP request handler to query the staking params values.
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
