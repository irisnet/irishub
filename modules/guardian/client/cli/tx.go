package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/irisnet/irishub/modules/guardian/types"
)

// GetTxCmd returns the transaction commands for the guardian module.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "guardian transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdCreateProfiler(),
		GetCmdDeleteProfiler(),
		GetCmdCreateTrustee(),
		GetCmdDeleteTrustee(),
	)
	return txCmd
}

// GetCmdCreateProfiler implements the create profiler command.
func GetCmdCreateProfiler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-profiler",
		Short: "Add a new profiler",
		Example: fmt.Sprintf("%s tx guardian add-profiler --chain-id=<chain-id> --from=<key-name> --fees=0.3iris "+
			"--address=<added address> --description=<name>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			fromAddr := clientCtx.GetFromAddress()

			paStr := viper.GetString(FlagAddress)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			description := viper.GetString(FlagDescription)
			msg := types.NewMsgAddProfiler(description, pAddr, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	_ = cmd.MarkFlagRequired(FlagAddress)
	_ = cmd.MarkFlagRequired(FlagDescription)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDeleteProfiler implements the delete profiler command.
func GetCmdDeleteProfiler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-profiler",
		Short: "Delete a profiler",
		Example: fmt.Sprintf("%s tx guardian delete-profiler --chain-id=<chain-id> --from=<key-name> --fees=0.3iris "+
			"--address=<deleted address>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			fromAddr := clientCtx.GetFromAddress()
			paStr := viper.GetString(FlagAddress)
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteProfiler(pAddr, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	_ = cmd.MarkFlagRequired(FlagAddress)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdCreateTrustee implements the create trustee command.
func GetCmdCreateTrustee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-trustee",
		Short: "Add a new trustee",
		Example: fmt.Sprintf("%s tx guardian add-trustee --chain-id=<chain-id> --from=<key-name> --fees=0.3iris "+
			"--address=<added address> --description=<name>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			fromAddr := clientCtx.GetFromAddress()
			taStr := viper.GetString(FlagAddress)
			if len(taStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			tAddr, err := sdk.AccAddressFromBech32(taStr)
			if err != nil {
				return err
			}
			description := viper.GetString(FlagDescription)
			msg := types.NewMsgAddTrustee(description, tAddr, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	_ = cmd.MarkFlagRequired(FlagDescription)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdDeleteTrustee implements the delete trustee command.
func GetCmdDeleteTrustee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-trustee",
		Short: "Delete a trustee",
		Example: fmt.Sprintf("%s tx guardian delete-trustee --chain-id=<chain-id> --from=<key-name> --fees=0.3iris "+
			"--address=<deleted address>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			fromAddr := clientCtx.GetFromAddress()
			taStr := viper.GetString(FlagAddress)
			tAddr, err := sdk.AccAddressFromBech32(taStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteTrustee(tAddr, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	_ = cmd.MarkFlagRequired(FlagAddress)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
