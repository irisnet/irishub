package cli

import (
	"fmt"

	"github.com/irisnet/irishub/client/context"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
)

func queryBonds(cliCtx context.CLIContext, route string, query string, params stake.QueryBondsParams) error {
	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, err := cliCtx.QueryWithData(fmt.Sprintf("%s/%s", route, query), bz)
	if err != nil {
		return err
	}

	switch query {
	case stake.QueryUnbondingDelegation:
		var unbondingDelegation types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegation); err != nil {
			return err
		}
		unbondingDelegationOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation)
		return cliCtx.PrintOutput(unbondingDelegationOutput)

	case stake.QueryDelegation:
		var delegation types.Delegation
		// parse out the validators
		if err = cdc.UnmarshalJSON(res, &delegation); err != nil {
			return err
		}
		delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
		return cliCtx.PrintOutput(delegationOutput)

	case stake.QueryDelegatorValidator:
		var validator types.Validator
		if err = cdc.UnmarshalJSON(res, &validator); err != nil {
			return err
		}

		validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
		return cliCtx.PrintOutput(validatorOutput)
	}
	return nil
}

func queryDelegator(cliCtx context.CLIContext, route string, query string, params stake.QueryDelegatorParams) error {

	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, err := cliCtx.QueryWithData(fmt.Sprintf("%s/%s", route, query), bz)
	if err != nil {
		return err
	}

	switch query {
	case stake.QueryDelegatorDelegations:
		var delegations []types.Delegation
		// parse out the validators
		if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
			return err
		}

		var delegationsOutput stakeClient.DelegationsOutput
		for _, delegation := range delegations {
			delegationsOutput = append(delegationsOutput, stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation))
		}
		return cliCtx.PrintOutput(delegationsOutput)

	case stake.QueryDelegatorUnbondingDelegations:
		var unbondingDelegations []types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
			return err
		}
		var unbondingDelegationsOutput stakeClient.UnbondingDelegationsOutput
		for _, unbondingDelegation := range unbondingDelegations {
			unbondingDelegationsOutput = append(unbondingDelegationsOutput,
				stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation))
		}
		return cliCtx.PrintOutput(unbondingDelegationsOutput)

	case stake.QueryDelegatorRedelegations:
		var relegations []types.Redelegation
		if err = cdc.UnmarshalJSON(res, &relegations); err != nil {
			return err
		}

		var relegationsOutputs stakeClient.RedelegationsOutput
		for _, relegation := range relegations {
			relegationsOutputs = append(relegationsOutputs, stakeClient.ConvertREDToREDOutput(cliCtx, relegation))
		}
		return cliCtx.PrintOutput(relegationsOutputs)

	case stake.QueryDelegatorValidators:
		var validators []types.Validator
		if err = cdc.UnmarshalJSON(res, &validators); err != nil {
			return err
		}

		var validatorOutputs stakeClient.ValidatorsOutput
		for _, validator := range validators {
			validatorOutputs = append(validatorOutputs, stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator))
		}
		return cliCtx.PrintOutput(validatorOutputs)
	}

	return nil
}

func queryValidator(cliCtx context.CLIContext, route string, query string, params stake.QueryValidatorParams) error {
	cdc := cliCtx.Codec
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, err := cliCtx.QueryWithData(fmt.Sprintf("%s/%s", route, query), bz)
	if err != nil {
		return err
	}

	switch query {
	case stake.QueryValidator:
		var validator types.Validator
		if err = cdc.UnmarshalJSON(res, &validator); err != nil {
			return err
		}

		validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
		return cliCtx.PrintOutput(validatorOutput)

	case stake.QueryValidatorUnbondingDelegations:
		var unbondingDelegations []types.UnbondingDelegation
		if err = cdc.UnmarshalJSON(res, &unbondingDelegations); err != nil {
			return err
		}
		var unbondingDelegationsOutputs stakeClient.UnbondingDelegationsOutput
		for _, unbondingDelegation := range unbondingDelegations {
			unbondingDelegationsOutputs = append(unbondingDelegationsOutputs,
				stakeClient.ConvertUBDToUBDOutput(cliCtx, unbondingDelegation))
		}
		return cliCtx.PrintOutput(unbondingDelegationsOutputs)

	case stake.QueryValidatorRedelegations:
		var redelegations []types.Redelegation
		if err = cdc.UnmarshalJSON(res, &redelegations); err != nil {
			return err
		}

		var redelegationsOutputs stakeClient.RedelegationsOutput
		for _, redelegation := range redelegations {
			redelegationsOutputs = append(redelegationsOutputs,
				stakeClient.ConvertREDToREDOutput(cliCtx, redelegation))
		}
		return cliCtx.PrintOutput(redelegationsOutputs)

	case stake.QueryValidatorDelegations:
		var delegations []types.Delegation
		if err = cdc.UnmarshalJSON(res, &delegations); err != nil {
			return err
		}

		var delegationOutputs stakeClient.DelegationsOutput
		for _, delegation := range delegations {
			delegationOutputs = append(delegationOutputs,
				stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation))
		}
		return cliCtx.PrintOutput(delegationOutputs)
	}
	return nil
}
