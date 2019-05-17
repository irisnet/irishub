package cli

import (
	"fmt"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/distribution"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	distrClient "github.com/irisnet/irishub/client/distribution"
	"github.com/irisnet/irishub/app/protocol"
)

// GetWithdrawAddress returns withdraw address of a given delegator address
func GetWithdrawAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "withdraw-address",
		Short:   "Query withdraw address",
		Example: "iriscli distribution withdraw-address <account address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			addrString := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr, err := sdk.AccAddressFromBech32(addrString)
			if err != nil {
				return err
			}

			params := distribution.NewQueryDelegatorParams(delAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryWithdrawAddr),
				bz)
			if err != nil {
				return err
			}

			var acc sdk.AccAddress
			err = cdc.UnmarshalJSON(res, &acc)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(acc)
		},
	}
}

// GetDelegationDistInfo returns the delegation distribution information of a given delegation
func GetDelegationDistInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegation-distr-info",
		Short:   "Query delegation distribution information",
		Example: "iriscli distribution delegation-distr-info --address-delegator=<delegator address> --address-validator=<validator address>",
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

			params := distribution.NewQueryDelegationDistInfoParams(delAddr, valAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryDelegationDistInfo),
				bz)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().String(FlagAddressDelegator, "", "bech address of the delegator")
	cmd.Flags().String(FlagAddressValidator, "", "bech address of the validator")
	cmd.MarkFlagRequired(FlagAddressDelegator)
	cmd.MarkFlagRequired(FlagAddressValidator)
	return cmd
}

// GetAllDelegationDistInfo returns all delegation distribution information of a given delegator
func GetAllDelegationDistInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegator-distr-info",
		Short:   "Query delegator distribution information",
		Example: "iriscli distribution delegator-distr-info <delegator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addrString := args[0]

			delAddr, err := sdk.AccAddressFromBech32(addrString)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := distribution.NewQueryDelegatorParams(delAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryAllDelegationDistInfo),
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

// GetValidatorDistInfo returns the validator distribution information of a given validator
func GetValidatorDistInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validator-distr-info",
		Short:   "Query validator distribution information",
		Example: "iriscli distribution validator-distr-info <validator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			addrString := args[0]

			valAddr, err := sdk.ValAddressFromBech32(addrString)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := distribution.NewQueryValidatorDistInfoParams(valAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryValidatorDistInfo),
				bz)
			if err != nil {
				return err
			}

			var vddInfo distribution.ValidatorDistInfo
			if err = cliCtx.Codec.UnmarshalJSON(res, &vddInfo); err != nil {
				return err
			}

			output := distrClient.ConvertToValidatorDistInfoOutput(cliCtx, vddInfo)
			res, err = codec.MarshalJSONIndent(cdc, output)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetRewards returns the all the rewards of validator or delegator
func GetRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rewards",
		Short:   "Query all the rewards of validator or delegator",
		Example: "iriscli distribution rewards <address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			addrString := args[0]
			address, err := sdk.AccAddressFromBech32(addrString)
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := distribution.NewQueryRewardsParams(address)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryRewards),
				bz)
			if err != nil {
				return err
			}

			var rewards distribution.Rewards
			err = cdc.UnmarshalJSON(res, &rewards)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(rewards)
		},
	}
	return cmd
}

// GetCommunityTax returns the community tax coins
func GetCommunityTax(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "community-tax",
		Short:   "Query community tax",
		Example: "iriscli distribution community-tax",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", protocol.DistrRoute, distribution.QueryCommunityTax),
				nil)
			if err != nil {
				return err
			}

			var communityTax distribution.CommunityTax
			err = cdc.UnmarshalJSON(res, &communityTax)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(communityTax)
		},
	}
	return cmd
}
