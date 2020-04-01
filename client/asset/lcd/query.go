package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Query token by symbol
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}", RestParamSymbol),
		queryTokenHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query tokens by owner
	r.HandleFunc(
		"/asset/tokens",
		queryTokensHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query token fees
	r.HandleFunc(
		fmt.Sprintf("/asset/tokens/{%s}/fee", RestParamSymbol),
		queryTokenFeesHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// queryTokenHandlerFn is the HTTP request handler to query token
func queryTokenHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]

		if err := asset.CheckSymbol(symbol); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := asset.QueryTokenParams{
			Symbol: symbol,
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

// queryTokensHandlerFn is the HTTP request handler to query tokens
func queryTokensHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerStr := r.FormValue(RestParamOwner)

		var err error
		var owner sdk.AccAddress

		if len(ownerStr) > 0 {
			owner, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

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

// queryTokenFeesHandlerFn is the HTTP request handler to query token fees
func queryTokenFeesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]

		if err := asset.CheckSymbol(symbol); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

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
