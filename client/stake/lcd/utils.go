package lcd

import (
	"bytes"
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/x/stake/tags"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/irisnet/irishub/client/context"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/client/tendermint/tx"
	"github.com/pkg/errors"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
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

func getDelegatorValidator(cliCtx context.CLIContext, cdc *wire.Codec, delegatorAddr sdk.AccAddress, validatorAccAddr sdk.AccAddress) (
	validator stakeClient.ValidatorOutput, httpStatusCode int, errMsg string, err error) {

	// check if the delegator is bonded or redelegated to the validator
	keyDel := stake.GetDelegationKey(delegatorAddr, validatorAccAddr)

	res, err := cliCtx.QueryStore(keyDel, storeName)
	if err != nil {
		return stakeClient.ValidatorOutput{}, http.StatusInternalServerError, "couldn't query delegation. Error: ", err
	}

	if len(res) == 0 {
		return stakeClient.ValidatorOutput{}, http.StatusNoContent, "", nil
	}

	kvs, err := cliCtx.QuerySubspace(stake.ValidatorsKey, storeName)
	if err != nil {
		return stakeClient.ValidatorOutput{}, http.StatusInternalServerError, "Error: ", err
	}
	if len(kvs) == 0 {
		// the query will return empty if there are no delegations
		return stakeClient.ValidatorOutput{}, http.StatusNoContent, "", nil
	}

	validator, errVal := getValidatorFromAccAdrr(validatorAccAddr, kvs, cliCtx, cdc)
	if errVal != nil {
		return stakeClient.ValidatorOutput{}, http.StatusInternalServerError, "Couldn't get info from validator. Error: ", errVal
	}
	return validator, http.StatusOK, "", nil
}

func getDelegatorDelegations(cliCtx context.CLIContext, cdc *wire.Codec, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress) (
	outputDelegation stakeClient.DelegationOutput, httpStatusCode int, errMsg string, err error) {
	delegationKey := stake.GetDelegationKey(delegatorAddr, validatorAddr)
	marshalledDelegation, err := cliCtx.QueryStore(delegationKey, storeName)
	if err != nil {
		return stakeClient.DelegationOutput{}, http.StatusInternalServerError, "couldn't query delegation. Error: ", err
	}

	// the query will return empty if there is no data for this record
	if len(marshalledDelegation) == 0 {
		return stakeClient.DelegationOutput{}, http.StatusNoContent, "", nil
	}

	delegation, err := types.UnmarshalDelegation(cdc, delegationKey, marshalledDelegation)
	if err != nil {
		return stakeClient.DelegationOutput{}, http.StatusInternalServerError, "couldn't unmarshall delegation. Error: ", err
	}

	outputDelegation = stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)

	return outputDelegation, http.StatusOK, "", nil
}

func getDelegatorUndelegations(cliCtx context.CLIContext, cdc *wire.Codec, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress) (
	unbonds types.UnbondingDelegation, httpStatusCode int, errMsg string, err error) {
	undelegationKey := stake.GetUBDKey(delegatorAddr, validatorAddr)
	marshalledUnbondingDelegation, err := cliCtx.QueryStore(undelegationKey, storeName)
	if err != nil {
		return types.UnbondingDelegation{}, http.StatusInternalServerError, "couldn't query unbonding-delegation. Error: ", err
	}

	// the query will return empty if there is no data for this record
	if len(marshalledUnbondingDelegation) == 0 {
		return types.UnbondingDelegation{}, http.StatusNoContent, "", nil
	}

	unbondingDelegation, err := types.UnmarshalUBD(cdc, undelegationKey, marshalledUnbondingDelegation)
	if err != nil {
		return types.UnbondingDelegation{}, http.StatusInternalServerError, "couldn't unmarshall unbonding-delegation. Error: ", err
	}
	return unbondingDelegation, http.StatusOK, "", nil
}

func getDelegatorRedelegations(cliCtx context.CLIContext, cdc *wire.Codec, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress) (
	regelegations types.Redelegation, httpStatusCode int, errMsg string, err error) {

	keyRedelegateTo := stake.GetREDsByDelToValDstIndexKey(delegatorAddr, validatorAddr)
	marshalledRedelegations, err := cliCtx.QueryStore(keyRedelegateTo, storeName)
	if err != nil {
		return types.Redelegation{}, http.StatusInternalServerError, "couldn't query redelegation. Error: ", err
	}

	if len(marshalledRedelegations) == 0 {
		return types.Redelegation{}, http.StatusNoContent, "", nil
	}

	redelegations, err := types.UnmarshalRED(cdc, keyRedelegateTo, marshalledRedelegations)
	if err != nil {
		return types.Redelegation{}, http.StatusInternalServerError, "couldn't unmarshall redelegations. Error: ", err
	}

	return redelegations, http.StatusOK, "", nil
}

// queries staking txs
func queryTxs(node rpcclient.Client, cliCtx context.CLIContext, cdc *wire.Codec, tag string, delegatorAddr string) ([]tx.Info, error) {
	page := 0
	perPage := 100
	prove := !cliCtx.TrustNode
	query := fmt.Sprintf("%s='%s' AND %s='%s'", tags.Action, tag, tags.Delegator, delegatorAddr)
	res, err := node.TxSearch(query, prove, page, perPage)
	if err != nil {
		return nil, err
	}

	return tx.FormatTxResults(cdc, res.Txs)
}

// gets all validators
func getValidators(cliCtx context.CLIContext, cdc *wire.Codec, validatorKVs []sdk.KVPair) ([]stakeClient.ValidatorOutput, error) {
	validators := make([]stakeClient.ValidatorOutput, len(validatorKVs))
	for i, kv := range validatorKVs {

		addr := kv.Key[1:]
		validator, err := types.UnmarshalValidator(cdc, addr, kv.Value)
		if err != nil {
			return nil, err
		}

		validatorOutput, err := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
		if err != nil {
			return nil, err
		}
		validators[i] = validatorOutput
	}
	return validators, nil
}

// gets a validator given a ValAddress
func getValidator(address sdk.AccAddress, validatorKVs []sdk.KVPair, cliCtx context.CLIContext, cdc *wire.Codec) (stakeClient.ValidatorOutput, error) {
	// parse out the validators
	for _, kv := range validatorKVs {
		addr := kv.Key[1:]
		validator, err := types.UnmarshalValidator(cdc, addr, kv.Value)
		if err != nil {
			return stakeClient.ValidatorOutput{}, err
		}

		ownerAddress := validator.Owner
		if bytes.Equal(ownerAddress.Bytes(), address.Bytes()) {
			validatorOutput, err := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
			if err != nil {
				return stakeClient.ValidatorOutput{}, err
			}

			return validatorOutput, nil
		}
	}
	return stakeClient.ValidatorOutput{}, errors.Errorf("Couldn't find validator")
}

// gets a validator given an AccAddress
func getValidatorFromAccAdrr(address sdk.AccAddress, validatorKVs []sdk.KVPair, cliCtx context.CLIContext, cdc *wire.Codec) (stakeClient.ValidatorOutput, error) {
	// parse out the validators
	for _, kv := range validatorKVs {
		addr := kv.Key[1:]
		validator, err := types.UnmarshalValidator(cdc, addr, kv.Value)
		if err != nil {
			return stakeClient.ValidatorOutput{}, err
		}

		ownerAddress := validator.Owner
		if bytes.Equal(ownerAddress.Bytes(), address.Bytes()) {
			validatorOutput, err := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
			if err != nil {
				return stakeClient.ValidatorOutput{}, err
			}

			return validatorOutput, nil
		}
	}
	return stakeClient.ValidatorOutput{}, errors.Errorf("Couldn't find validator")
}

//  gets all Bech32 validators from a key
func getValidatorOutputs(storeName string, cliCtx context.CLIContext, cdc *wire.Codec) (
	validators []stakeClient.ValidatorOutput, httpStatusCode int, errMsg string, err error) {
	// Get all validators using key
	kvs, err := cliCtx.QuerySubspace(stake.ValidatorsKey, storeName)
	if err != nil {
		return nil, http.StatusInternalServerError, "couldn't query validators. Error: ", err
	}

	// the query will return empty if there are no validators
	if len(kvs) == 0 {
		return nil, http.StatusNoContent, "", nil
	}

	validators, err = getValidators(cliCtx, cdc, kvs)
	if err != nil {
		return nil, http.StatusInternalServerError, "Error: ", err
	}
	return validators, http.StatusOK, "", nil
}
