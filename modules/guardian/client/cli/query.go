package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

// GetQueryCmd returns the cli query commands for the guardian module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the guardian module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(flags.GetCommands(
		GetCmdQueryProfilers(cdc),
		GetCmdQueryTrustees(cdc),
	)...)
	return txCmd
}

// GetCmdQueryProfilers implements the query profilers command.
func GetCmdQueryProfilers(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profilers",
		Short:   "Query for all profilers",
		Example: "iriscli query guardian profilers",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProfilers), nil)
			if err != nil {
				return err
			}

			var profilers types.Profilers
			if err := cdc.UnmarshalJSON(res, &profilers); err != nil {
				return err
			}

			return cliCtx.PrintOutput(profilers)
		},
	}
	return cmd
}

// GetCmdQueryTrustees implements the query trustees command.
func GetCmdQueryTrustees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trustees",
		Short:   "Query for all trustees",
		Example: "iriscli query guardian trustees",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTrustees), nil)
			if err != nil {
				return err
			}

			var trustees types.Trustees
			if err := cdc.UnmarshalJSON(res, &trustees); err != nil {
				return err
			}

			return cliCtx.PrintOutput(trustees)
		},
	}
	return cmd
}
