package cli

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/modules/random/types"
)

// GetQueryCmd returns the cli query commands for the rand module.
func GetQueryCmd(clientCtx client.Context) *cobra.Command {
	// Group rand queries under a subcommand
	randQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the rand module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	randQueryCmd.AddCommand(
		GetCmdQueryRandom(clientCtx),
		GetCmdQueryRandomRequestQueue(clientCtx),
	)
	return randQueryCmd
}

// GetCmdQueryRandom implements the query rand command.
func GetCmdQueryRandom(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rand",
		Short:   "Query a random number by the request id",
		Example: fmt.Sprintf("%s query rand rand <request id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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

			readableRandom := types.ReadableRandom{
				RequestTxHash: hex.EncodeToString(res.Random.RequestTxHash),
				Height:        res.Random.Height,
				Value:         res.Random.Value,
			}

			return clientCtx.PrintOutput(readableRandom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryRandomRequestQueue implements the query queue command.
func GetCmdQueryRandomRequestQueue(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-queue",
		Short:   "Query the random number request queue with an optional height",
		Example: "iriscli query rand queue [gen-height]",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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

			return clientCtx.PrintOutput(res.Requests)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
