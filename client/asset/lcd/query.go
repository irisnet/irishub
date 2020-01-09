package lcd

import (
	"fmt"
	"net/http"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/utils"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get token by id
	r.HandleFunc(
		"/asset/tokens/{token-id}",
		queryTokenHandlerFn(cliCtx, cdc),
	).Methods("GET")
	// Search tokens
	r.HandleFunc(
		"/asset/tokens",
		queryTokensHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get token fees
	r.HandleFunc(
		"/asset/tokens/{symbol}/fee",
		tokenFeesHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// queryTokenHandlerFn performs token information query
func queryTokenHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tokenId := vars["token-id"]

		params := asset.QueryTokenParams{
			TokenId: tokenId,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryToken), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// queryTokenHandlerFn performs token information query
func queryTokensHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := r.FormValue("owner")

		// TODO: pagination support

		params := asset.QueryTokensParams{
			Owner: owner,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryTokens), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// tokenFeesHandlerFn is the HTTP request handler to query token fees
func tokenFeesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		symbol := vars["symbol"]

		params := asset.QueryTokenFeesParams{
			Symbol: symbol,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s/tokens", protocol.AssetRoute, asset.QueryFees), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}
