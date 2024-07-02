package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"mods.irisnet.org/modules/record/types"
)

// NewTxCmd returns the transaction commands for the record module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Record transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdCreateRecord(),
	)
	return txCmd
}

// GetCmdCreateRecord implements the create record command.
func GetCmdCreateRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [digest] [digest-algo]",
		Short: "Create a new record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress().String()
			uri, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}
			meta, err := cmd.Flags().GetString(FlagMeta)
			if err != nil {
				return err
			}

			content := types.Content{
				Digest:     args[0],
				DigestAlgo: args[1],
				URI:        uri,
				Meta:       meta,
			}

			msg := types.NewMsgCreateRecord([]types.Content{content}, from)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreateRecord)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
