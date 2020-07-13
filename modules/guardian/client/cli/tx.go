package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	txCmd.AddCommand(flags.PostCommands(
		GetCmdCreateProfiler(),
		GetCmdDeleteProfiler(),
		GetCmdCreateTrustee(),
		GetCmdDeleteTrustee(),
	)...)
	return txCmd
}

// GetCmdCreateProfiler implements the create profiler command.
func GetCmdCreateProfiler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-profiler",
		Short: "Add a new profiler",
		Example: "iriscli tx guardian add-profiler --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<added address> --description=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

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
			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	_ = cmd.MarkFlagRequired(FlagAddress)
	_ = cmd.MarkFlagRequired(FlagDescription)
	return cmd
}

// GetCmdDeleteProfiler implements the delete profiler command.
func GetCmdDeleteProfiler() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-profiler",
		Short: "Delete a profiler",
		Example: "iriscli tx guardian delete-profiler --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<deleted address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			fromAddr := clientCtx.GetFromAddress()
			paStr := viper.GetString(FlagAddress)
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteProfiler(pAddr, fromAddr)
			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	_ = cmd.MarkFlagRequired(FlagAddress)
	return cmd
}

// GetCmdCreateTrustee implements the create trustee command.
func GetCmdCreateTrustee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-trustee",
		Short: "Add a new trustee",
		Example: "iriscli tx guardian add-trustee --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<added address> --description=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

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
			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	_ = cmd.MarkFlagRequired(FlagDescription)
	return cmd
}

// GetCmdDeleteTrustee implements the delete trustee command.
func GetCmdDeleteTrustee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-trustee",
		Short: "Delete a trustee",
		Example: "iriscli tx guardian delete-trustee --chain-id=<chain-id> --from=<key-name> --fees=0.3iris " +
			"--address=<deleted address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			fromAddr := clientCtx.GetFromAddress()
			taStr := viper.GetString(FlagAddress)
			tAddr, err := sdk.AccAddressFromBech32(taStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteTrustee(tAddr, fromAddr)
			return tx.GenerateOrBroadcastTx(clientCtx, msg)
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	_ = cmd.MarkFlagRequired(FlagAddress)
	return cmd
}
