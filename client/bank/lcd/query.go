package lcd

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/gorilla/mux"
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
func QueryBalancesRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
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
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, account.GetCoins(), cliCtx.Indent)
	}
}

// QueryAccountRequestHandlerFn performs account information query
func QueryAccountRequestHandlerFn(storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), storeName)
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

		accountRes, err := bank.ConvertAccountCoin(cliCtx, account)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, accountRes, cliCtx.Indent)
	}
}

// QueryCoinTypeRequestHandlerFn performs coin type query
func QueryCoinTypeRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext,
) http.HandlerFunc {
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

// QueryCoinTypeRequestHandlerFn performs coin type query
func QueryTokenStatsRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, accStore, stakeStore string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query acc store
		var loosenToken sdk.Coins
		var burnedToken sdk.Coins
		res, err := cliCtx.QueryStore(auth.TotalLoosenTokenKey, accStore)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if res == nil {
			loosenToken = nil
		} else {
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &loosenToken)
		}
		res, err = cliCtx.QueryStore(auth.BurnTokenKey, accStore)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if res == nil {
			burnedToken = nil
		} else {
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &burnedToken)
		}

		// Query stake store
		var bondedPool stake.BondedPool
		res, err = cliCtx.QueryStore(stake.PoolKey, stakeStore)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if res != nil {
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &bondedPool)
		}
		if !bondedPool.BondedTokens.Equal(bondedPool.BondedTokens.TruncateDec()) {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "get invalid bonded token amount")
			return
		}
		bondedToken := sdk.NewCoin(stakeTypes.StakeDenom, bondedPool.BondedTokens.TruncateInt())

		//Convert to main coin unit
		loosenTokenStr, err := cliCtx.ConvertCoinToMainUnit(loosenToken.String())
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		burnedTokenStr, err := cliCtx.ConvertCoinToMainUnit(burnedToken.String())
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		bondedTokenStr, err := cliCtx.ConvertCoinToMainUnit(bondedToken.String())
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		tokenStats := bank.TokenStats{
			LoosenToken: loosenTokenStr,
			BurnedToken: burnedTokenStr,
			BondedToken: bondedTokenStr[0],
		}

		utils.PostProcessResponse(w, cdc, tokenStats, cliCtx.Indent)
	}
}
