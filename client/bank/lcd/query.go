package lcd

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/stake"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/auth"
	stakeTypes "github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// QueryAccountRequestHandlerFn performs account information query
func QueryAccountRequestHandlerFn(cdc *codec.Codec, decoder auth.AccountDecoder,
	cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		cliCtx = cliCtx.WithAccountDecoder(decoder)

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), protocol.AccountStore)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value
		account, err := decoder(res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error()))
			return
		}

		utils.PostProcessResponse(w, cdc, account, cliCtx.Indent)
	}
}

// QueryCoinTypeRequestHandlerFn performs coin type query
func QueryCoinTypeRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coinType := vars["type"]
		res, err := cliCtx.GetCoinType(coinType)
		if err != nil && strings.Contains(err.Error(), "unsupported coin type") {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// QueryTokenStatsRequestHandlerFn performs token statistic query
func QueryTokenStatsRequestHandlerFn(cdc *codec.Codec, decoder auth.AccountDecoder, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		assetId := vars["id"]

		if len(assetId) == 0 || assetId == "iris" || assetId == "iris-atto" {
			resToken, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, bank.QueryTokenStats), nil)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			}

			var tokenStats bank.TokenStats
			err = cdc.UnmarshalJSON(resToken, &tokenStats)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			resPool, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.StakeRoute, stake.QueryPool), nil)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			var poolStatus stakeTypes.PoolStatus
			err = cdc.UnmarshalJSON(resPool, &poolStatus)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			tokenStats.BondedTokens = sdk.Coins{sdk.Coin{Denom: stakeTypes.StakeDenom, Amount: poolStatus.BondedTokens.TruncateInt()}}

			utils.PostProcessResponse(w, cdc, tokenStats, cliCtx.Indent)
		} else {
			cliCtx = cliCtx.WithAccountDecoder(decoder)
			params := asset.QueryAssetParams{
				Asset: assetId,
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryAsset), bz)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			var nAsset asset.Asset
			err = cdc.UnmarshalJSON(res, &nAsset)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			var tokenStats bank.TokenStats

			//get loose token from asset

			looseToken := sdk.Coin{}
			looseToken.Denom = nAsset.GetDenom()
			looseToken.Amount = nAsset.GetTotalSupply()
			tokenStats.LooseTokens = append(tokenStats.LooseTokens, looseToken)

			//get burned token from burnAddress
			burnedAcc, err := cliCtx.GetAccount(bank.BurnedCoinsAccAddr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			burnToken := sdk.Coin{nAsset.GetDenom(), burnedAcc.Coins.AmountOf(nAsset.GetDenom())}
			tokenStats.BurnedTokens = append(tokenStats.BurnedTokens, burnToken)
			utils.PostProcessResponse(w, cdc, tokenStats, cliCtx.Indent)
		}

	}
}
