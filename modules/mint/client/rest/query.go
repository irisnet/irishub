package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// get the current mint parameter values
	//r.HandleFunc(fmt.Sprintf("/%s/params", types.ModuleName), queryParamsHandlerFn(cliCtx)).Methods("GET")
}

// // HTTP request handler to get the current mint parameter values
// func queryParamsHandlerFn(cliCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)

// 		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
// 		if !ok {
// 			return
// 		}

// 		res, height, err := cliCtx.QueryWithData(route, nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		cliCtx = cliCtx.WithHeight(height)
// 		rest.PostProcessResponse(w, cliCtx, res)
// 	}
// }
