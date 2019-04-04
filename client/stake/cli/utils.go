package cli

import (
	"fmt"

	"github.com/irisnet/irishub/client/context"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
)

func queryBonds(cliCtx context.CLIContext, route string, query string, params stake.QueryBondsParams) ([]byte, error) {
	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	res, err := cliCtx.QueryWithData(fmt.Sprintf("%s/%s", route, query), bz)
	if err != nil {
		return nil, err
	}

	switch query {
	case stake.QueryUnbondingDelegation:
		var unbondingDelegation types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegation); err != nil {
			return nil, err
		}
		unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
		if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationOutput); err != nil {
			return nil, err
		}
	case stake.QueryDelegation:
		var delegation types.Delegation
		// parse out the validators
		if err = cdc.UnmarshalJSON(res, &delegation); err != nil {
			return nil, err
		}
		delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
		if res, err = codec.MarshalJSONIndent(cdc, delegationOutput); err != nil {
			return nil, err
		}
	case stake.QueryDelegatorValidator:
		var validator types.Validator
		if err = cdc.UnmarshalJSON(res, &validator); err != nil {
			return nil, err
		}

		validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
		if res, err = codec.MarshalJSONIndent(cdc, validatorOutput); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func queryDelegator(cliCtx context.CLIContext, route string, query string, params stake.QueryDelegatorParams) ([]byte, error) {

	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	res, err := cliCtx.QueryWithData(fmt.Sprintf("%s/%s", route, query), bz)
	if err != nil {
		return nil, err
	}

	switch query {
	case stake.QueryDelegatorDelegations:
		var delegations []types.Delegation
		// parse out the validators
		if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
			return nil, err
		}

		delegationOutputs := make([]stakeClient.DelegationOutput, len(delegations))
		for index, delegation := range delegations {
			delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
			delegationOutputs[index] = delegationOutput
		}
		if res, err = codec.MarshalJSONIndent(cdc, delegationOutputs); err != nil {
			return nil, err
		}

	case stake.QueryDelegatorUnbondingDelegations:
		var unbondingDelegations []types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
			return nil, err
		}
		unbondingDelegationsOutputs := make([]stakeClient.UnbondingDelegationOutput, len(unbondingDelegations))
		for index, unbondingDelegation := range unbondingDelegations {
			unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
			unbondingDelegationsOutputs[index] = unbondingDelegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationsOutputs); err != nil {
			return nil, err
		}

	case stake.QueryDelegatorRedelegations:
		var relegations []types.Redelegation
		if err = cdc.UnmarshalJSON(res, &relegations); err != nil {
			return nil, err
		}

		relegationsOutputs := make([]stakeClient.RedelegationOutput, len(relegations))
		for index, relegation := range relegations {
			relegationOutput := stakeClient.ConvertREDToREDOutput(cliCtx, relegation)
			relegationsOutputs[index] = relegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, relegationsOutputs); err != nil {
			return nil, err
		}

	case stake.QueryDelegatorValidators:
		var validators []types.Validator
		if err = cdc.UnmarshalJSON(res, &validators); err != nil {
			return nil, err
		}

		validatorOutputs := make([]stakeClient.ValidatorOutput, len(validators))
		for index, validator := range validators {
			validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
			validatorOutputs[index] = validatorOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, validatorOutputs); err != nil {
			return nil, err
		}

	}

	return res, nil
}

func queryValidator(cliCtx context.CLIContext, route string, query string, params stake.QueryValidatorParams) ([]byte, error) {
	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	res, err := cliCtx.QueryWithData(fmt.Sprintf("%s/%s", route, query), bz)
	if err != nil {
		return nil, err
	}

	switch query {
	case stake.QueryValidator:
		var validator types.Validator
		if err = cdc.UnmarshalJSON(res, &validator); err != nil {
			return nil, err
		}

		validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
		if res, err = codec.MarshalJSONIndent(cdc, validatorOutput); err != nil {
			return nil, err
		}

	case stake.QueryValidatorUnbondingDelegations:
		var unbondingDelegations []types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
			return nil, err
		}
		unbondingDelegationsOutputs := make([]stakeClient.UnbondingDelegationOutput, len(unbondingDelegations))
		for index, unbondingDelegation := range unbondingDelegations {
			unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
			unbondingDelegationsOutputs[index] = unbondingDelegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationsOutputs); err != nil {
			return nil, err
		}

	case stake.QueryValidatorRedelegations:
		var redelegations []types.Redelegation
		if err = cdc.UnmarshalJSON(res, &redelegations); err != nil {
			return nil, err
		}

		redelegationsOutputs := make([]stakeClient.RedelegationOutput, len(redelegations))
		for index, redelegation := range redelegations {
			redelegationOutput := stakeClient.ConvertREDToREDOutput(cliCtx, redelegation)
			redelegationsOutputs[index] = redelegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, redelegationsOutputs); err != nil {
			return nil, err
		}
	case stake.QueryValidatorDelegations:
		var delegations []types.Delegation
		if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
			return nil, err
		}

		delegationOutputs := make([]stakeClient.DelegationOutput, len(delegations))
		for index, delegation := range delegations {
			delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
			delegationOutputs[index] = delegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, delegationOutputs); err != nil {
			return nil, err
		}
	}
	return res, nil
}
