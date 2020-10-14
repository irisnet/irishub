package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/irisnet/irismod/modules/record/types"
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			fromAddr := clientCtx.GetFromAddress()
			uri, _ := cmd.Flags().GetString(FlagURI)
			meta, _ := cmd.Flags().GetString(FlagMeta)

			content := types.Content{
				Digest:     args[0],
				DigestAlgo: args[1],
				URI:        uri,
				Meta:       meta,
			}

			msg := types.NewMsgCreateRecord([]types.Content{content}, fromAddr)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreateRecord)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
