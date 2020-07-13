package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/irisnet/irishub/modules/guardian/types"
)

// GetQueryCmd returns the cli query commands for the guardian module.
func GetQueryCmd(clientCtx client.Context) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the guardian module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(flags.GetCommands(
		GetCmdQueryProfilers(clientCtx),
		GetCmdQueryTrustees(clientCtx),
	)...)
	return txCmd
}

// GetCmdQueryProfilers implements the query profilers command.
func GetCmdQueryProfilers(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profilers",
		Short:   "Query for all profilers",
		Example: "iriscli query guardian profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProfilers), nil)
			if err != nil {
				return err
			}

			var profilers types.Profilers
			if err := clientCtx.Codec.UnmarshalJSON(res, &profilers); err != nil {
				return err
			}

			return clientCtx.PrintOutput(profilers)
		},
	}
	return cmd
}

// GetCmdQueryTrustees implements the query trustees command.
func GetCmdQueryTrustees(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trustees",
		Short:   "Query for all trustees",
		Example: "iriscli query guardian trustees",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTrustees), nil)
			if err != nil {
				return err
			}

			var trustees types.Trustees
			if err := clientCtx.Codec.UnmarshalJSON(res, &trustees); err != nil {
				return err
			}

			return clientCtx.PrintOutput(trustees)
		},
	}
	return cmd
}
