package cli

import (
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"fmt"
)

func queryBonds(cliCtx context.CLIContext,endpoint string, params stake.QueryBondsParams) error {
	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, err := cliCtx.QueryWithData(endpoint, bz)
	if err != nil {
		return err
	}

	switch endpoint {
	case "custom/stake/unbondingDelegation":
		var unbondingDelegation types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegation); err != nil {
			return err
		}
		unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
		if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationOutput); err != nil {
			return err
		}
	case "custom/stake/delegation":
		var delegation types.Delegation
		// parse out the validators
		if err = cdc.UnmarshalJSON(res, &delegation); err != nil {
			return err
		}
		delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
		if res, err = codec.MarshalJSONIndent(cdc, delegationOutput); err != nil {
			return err
		}
	case "custom/stake/delegatorValidator":
		var validator types.Validator
		if err = cdc.UnmarshalJSON(res, &validator); err != nil {
			return err
		}

		validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
		if res, err = codec.MarshalJSONIndent(cdc, validatorOutput); err != nil {
			return err
		}
	}
	fmt.Println(string(res))
	return nil
}

func queryDelegator(cliCtx context.CLIContext, endpoint string, params stake.QueryDelegatorParams) (error) {

	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, err := cliCtx.QueryWithData(endpoint, bz)
	if err != nil {
		return err
	}

	switch endpoint {
	case "custom/stake/delegatorDelegations":
		var delegations []types.Delegation
		// parse out the validators
		if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
			return err
		}

		delegationOutputs := make([]stakeClient.DelegationOutput, len(delegations))
		for index, delegation := range delegations {
			delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
			delegationOutputs[index] = delegationOutput
		}
		if res, err = codec.MarshalJSONIndent(cdc, delegationOutputs); err != nil {
			return err
		}

	case "custom/stake/delegatorUnbondingDelegations":
		var unbondingDelegations []types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
			return err
		}
		unbondingDelegationsOutputs := make([]stakeClient.UnbondingDelegationOutput, len(unbondingDelegations))
		for index, unbondingDelegation := range unbondingDelegations {
			unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
			unbondingDelegationsOutputs[index] = unbondingDelegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationsOutputs); err != nil {
			return err
		}

	case "custom/stake/delegatorRedelegations":
		var relegations []types.Redelegation
		if err = cdc.UnmarshalJSON(res, &relegations); err != nil {
			return err
		}

		relegationsOutputs := make([]stakeClient.RedelegationOutput, len(relegations))
		for index, relegation := range relegations {
			relegationOutput := stakeClient.ConvertREDToREDOutput(cliCtx, relegation)
			relegationsOutputs[index] = relegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, relegationsOutputs); err != nil {
			return err
		}

	case "custom/stake/delegatorValidators":
		var validators []types.Validator
		if err = cdc.UnmarshalJSON(res, &validators); err != nil {
			return err
		}

		validatorOutputs := make([]stakeClient.ValidatorOutput, len(validators))
		for index, validator := range validators {
			validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
			validatorOutputs[index] = validatorOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, validatorOutputs); err != nil {
			return err
		}

	}

	fmt.Println(string(res))
	return nil
}

func queryValidator(cliCtx context.CLIContext,endpoint string, params stake.QueryValidatorParams) error {
	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, err := cliCtx.QueryWithData(endpoint, bz)
	if err != nil {
		return err
	}

	switch endpoint {
	case "custom/stake/validator":
		var validator types.Validator
		if err = cdc.UnmarshalJSON(res, &validator); err != nil {
			return err
		}

		validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
		if res, err = codec.MarshalJSONIndent(cdc, validatorOutput); err != nil {
			return err
		}

	case "custom/stake/validatorUnbondingDelegations":
		var unbondingDelegations []types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
			return err
		}
		unbondingDelegationsOutputs := make([]stakeClient.UnbondingDelegationOutput, len(unbondingDelegations))
		for index, unbondingDelegation := range unbondingDelegations {
			unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
			unbondingDelegationsOutputs[index] = unbondingDelegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, unbondingDelegationsOutputs); err != nil {
			return err
		}

	case "custom/stake/validatorRedelegations":
		var redelegations []types.Redelegation
		if err = cdc.UnmarshalJSON(res, &redelegations); err != nil {
			return err
		}

		redelegationsOutputs := make([]stakeClient.RedelegationOutput, len(redelegations))
		for index, redelegation := range redelegations {
			redelegationOutput := stakeClient.ConvertREDToREDOutput(cliCtx, redelegation)
			redelegationsOutputs[index] = redelegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, redelegationsOutputs); err != nil {
			return err
		}
	case "custom/stake/validatorDelegations":
		var delegations []types.Delegation
		if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
			return err
		}

		delegationOutputs := make([]stakeClient.DelegationOutput, len(delegations))
		for index, delegation := range delegations {
			delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
			delegationOutputs[index] = delegationOutput
		}

		if res, err = codec.MarshalJSONIndent(cdc, delegationOutputs); err != nil {
			return err
		}
	}
	fmt.Println(string(res))
	return nil
}
