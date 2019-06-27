package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func queryToken(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tokenId := vars["id"]

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

func queryTokens(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		source := r.FormValue("source")
		gateway := r.FormValue("gateway")
		owner := r.FormValue("owner")

		// TODO: pagination support

		params := asset.QueryTokensParams{
			Source:  source,
			Gateway: gateway,
			Owner:   owner,
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

// queryGateway queries a gateway of the given moniker from the specified endpoint
func queryGateway(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		moniker := vars["moniker"]
		if len(moniker) < asset.MinimumGatewayMonikerSize || len(moniker) > asset.MaximumGatewayMonikerSize {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("the length of the moniker must be [%d,%d]", asset.MinimumGatewayMonikerSize, asset.MaximumGatewayMonikerSize))
			return
		}

		if !asset.IsAlpha(moniker) {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("the moniker must contain only letters"))
			return
		}

		params := asset.QueryGatewayParams{
			Moniker: moniker,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryGateway), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// queryGateways queries all gateways with an optional owner from the specified endpoint
func queryGateways(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerStr := r.FormValue("owner")

		var (
			owner sdk.AccAddress
			err   error
		)

		if ownerStr != "" {
			owner, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := asset.QueryGatewaysParams{
			Owner: owner,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryGateways), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// queryGatewayFee queries the gateway creation fee from the specified endpoint
func queryGatewayFee(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		moniker := vars["moniker"]
		if err := asset.ValidateMoniker(moniker); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := asset.QueryGatewayFeeParams{
			Moniker: moniker,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s/gateways", protocol.AssetRoute, asset.QueryFees), bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// queryTokenFees queries the token related fees from the specified endpoint
func queryTokenFees(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id := vars["id"]
		if err := asset.CheckAssetID(id); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := asset.QueryTokenFeesParams{
			ID: id,
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
