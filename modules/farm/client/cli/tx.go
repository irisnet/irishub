package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"mods.irisnet.org/modules/farm/types"
)

// NewTxCmd returns the transaction commands for the farm module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Record transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdCreateFarmPool(),
		GetCreatePoolWithCommunityPool(),
		GetCmdDestroyFarmPool(),
		GetCmdAdjustPool(),
		GetCmdStake(),
		GetCmdUnstake(),
		GetCmdHarvest(),
	)
	return txCmd
}

// GetCmdCreateFarmPool implements the create a new farm pool command.
func GetCmdCreateFarmPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new farm pool",
		Example: fmt.Sprintf("$ %s tx farm create [flags]", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			description, _ := cmd.Flags().GetString(FlagDescription)
			lpTokenDenom, _ := cmd.Flags().GetString(FlagLPTokenDenom)
			startHeight, err := cmd.Flags().GetInt64(FlagStartHeight)
			if err != nil {
				return err
			}
			editable, _ := cmd.Flags().GetBool(FlagEditable)

			rewardPerBlockStr, _ := cmd.Flags().GetString(FlagRewardPerBlock)
			rewardPerBlock, err := sdk.ParseCoinsNormalized(rewardPerBlockStr)
			if err != nil {
				return err
			}

			totalRewardStr, _ := cmd.Flags().GetString(FlagTotalReward)
			totalReward, err := sdk.ParseCoinsNormalized(totalRewardStr)
			if err != nil {
				return err
			}

			msg := types.MsgCreatePool{
				Description:    description,
				LptDenom:       lpTokenDenom,
				StartHeight:    startHeight,
				RewardPerBlock: rewardPerBlock,
				TotalReward:    totalReward,
				Editable:       editable,
				Creator:        clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			if res.Params.PoolCreationFee.IsPositive() {
				fmt.Printf(
					"The farm creation transaction will consume extra fee: %s\n",
					res.Params.PoolCreationFee.String(),
				)
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreateFarmPool)
	_ = cmd.MarkFlagRequired(FlagStartHeight)
	_ = cmd.MarkFlagRequired(FlagRewardPerBlock)
	_ = cmd.MarkFlagRequired(FlagLPTokenDenom)
	_ = cmd.MarkFlagRequired(FlagTotalReward)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCreatePoolWithCommunityPool implements the create a new farm pool with communityPool command.
func GetCreatePoolWithCommunityPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-proposal",
		Short:   "Create a gov proposal to create a farm pool with community funds",
		Example: fmt.Sprintf("$ %s tx farm create-proposal [flags]", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolDescription, _ := cmd.Flags().GetString(FlagDescription)
			description, _ := cmd.Flags().GetString(FlagProposalDescription)
			title, _ := cmd.Flags().GetString(FlagProposalTitle)
			lpTokenDenom, _ := cmd.Flags().GetString(FlagLPTokenDenom)

			rewardPerBlockStr, _ := cmd.Flags().GetString(FlagRewardPerBlock)
			rewardPerBlock, err := sdk.ParseCoinsNormalized(rewardPerBlockStr)
			if err != nil {
				return err
			}

			depositStr, _ := cmd.Flags().GetString(FlagProposaldeposit)
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			fundAppliedStr, _ := cmd.Flags().GetString(FlagFundApplied)
			fundApplied, err := sdk.ParseCoinsNormalized(fundAppliedStr)
			if err != nil {
				return err
			}

			fundSelfBondStr, _ := cmd.Flags().GetString(FlagFundSelfBond)
			fundSelfBond, err := sdk.ParseCoinsNormalized(fundSelfBondStr)
			if err != nil {
				return err
			}

			msg := types.MsgCreatePoolWithCommunityPool{
				Content: types.CommunityPoolCreateFarmProposal{
					Title:           title,
					Description:     description,
					PoolDescription: poolDescription,
					LptDenom:        lpTokenDenom,
					RewardPerBlock:  rewardPerBlock,
					FundApplied:     fundApplied,
					FundSelfBond:    fundSelfBond,
				},
				InitialDeposit: deposit,
				Proposer:       clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreatePoolWithCommunityPool)
	_ = cmd.MarkFlagRequired(FlagRewardPerBlock)
	_ = cmd.MarkFlagRequired(FlagLPTokenDenom)
	_ = cmd.MarkFlagRequired(FlagFundApplied)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDestroyFarmPool implements the destroy a farm pool command.
func GetCmdDestroyFarmPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy",
		Short:   "Destroy a new farm pool",
		Example: fmt.Sprintf("$ %s tx farm destroy <Farm Pool ID> [flags]", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.MsgDestroyPool{
				PoolId:  args[0],
				Creator: clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdAdjustPool implements the append some reward for farm pool command.
func GetCmdAdjustPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "adjust",
		Short:   "Adjust farm pool parameters",
		Example: fmt.Sprintf("$ %s tx farm adjust <pool-id> [flags]", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var additionalReward, rewardPerBlock sdk.Coins
			if cmd.Flags().Changed(FlagRewardPerBlock) {
				rewardPerBlockStr, _ := cmd.Flags().GetString(FlagRewardPerBlock)
				rewardPerBlock, err = sdk.ParseCoinsNormalized(rewardPerBlockStr)
				if err != nil {
					return err
				}
			}
			if cmd.Flags().Changed(FlagAdditionalReward) {
				additionalRewardStr, _ := cmd.Flags().GetString(FlagAdditionalReward)
				additionalReward, err = sdk.ParseCoinsNormalized(additionalRewardStr)
				if err != nil {
					return err
				}
			}
			msg := types.MsgAdjustPool{
				PoolId:           args[0],
				AdditionalReward: additionalReward,
				RewardPerBlock:   rewardPerBlock,
				Creator:          clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAdjustFarmPool)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdStake implements the staking lp token to farm pool command.
func GetCmdStake() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stake",
		Short:   "Stake some lp token to farm pool",
		Example: fmt.Sprintf("$ %s tx farm stake <Farm Pool ID> <lp token> [flags]", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgStake{
				PoolId: args[0],
				Amount: amount,
				Sender: clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdUnstake implements the unstaking some lp token from farm pool command.
func GetCmdUnstake() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unstake",
		Short:   "Unstake some lp token from farm pool",
		Example: fmt.Sprintf("$ %s tx farm unstake <Farm Pool ID> <lp token> [flags]", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgUnstake{
				PoolId: args[0],
				Amount: amount,
				Sender: clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdHarvest implements the withdrawing some reward from the farm pool.
func GetCmdHarvest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "harvest",
		Short:   "withdraw some reward from the farm pool",
		Example: fmt.Sprintf("$ %s tx farm harvest <Farm Pool ID>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgHarvest{
				PoolId: args[0],
				Sender: clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
