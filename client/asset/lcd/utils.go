package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/asset/tags"
	"github.com/irisnet/irishub/client/context"
	tmtx "github.com/irisnet/irishub/client/tendermint/tx"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
		sourceStr := r.FormValue("source")
		gateway := r.FormValue("gateway")
		owner := r.FormValue("owner")

		// TODO: pagination support
		page := 0
		size := 100

		var source asset.AssetSource
		if len(sourceStr) > 0 {
			_source, ok := asset.StringToAssetSourceMap[sourceStr]
			if !ok {
				utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid source %s", sourceStr))
				return
			}
			source = _source
		} else if len(owner) > 0 {
			source, _ = asset.StringToAssetSourceMap["native"]
		}

		queryTags := []string{fmt.Sprintf("%s='%s'", tags.Action, string(tags.ActionIssueToken))}

		if len(sourceStr) > 0 {
			queryTags = append(queryTags, fmt.Sprintf("%s='%s'", tags.Source, source.String()))
		}

		if len(gateway) > 0 {
			queryTags = append(queryTags, fmt.Sprintf("%s='%s'", tags.Gateway, gateway))
		}

		if len(owner) > 0 {
			queryTags = append(queryTags, fmt.Sprintf("%s='%s'", tags.Owner, owner))
		}

		infos, err := tmtx.SearchTxs(cliCtx, cdc, queryTags, page, size)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var tokens []asset.FungibleToken
		for _, info := range infos {
			if info.Result.Code != abci.CodeTypeOK {
				continue
			}

			for _, msg := range info.Tx.GetMsgs() {
				if msg.Type() == asset.MsgTypeIssueToken {
					msgIssueAsset := msg.(asset.MsgIssueToken)

					var token asset.FungibleToken
					switch msgIssueAsset.Family {
					case asset.FUNGIBLE:
						totalSupply := msgIssueAsset.InitialSupply
						decimal := int(msgIssueAsset.Decimal)
						token = asset.NewFungibleToken(msgIssueAsset.Source, msgIssueAsset.Gateway, msgIssueAsset.Symbol, msgIssueAsset.Name, msgIssueAsset.Decimal, msgIssueAsset.SymbolAtSource, msgIssueAsset.SymbolMinAlias, sdk.NewIntWithDecimal(int64(msgIssueAsset.InitialSupply), decimal), sdk.NewIntWithDecimal(int64(totalSupply), decimal), sdk.NewIntWithDecimal(int64(msgIssueAsset.MaxSupply), decimal), msgIssueAsset.Mintable, msgIssueAsset.Owner)
					default:
						continue
					}

					tokens = append(tokens, token)
				}
			}
		}

		if len(tokens) == 0 {
			tokens = make([]asset.FungibleToken, 0)
		}

		utils.PostProcessResponse(w, cliCtx.Codec, tokens, cliCtx.Indent)
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
		vars := mux.Vars(r)
		ownerStr := vars["owner"]

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
		if ok, err := asset.CheckAssetID(id); !ok {
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
