package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"mods.irisnet.org/modules/random/types"
)

// GetQueryCmd returns the cli query commands for the random module.
func GetQueryCmd() *cobra.Command {
	// Group random queries under a subcommand
	randQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the random module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	randQueryCmd.AddCommand(
		GetCmdQueryRandom(),
		GetCmdQueryRandomRequestQueue(),
	)
	return randQueryCmd
}

// GetCmdQueryRandom implements the query random command.
func GetCmdQueryRandom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "random [request-id]",
		Short:   "Query a random number by the request id",
		Example: fmt.Sprintf("%s query random random <request-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.CheckReqID(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Random(context.Background(), &types.QueryRandomRequest{ReqId: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Random)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryRandomRequestQueue implements the query queue command.
func GetCmdQueryRandomRequestQueue() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "queue [gen-height]",
		Short:   "Query the random number request queue with an optional height",
		Example: fmt.Sprintf("%s query random queue <gen-height>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			var height int64

			if len(args) > 0 {
				if height, err = strconv.ParseInt(args[0], 10, 64); err != nil {
					return err
				}
			}

			if height < 0 {
				return fmt.Errorf("the height must not be less than 0: %d", height)
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.RandomRequestQueue(context.Background(), &types.QueryRandomRequestQueueRequest{Height: height})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
