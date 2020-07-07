package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/irisnet/irishub/modules/mint/types"
)

// GetQueryCmd returns the cli query commands for the mint module.
func GetQueryCmd(clientCtx client.Context) *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the minting module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	mintingQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryParams(clientCtx),
		)...,
	)
	return mintingQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting parameters.
func GetCmdQueryParams(clientCtx client.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			res, _, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := clientCtx.Codec.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return clientCtx.PrintOutput(params)
		},
	}
}
