package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

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
	randQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryRandom(clientCtx),
		GetCmdQueryRandomRequestQueue(clientCtx),
	)...)
	return randQueryCmd
}

// GetCmdQueryRandom implements the query rand command.
func GetCmdQueryRandom(clientCtx client.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rand",
		Short:   "Query a random number by the request id",
		Example: "iriscli query rand rand <request id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			if err := types.CheckReqID(args[0]); err != nil {
				return err
			}

			params := types.QueryRandomParams{
				ReqID: args[0],
			}

			bz, err := clientCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandom)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var rawRandom types.Random
			if err := clientCtx.Codec.UnmarshalJSON(res, &rawRandom); err != nil {
				return err
			}

			readableRandom := types.ReadableRandom{
				RequestTxHash: hex.EncodeToString(rawRandom.RequestTxHash),
				Height:        rawRandom.Height,
				Value:         rawRandom.Value,
			}

			return clientCtx.PrintOutput(readableRandom)
		},
	}

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

			params := types.QueryRandomRequestQueueParams{
				Height: height,
			}

			bz, err := clientCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRandomRequestQueue)
			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var requests types.Requests
			if err := clientCtx.Codec.UnmarshalJSON(res, &requests); err != nil {
				return err
			}

			return clientCtx.PrintOutput(requests)
		},
	}

	return cmd
}
