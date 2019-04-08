package lcd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/stake"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

// query accountREST Handler
func QueryBalancesRequestHandlerFn(cdc *codec.Codec, decoder auth.AccountDecoder,
	cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		cliCtx = cliCtx.WithAccountDecoder(decoder)

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := cliCtx.EnsureAccountExistsFromAddr(addr); err != nil {
			utils.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		acc, err := cliCtx.GetAccount(addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		utils.PostProcessResponse(w, cdc, acc.GetCoins(), cliCtx.Indent)
	}
}

// QueryAccountRequestHandlerFn performs account information query
func QueryAccountRequestHandlerFn(cdc *codec.Codec, decoder auth.AccountDecoder,
	cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		cliCtx = cliCtx.WithAccountDecoder(decoder)

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := cliCtx.EnsureAccountExistsFromAddr(addr); err != nil {
			utils.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		acc, err := cliCtx.GetAccount(addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		accountRes, err := bank.ConvertAccountCoin(cliCtx, acc)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, accountRes, cliCtx.Indent)
	}
}

// QueryCoinTypeRequestHandlerFn performs coin type query
func QueryCoinTypeRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coinType := vars["coin-type"]
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
func QueryTokenStatsRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resToken, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, auth.QueryTokenStats), nil)
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

		tokenStats.BondedToken = poolStatus.BondedTokens

		utils.PostProcessResponse(w, cdc, tokenStats, cliCtx.Indent)
	}
}
