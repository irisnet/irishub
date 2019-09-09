package lcd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v1/stake"
	"github.com/irisnet/irishub/app/v1/stake/tags"
	"github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/client/context"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/client/tendermint/tx"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// contains checks if the a given query contains one of the tx types
func contains(stringSlice []string, txType string) bool {
	for _, word := range stringSlice {
		if word == txType {
			return true
		}
	}
	return false
}

// queries staking txs
func queryTxs(cliCtx context.CLIContext, cdc *codec.Codec, tag string, delegatorAddr string) ([]tx.Info, error) {
	page := 0
	perPage := 100
	tags := []string{
		fmt.Sprintf("%s='%s'", tags.Action, tag),
		fmt.Sprintf("%s='%s'", tags.Delegator, delegatorAddr),
	}
	result, err := tx.SearchTxs(cliCtx, cdc, tags, page, perPage)
	if result != nil {
		return result.Txs, err
	}
	return []tx.Info{}, err
}

func queryBonds(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]
		bech32validator := vars["validatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		validatorAddr, err := sdk.ValAddressFromBech32(bech32validator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := stake.NewQueryBondsParams(delegatorAddr, validatorAddr)

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		switch endpoint {
		case "custom/stake/unbondingDelegation":
			var unbondingDelegation types.UnbondingDelegation
			if err = cdc.UnmarshalJSON(res, &unbondingDelegation); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
			if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationOutput); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		case "custom/stake/delegation":
			var delegation types.Delegation
			// parse out the validators
			if err = cdc.UnmarshalJSON(res, &delegation); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
			if res, err = codec.MarshalJSONIndent(cdc, delegationOutput); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		case "custom/stake/delegatorValidator":
			var validator types.Validator
			if err = cdc.UnmarshalJSON(res, &validator); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
			if res, err = codec.MarshalJSONIndent(cdc, validatorOutput); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryDelegator(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := stake.NewQueryDelegatorParams(delegatorAddr)

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		switch endpoint {
		case "custom/stake/delegatorDelegations":
			var delegations []types.Delegation
			// parse out the validators
			if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			delegationOutputs := make([]stakeClient.DelegationOutput, len(delegations))
			for index, delegation := range delegations {
				delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
				delegationOutputs[index] = delegationOutput
			}
			if res, err = codec.MarshalJSONIndent(cdc, delegationOutputs); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

		case "custom/stake/delegatorUnbondingDelegations":
			var unbondingDelegations []types.UnbondingDelegation
			if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			unbondingDelegationsOutputs := make([]stakeClient.UnbondingDelegationOutput, len(unbondingDelegations))
			for index, unbondingDelegation := range unbondingDelegations {
				unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
				unbondingDelegationsOutputs[index] = unbondingDelegationOutput
			}

			if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationsOutputs); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

		case "custom/stake/delegatorRedelegations":
			var relegations []types.Redelegation
			if err = cdc.UnmarshalJSON(res, &relegations); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			relegationsOutputs := make([]stakeClient.RedelegationOutput, len(relegations))
			for index, relegation := range relegations {
				relegationOutput := stakeClient.ConvertREDToREDOutput(cliCtx, relegation)
				relegationsOutputs[index] = relegationOutput
			}

			if res, err = codec.MarshalJSONIndent(cdc, relegationsOutputs); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

		case "custom/stake/delegatorValidators":
			var validators []types.Validator
			if err = cdc.UnmarshalJSON(res, &validators); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			validatorOutputs := make([]stakeClient.ValidatorOutput, len(validators))
			for index, validator := range validators {
				validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
				validatorOutputs[index] = validatorOutput
			}

			if res, err = codec.MarshalJSONIndent(cdc, validatorOutputs); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryValidator(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32validatorAddr := vars["validatorAddr"]

		validatorAddr, err := sdk.ValAddressFromBech32(bech32validatorAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := stake.NewQueryValidatorParams(validatorAddr)
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(endpoint, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		switch endpoint {
		case "custom/stake/validator":
			var validator types.Validator
			if err = cdc.UnmarshalJSON(res, &validator); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
			if res, err = codec.MarshalJSONIndent(cdc, validatorOutput); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

		case "custom/stake/validatorUnbondingDelegations":
			var unbondingDelegations []types.UnbondingDelegation
			if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			unbondingDelegationsOutputs := make([]stakeClient.UnbondingDelegationOutput, len(unbondingDelegations))
			for index, unbondingDelegation := range unbondingDelegations {
				unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
				unbondingDelegationsOutputs[index] = unbondingDelegationOutput
			}

			if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationsOutputs); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

		case "custom/stake/validatorRedelegations":
			var redelegations []types.Redelegation
			if err = cdc.UnmarshalJSON(res, &redelegations); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			redelegationsOutputs := make([]stakeClient.RedelegationOutput, len(redelegations))
			for index, redelegation := range redelegations {
				redelegationOutput := stakeClient.ConvertREDToREDOutput(cliCtx, redelegation)
				redelegationsOutputs[index] = redelegationOutput
			}

			if res, err = codec.MarshalJSONIndent(cdc, redelegationsOutputs); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		case "custom/stake/validatorDelegations":
			var delegations []types.Delegation
			if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			delegationOutputs := make([]stakeClient.DelegationOutput, len(delegations))
			for index, delegation := range delegations {
				delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
				delegationOutputs[index] = delegationOutput
			}

			if res, err = codec.MarshalJSONIndent(cdc, delegationOutputs); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func ConvertPaginationParams(pageString, sizeString string) (paginationParams sdk.PaginationParams, err error) {
	page := uint64(0)
	size := uint16(100)
	if pageString != "" {
		page, err = strconv.ParseUint(pageString, 10, 64)
		if err != nil {
			err = fmt.Errorf("page '%s' is not a valid uint64", pageString)
			return paginationParams, err
		}
	}
	if sizeString != "" {
		sizeUint64, err := strconv.ParseUint(sizeString, 10, 16)
		if err != nil {
			err = fmt.Errorf("size '%s' is not a valid uint16", sizeString)
			return paginationParams, err
		}
		size = uint16(sizeUint64)
	}
	paginationParams = sdk.NewPaginationParams(page, size)
	return paginationParams, err
}
