package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/irisnet/irishub/modules/rand/client/types"
	"github.com/irisnet/irishub/modules/rand/internal/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group rand queries under a subcommand
	randQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the rand module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	randQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryRand(queryRoute, cdc),
		GetCmdQueryRandRequestQueue(queryRoute, cdc),
	)...)

	return randQueryCmd
}

// GetCmdQueryRand implements the query rand command.
func GetCmdQueryRand(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rand",
		Short:   "Query a random number by the request id",
		Example: "iriscli query rand rand <request id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if err := types.CheckReqID(args[0]); err != nil {
				return err
			}

			params := types.QueryRandParams{
				ReqID: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRand)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var rawRand types.Rand
			err = cdc.UnmarshalJSON(res, &rawRand)
			if err != nil {
				return err
			}

			readableRand := clienttypes.ReadableRand{
				RequestTxHash: hex.EncodeToString(rawRand.RequestTxHash),
				Height:        rawRand.Height,
				Value:         rawRand.Value.FloatString(types.RandPrec),
			}

			return cliCtx.PrintOutput(readableRand)
		},
	}

	return cmd
}

// GetCmdQueryRandRequestQueue implements the query queue command.
func GetCmdQueryRandRequestQueue(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-queue",
		Short:   "Query the random number request queue with an optional height",
		Example: "iriscli query rand queue [gen-height]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var height int64
			var err error

			if len(args) > 0 {
				height, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
			}

			if height < 0 {
				return fmt.Errorf("the height must not be less than 0: %d", height)
			}

			params := types.QueryRandRequestQueueParams{
				Height: height,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRandRequestQueue)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var requests types.Requests
			err = cdc.UnmarshalJSON(res, &requests)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	return cmd
}
