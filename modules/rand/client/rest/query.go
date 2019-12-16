package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	clienttypes "github.com/irisnet/irishub/modules/rand/client/types"
	"github.com/irisnet/irishub/modules/rand/internal/types"
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
		if err := types.CheckReqID(reqID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryRandParams{
			ReqID: reqID,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRand)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var rawRand types.Rand
		err = cliCtx.Codec.UnmarshalJSON(res, &rawRand)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		readableRand := clienttypes.ReadableRand{
			RequestTxHash: hex.EncodeToString(rawRand.RequestTxHash),
			Height:        rawRand.Height,
			Value:         rawRand.Value.FloatString(types.RandPrec),
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

		params := types.QueryRandRequestQueueParams{
			Height: genHeight,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandRequestQueue)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
