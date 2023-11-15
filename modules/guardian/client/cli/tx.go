package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/v2/modules/guardian/types"
)

// NewTxCmd returns the transaction commands for the guardian module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "guardian transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdCreateSuper(),
		GetCmdDeleteSuper(),
	)
	return txCmd
}

// GetCmdCreateSuper implements the create super command.
func GetCmdCreateSuper() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-super",
		Short: "Add a new super",
		Example: fmt.Sprintf(
			"%s tx guardian add-super --chain-id=<chain-id> --from=<key-name> --fees=0.3iris --address=<added address> --description=<name>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromAddr := clientCtx.GetFromAddress()

			paStr, _ := cmd.Flags().GetString(FlagAddress)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			description, _ := cmd.Flags().GetString(FlagDescription)
			msg := types.NewMsgAddSuper(description, pAddr, fromAddr)
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

// GetCmdDeleteSuper implements the delete super command.
func GetCmdDeleteSuper() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-super",
		Short: "Delete a super",
		Example: fmt.Sprintf(
			"%s tx guardian delete-super --chain-id=<chain-id> --from=<key-name> --fees=0.3iris --address=<deleted address>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromAddr := clientCtx.GetFromAddress()
			paStr, _ := cmd.Flags().GetString(FlagAddress)
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgDeleteSuper(pAddr, fromAddr)
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
