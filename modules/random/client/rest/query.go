package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/modules/random/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// query rand by the request id
	r.HandleFunc(fmt.Sprintf("/rand/rands/{%s}", RestRequestID), queryRandomHandlerFn(cliCtx)).Methods("GET")
	// query rand request queue by an optional heigth
	r.HandleFunc("/rand/queue", queryQueueHandlerFn(cliCtx)).Methods("GET")
}

// HTTP request handler to query rand by the request id.
func queryRandomHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		reqID := vars[RestRequestID]
		if err := types.CheckReqID(reqID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryRandomParams{
			ReqID: reqID,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandom)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var rawRandom types.Random
		if err := cliCtx.Codec.UnmarshalJSON(res, &rawRandom); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		readableRandom := types.ReadableRandom{
			RequestTxHash: hex.EncodeToString(rawRandom.RequestTxHash),
			Height:        rawRandom.Height,
			Value:         rawRandom.Value,
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, readableRandom)
	}
}

// HTTP request handler to query request queue by an optional heigth.
func queryQueueHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genHeightStr := r.FormValue("gen-height")

		genHeight, err := strconv.ParseInt(genHeightStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		params := types.QueryRandomRequestQueueParams{
			Height: genHeight,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandomRequestQueue)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
