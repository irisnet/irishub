package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/distribution"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	sdk "github.com/irisnet/irishub/types"
)

// QueryWithdrawAddressHandlerFn performs withdraw address query
func QueryWithdrawAddressHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["delegatorAddr"]

		delAddr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := distribution.NewQueryDelegatorParams(delAddr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryWithdrawAddr),
			bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// QueryRewardsHandlerFn query the all the rewards of validator or delegator
func QueryRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		AddrStr := vars["address"]
		accAddress, err := sdk.AccAddressFromBech32(AddrStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := distribution.NewQueryRewardsParams(accAddress)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryRewards),
			bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}
