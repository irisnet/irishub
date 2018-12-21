package cli

import (
	"fmt"

	"github.com/irisnet/irishub/client/context"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// GetCmdQueryValidator implements the validator query command.
func GetCmdQueryValidator(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validator [validator-address]",
		Short:   "Query a validator",
		Example: "iriscli stake validator <validator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			key := stake.GetValidatorKey(addr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(res) == 0 {
				return fmt.Errorf("No validator found with address %s", args[0])
			}

			validator := types.MustUnmarshalValidator(cdc, addr, res)
			validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
			if err != nil {
				return err
			}

			switch viper.Get(cli.OutputFlag) {
			case "text":
				human, err := validatorOutput.HumanReadableString()
				if err != nil {
					return err
				}
				fmt.Println(human)

			case "json":
				// parse out the validator
				output, err := codec.MarshalJSONIndent(cdc, validatorOutput)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
			}

			// TODO: output with proofs / machine parseable etc.
			return nil
		},
	}

	return cmd
}

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validators",
		Short:   "Query for all validators",
		Example: "iriscli stake validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := stake.ValidatorsKey
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(key, storeName)
			if err != nil {
				return err
			}

			// parse out the validators
			var validators []stakeClient.ValidatorOutput
			for _, kv := range resKVs {
				addr := kv.Key[1:]
				validator := types.MustUnmarshalValidator(cdc, addr, kv.Value)
				validatorOutput := stakeClient.ConvertValidatorToValidatorOutput(cliCtx, validator)
				if err != nil {
					return err
				}
				validators = append(validators, validatorOutput)
			}

			switch viper.Get(cli.OutputFlag) {
			case "text":
				for _, validator := range validators {
					resp, err := validator.HumanReadableString()
					if err != nil {
						return err
					}

					fmt.Println(resp)
				}
			case "json":
				output, err := codec.MarshalJSONIndent(cdc, validators)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
				return nil
			}

			// TODO: output with proofs / machine parseable etc.
			return nil
		},
	}

	return cmd
}

// GetCmdQueryValidatorUnbondingDelegations implements the query all unbonding delegatations from a validator command.
func GetCmdQueryValidatorUnbondingDelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/validatorUnbondingDelegations", queryRoute),
				bz)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetCmdQueryValidatorRedelegations implements the query all redelegatations from a validator command.
func GetCmdQueryValidatorRedelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/validatorRedelegations", queryRoute),
				bz)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetCmdQueryDelegation the query delegation command.
func GetCmdQueryDelegation(storeName string, cdc *codec.Codec) *cobra.Command {
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

			key := stake.GetDelegationKey(delAddr, valAddr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(res) == 0 {
				return fmt.Errorf("no delegation found with delegator %s on validator %s", delAddr, valAddr)
			}

			// parse out the unbonding delegation
			delegation := types.MustUnmarshalDelegation(cdc, key, res)
			delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
			switch viper.Get(cli.OutputFlag) {
			case "text":
				resp, err := delegationOutput.HumanReadableString()
				if err != nil {
					return err
				}

				fmt.Println(resp)
			case "json":
				output, err := codec.MarshalJSONIndent(cdc, delegationOutput)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
				return nil
			}

			return nil
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)

	return cmd
}

// GetCmdQueryDelegations implements the command to query all the delegations
// made from one delegator.
func GetCmdQueryDelegations(storeName string, cdc *codec.Codec) *cobra.Command {
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

			key := stake.GetDelegationsKey(delegatorAddr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(key, storeName)
			if err != nil {
				return err
			}

			// parse out the validators
			var delegations []stakeClient.DelegationOutput
			for _, kv := range resKVs {
				delegation := types.MustUnmarshalDelegation(cdc, kv.Key, kv.Value)
				delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
				delegations = append(delegations, delegationOutput)
			}

			output, err := codec.MarshalJSONIndent(cdc, delegations)
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			// TODO: output with proofs / machine parseable etc.
			return nil
		},
	}

	return cmd
}

// GetCmdQueryValidatorDelegations implements the command to query all the
// delegations to a specific validator.
func GetCmdQueryValidatorDelegations(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
			params := stake.NewQueryValidatorParams(validatorAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/validatorDelegations", queryRoute), bz)
			if err != nil {
				return err
			}
			var delegations []stake.Delegation
			err = cdc.UnmarshalJSON(res, &delegations)
			if err != nil {
				return err
			}
			// parse out the validators
			var delegationsOutput []stakeClient.DelegationOutput
			for _, delegation := range delegations {
				delegationOutput := stakeClient.ConvertDelegationToDelegationOutput(cliCtx, delegation)
				delegationsOutput = append(delegationsOutput, delegationOutput)
			}
			output, err := codec.MarshalJSONIndent(cdc, delegationsOutput)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}

// GetCmdQueryUnbondingDelegation implements the command to query a single
// unbonding-delegation record.
func GetCmdQueryUnbondingDelegation(storeName string, cdc *codec.Codec) *cobra.Command {
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

			key := stake.GetUBDKey(delAddr, valAddr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(res) == 0 {
				return fmt.Errorf("no unbonding-delegation found with delegator %s from validator %s", delAddr, valAddr)
			}

			// parse out the unbonding delegation
			ubd := types.MustUnmarshalUBD(cdc, key, res)
			ubdOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, ubd)

			switch viper.Get(cli.OutputFlag) {
			case "text":
				resp, err := ubdOutput.HumanReadableString()
				if err != nil {
					return err
				}

				fmt.Println(resp)
			case "json":
				output, err := codec.MarshalJSONIndent(cdc, ubdOutput)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
				return nil
			}

			return nil
		},
	}

	cmd.Flags().AddFlagSet(fsValidator)
	cmd.Flags().AddFlagSet(fsDelegator)

	return cmd
}

// GetCmdQueryUnbondingDelegations implements the command to query all the
// unbonding-delegation records for a delegator.
func GetCmdQueryUnbondingDelegations(storeName string, cdc *codec.Codec) *cobra.Command {
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

			key := stake.GetUBDsKey(delegatorAddr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(key, storeName)
			if err != nil {
				return err
			}

			// parse out the validators
			var ubds []stakeClient.UnbondingDelegationOutput
			for _, kv := range resKVs {
				ubd := types.MustUnmarshalUBD(cdc, kv.Key, kv.Value)
				ubdOutput := stakeClient.ConvertUBDToUBDOutput(cliCtx, ubd)
				ubds = append(ubds, ubdOutput)
			}

			output, err := codec.MarshalJSONIndent(cdc, ubds)
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			// TODO: output with proofs / machine parseable etc.
			return nil
		},
	}

	return cmd
}

// GetCmdQueryRedelegation implements the command to query a single
// redelegation record.
func GetCmdQueryRedelegation(storeName string, cdc *codec.Codec) *cobra.Command {
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

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			} else if len(res) == 0 {
				return fmt.Errorf("no redelegation found with delegator %s from source validator %s to destination validator %s",
					delAddr, valSrcAddr, valDstAddr)
			}

			// parse out the unbonding delegation
			red := types.MustUnmarshalRED(cdc, key, res)
			redOutput := stakeClient.ConvertREDToREDOutput(cliCtx, red)
			switch viper.Get(cli.OutputFlag) {
			case "text":
				resp, err := redOutput.HumanReadableString()
				if err != nil {
					return err
				}

				fmt.Println(resp)
			case "json":
				output, err := codec.MarshalJSONIndent(cdc, redOutput)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
				return nil
			}

			return nil
		},
	}

	cmd.Flags().AddFlagSet(fsRedelegation)
	cmd.Flags().AddFlagSet(fsDelegator)

	return cmd
}

// GetCmdQueryRedelegations implements the command to query all the
// redelegation records for a delegator.
func GetCmdQueryRedelegations(storeName string, cdc *codec.Codec) *cobra.Command {
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

			key := stake.GetREDsKey(delegatorAddr)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(key, storeName)
			if err != nil {
				return err
			}

			// parse out the validators
			var reds []stakeClient.RedelegationOutput
			for _, kv := range resKVs {
				red := types.MustUnmarshalRED(cdc, kv.Key, kv.Value)
				redOutput := stakeClient.ConvertREDToREDOutput(cliCtx, red)
				reds = append(reds, redOutput)
			}

			output, err := codec.MarshalJSONIndent(cdc, reds)
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			// TODO: output with proofs / machine parseable etc.
			return nil
		},
	}

	return cmd
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Short:   "Query the current staking pool values",
		Example: "iriscli stake pool",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			key := stake.PoolKey
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}

			pool := types.MustUnmarshalPool(cdc, res)
			poolOutput := stakeClient.ConvertPoolToPoolOutput(cliCtx, pool)

			switch viper.Get(cli.OutputFlag) {
			case "text":
				human := poolOutput.HumanReadableString()

				fmt.Println(human)

			case "json":
				// parse out the pool
				output, err := codec.MarshalJSONIndent(cdc, poolOutput)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
			}
			return nil
		},
	}

	return cmd
}

// GetCmdQueryPool implements the params query command.
func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Short:   "Query the current staking parameters information",
		Example: "iriscli stake parameters",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, err := cliCtx.QueryWithData("custom/stake/"+stake.QueryParameters, nil)
			if err != nil {
				return err
			}

			var params stake.Params
			err = cdc.UnmarshalJSON(bz, &params)
			if err != nil {
				return err
			}

			switch viper.Get(cli.OutputFlag) {
			case "text":
				human := params.HumanReadableString()

				fmt.Println(human)

			case "json":
				// parse out the params
				output, err := codec.MarshalJSONIndent(cdc, params)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
			}
			return nil
		},
	}

	return cmd
}
