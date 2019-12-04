package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/modules/rand"
	"github.com/irisnet/irishub/modules/rand/client/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/rand/rands/{%s}", RestRequestID), queryRandHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/rand/queue", queryQueueHandlerFn(cliCtx)).Methods("GET")
}

// queryRandHandlerFn performs rand query by the request id
func queryRandHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		reqID := vars["request-id"]
		if err := rand.CheckReqID(reqID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := rand.QueryRandParams{
			ReqID: reqID,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", rand.QuerierRoute, rand.QueryRand)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var rawRand rand.Rand
		err = cliCtx.Codec.UnmarshalJSON(res, &rawRand)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		readableRand := types.ReadableRand{
			RequestTxHash: hex.EncodeToString(rawRand.RequestTxHash),
			Height:        rawRand.Height,
			Value:         rawRand.Value.FloatString(rand.RandPrec),
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, readableRand)
	}
}

// queryQueueHandlerFn performs rand request queue query by an optional heigth
func queryQueueHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genHeightStr := r.FormValue("gen-height")

		genHeight, ok := rest.ParseInt64OrReturnBadRequest(w, genHeightStr)
		if !ok {
			return
		}

		if genHeight < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "the generation height must not be less than 0")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := rand.QueryRandRequestQueueParams{
			Height: genHeight,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", rand.QuerierRoute, rand.QueryRandRequestQueue)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
