package lcd

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/distribution/types"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	distributionclient "github.com/irisnet/irishub/client/distribution"
	"github.com/irisnet/irishub/client/utils"
)

// QueryWithdrawAddressHandlerFn performs withdraw address query
func QueryWithdrawAddressHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["delegatorAddr"]

		delAddr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		key := distribution.GetDelegatorWithdrawAddrKey(delAddr)

		res, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}
		withdrawAddress := sdk.AccAddress(res)

		w.Write([]byte(withdrawAddress.String()))
	}
}

// QueryDelegatorDistInfoHandlerFn query all delegation distribution info of the specified delegator
func QueryDelegatorDistInfoHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["delegatorAddr"]

		delAddr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		key := distribution.GetDelegationDistInfosKey(delAddr)
		resKVs, err := cliCtx.QuerySubspace(key, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var ddiList []types.DelegationDistInfo
		for _, kv := range resKVs {
			var ddi types.DelegationDistInfo
			err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(kv.Value, &ddi)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			ddiList = append(ddiList, ddi)
		}
		utils.PostProcessResponse(w, cliCtx.Codec, ddiList, cliCtx.Indent)
	}
}

// QueryDelegationDistInfoHandlerFn query delegation distribution info
func QueryDelegationDistInfoHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
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

		key := distribution.GetDelegationDistInfoKey(delAddr, valAddr)
		res, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}

		var ddi types.DelegationDistInfo
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &ddi)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cliCtx.Codec, ddi, cliCtx.Indent)
	}
}

// QueryValidatorDistInfoHandlerFn query validator distribution info
func QueryValidatorDistInfoHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		validatorAddrStr := vars["validatorAddr"]
		valAddr, err := sdk.ValAddressFromBech32(validatorAddrStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		key := distribution.GetValidatorDistInfoKey(valAddr)

		res, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}

		var vdi types.ValidatorDistInfo
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &vdi)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		vdiOutput := distributionclient.ConvertToValidatorDistInfoOutput(cliCtx, vdi)

		utils.PostProcessResponse(w, cliCtx.Codec, vdiOutput, cliCtx.Indent)
	}
}
