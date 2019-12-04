package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/modules/rand/client/types"
	"github.com/irisnet/irishub/modules/rand"
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

		params := rand.QueryRandParams{
			ReqID: reqID,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", rand.QuerierRoute, rand.QueryRand), bz)
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

		rest.PostProcessResponse(w, cliCtx, readableRand)
	}
}

// queryQueueHandlerFn performs rand request queue query by an optional heigth
func queryQueueHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		heightStr := r.FormValue("height")

		var (
			height int64
			err    error
		)

		if len(heightStr) != 0 {
			height, ok := rest.ParseInt64OrReturnBadRequest(w, heightStr)
			if !ok {
				return
			}

			if height < 0 {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "the height must not be less than 0")
				return
			}
		}

		params := rand.QueryRandRequestQueueParams{
			Height: height,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", rand.QuerierRoute, rand.QueryRandRequestQueue), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
