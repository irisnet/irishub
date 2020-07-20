package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

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
	txCmd.AddCommand(
		GetCmdQueryProfilers(clientCtx),
		GetCmdQueryTrustees(clientCtx),
	)
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Profilers(context.Background(), &types.QueryProfilersRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Profilers)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
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

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Trustees(context.Background(), &types.QueryTrusteesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Trustees)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
