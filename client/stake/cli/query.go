package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/client/context"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryValidator implements the validator query command.
func GetCmdQueryValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validator [validator-address]",
		Short:   "Query a validator",
		Example: "iriscli stake validator <validator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := stake.NewQueryValidatorParams(addr)

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.StakeRoute, stake.QueryValidator)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var validator stake.Validator
			err = cdc.UnmarshalJSON(res, &validator)
			if err != nil {
				return err
			}

			validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)

			return cliCtx.PrintOutput(validatorOutput)
		},
	}

	return cmd
}

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validators",
		Short:   "Query for all validators",
		Example: "iriscli stake validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := stake.ValidatorsKey
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(key, protocol.StakeRoute)
			if err != nil {
				return err
			}

			var validatorsOutput stakeClient.ValidatorsOutput
			for _, kv := range resKVs {
				addr := kv.Key[1:]
				validator := types.MustUnmarshalValidator(cdc, addr, kv.Value)
				validatorsOutput = append(validatorsOutput, stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator))
			}

			return cliCtx.PrintOutput(validatorsOutput)
		},
	}

	return cmd
}

// GetCmdQueryValidatorUnbondingDelegations implements the query all unbonding delegatations from a validator command.
func GetCmdQueryValidatorUnbondingDelegations(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegations-from [validator-address]",
		Short:   "Query all unbonding delegatations from a validator",
		Example: "iriscli stake unbonding-delegations-from <validator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := stake.NewQueryValidatorParams(valAddr)

			return queryValidator(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryValidatorUnbondingDelegations, params)
		},
	}
	return cmd
}

// GetCmdQueryValidatorRedelegations implements the query all redelegatations from a validator command.
func GetCmdQueryValidatorRedelegations(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "redelegations-from [validator-address]",
		Short:   "Query all outgoing redelegatations from a validator",
		Example: "iriscli stake redelegations-from <validator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := stake.NewQueryValidatorParams(valAddr)

			return queryValidator(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryValidatorRedelegations, params)
		},
	}
	return cmd
}

// GetCmdQueryDelegation the query delegation command.
func GetCmdQueryDelegation(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegation",
		Short:   "Query a delegation based on address and validator address",
		Example: "iriscli stake delegation --address-validator=<validator address> --address-delegator=<delegator address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			valAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidator))
			if err != nil {
				return err
			}

			delAddr, err := sdk.AccAddressFromBech32(viper.GetString(FlagAddressDelegator))
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := stake.NewQueryBondsParams(delAddr, valAddr)

			return queryBonds(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryDelegation, params)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)

	return cmd
}

// GetCmdQueryDelegations implements the command to query all the delegations
// made from one delegator.
func GetCmdQueryDelegations(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegations [delegator-address]",
		Short:   "Query all delegations made from one delegator",
		Example: "iriscli stake delegations <delegator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			delegatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := stake.NewQueryDelegatorParams(delegatorAddr)

			return queryDelegator(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryDelegatorDelegations, params)
		},
	}

	return cmd
}

// GetCmdQueryValidatorDelegations implements the command to query all the
// delegations to a specific validator.
func GetCmdQueryValidatorDelegations(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegations-to [validator-address]",
		Short:   "Query all delegations made to one validator",
		Example: "iriscli stake delegations-to <validator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			validatorAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := stake.NewQueryValidatorParams(validatorAddr)

			return queryValidator(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryValidatorDelegations, params)
		},
	}
	return cmd
}

// GetCmdQueryUnbondingDelegation implements the command to query a single
// unbonding-delegation record.
func GetCmdQueryUnbondingDelegation(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegation",
		Short:   "Query an unbonding-delegation record based on delegator and validator address",
		Example: "iriscli stake unbonding-delegation --address-validator=<validator address> --address-delegator=<delegator address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			valAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidator))
			if err != nil {
				return err
			}

			delAddr, err := sdk.AccAddressFromBech32(viper.GetString(FlagAddressDelegator))
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := stake.NewQueryBondsParams(delAddr, valAddr)

			return queryBonds(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryUnbondingDelegation, params)
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)

	return cmd
}

// GetCmdQueryUnbondingDelegations implements the command to query all the
// unbonding-delegation records for a delegator.
func GetCmdQueryUnbondingDelegations(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbonding-delegations [delegator-address]",
		Short:   "Query all unbonding-delegations records for one delegator",
		Example: "iriscli stake unbonding-delegation <delegator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			delegatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := stake.NewQueryDelegatorParams(delegatorAddr)

			return queryDelegator(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryDelegatorUnbondingDelegations, params)
		},
	}

	return cmd
}

// GetCmdQueryRedelegation implements the command to query a single
// redelegation record.
func GetCmdQueryRedelegation(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "redelegation",
		Short:   "Query a redelegation record based on delegator and a source and destination validator address",
		Example: "iriscli stake redelegation --address-validator-source=<source validator address> --address-validator-dest=<destination validator address> --address-delegator=<delegator address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			valSrcAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidatorSrc))
			if err != nil {
				return err
			}

			valDstAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidatorDst))
			if err != nil {
				return err
			}

			delAddr, err := sdk.AccAddressFromBech32(viper.GetString(FlagAddressDelegator))
			if err != nil {
				return err
			}

			key := stake.GetREDKey(delAddr, valSrcAddr, valDstAddr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, protocol.StakeStore)
			if err != nil {
				return err
			}

			red, err := types.UnmarshalRED(cdc, key, res)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(stakeClient.ConvertREDToREDOutput(cliCtx, red))
		},
	}

	cmd.Flags().AddFlagSet(fsRedelegation)
	cmd.Flags().AddFlagSet(fsDelegator)

	return cmd
}

// GetCmdQueryRedelegations implements the command to query all the
// redelegation records for a delegator.
func GetCmdQueryRedelegations(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "redelegations [delegator-address]",
		Short:   "Query all redelegations records for one delegator",
		Example: "iriscli stake redelegations <delegator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			delegatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := stake.NewQueryDelegatorParams(delegatorAddr)

			return queryDelegator(cliCtx, fmt.Sprintf("custom/%s", protocol.StakeRoute),
				stake.QueryDelegatorRedelegations, params)
		},
	}

	return cmd
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Short:   "Query the current staking pool values",
		Example: "iriscli stake pool",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.StakeRoute, stake.QueryPool), nil)
			if err != nil {
				return err
			}
			var poolStatus stake.PoolStatus
			err = cdc.UnmarshalJSON(res, &poolStatus)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(stakeClient.ConvertPoolToPoolOutput(cliCtx, poolStatus))
		},
	}

	return cmd
}

// GetCmdQueryPool implements the params query command.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Short:   "Query the current staking parameters information",
		Example: "iriscli stake parameters",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
				protocol.StakeRoute, stake.QueryParameters), nil)
			if err != nil {
				return err
			}

			var params stake.Params
			err = cdc.UnmarshalJSON(bz, &params)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}
