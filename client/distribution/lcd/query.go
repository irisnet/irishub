package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/distribution"
	distrClient "github.com/irisnet/irishub/client/distribution"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/app/protocol"
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

// QueryDelegatorDistInfoHandlerFn query all delegation distribution info of the specified delegator
func QueryDelegatorDistInfoHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
			fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryAllDelegationDistInfo),
			bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// QueryDelegationDistInfoHandlerFn query delegation distribution info
func QueryDelegationDistInfoHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		delegatorAddrStr := vars["delegatorAddr"]
		delAddr, err := sdk.AccAddressFromBech32(delegatorAddrStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		validatorAddrStr := vars["validatorAddr"]
		valAddr, err := sdk.ValAddressFromBech32(validatorAddrStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := distribution.NewQueryDelegationDistInfoParams(delAddr, valAddr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryDelegationDistInfo),
			bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

// QueryValidatorDistInfoHandlerFn query validator distribution info
func QueryValidatorDistInfoHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		validatorAddrStr := vars["validatorAddr"]
		valAddr, err := sdk.ValAddressFromBech32(validatorAddrStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := distribution.NewQueryValidatorDistInfoParams(valAddr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryValidatorDistInfo),
			bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var vddInfo distribution.ValidatorDistInfo
		if err = cliCtx.Codec.UnmarshalJSON(res, &vddInfo); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		output := distrClient.ConvertToValidatorDistInfoOutput(cliCtx, vddInfo)

		utils.PostProcessResponse(w, cliCtx.Codec, output, cliCtx.Indent)
	}
}

// QueryRewardsHandlerFn query the all the rewards of validator or delegator
func QueryRewardsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		AddrStr := vars["address"]
		accAddress, err := sdk.AccAddressFromBech32(AddrStr)
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
