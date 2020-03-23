package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/distribution"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
)

// GetWithdrawAddress returns withdraw address of a given delegator address
func GetWithdrawAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "withdraw-address [account-address]",
		Short:   "Query withdraw address",
		Example: "iriscli distribution withdraw-address <account-address>",
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

// GetRewards returns the all the rewards of validator or delegator
func GetRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rewards [address]",
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
