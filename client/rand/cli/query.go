package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/rand"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/rand/types"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryRand implements the query-rand command.
func GetCmdQueryRand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-rand",
		Short:   "Query a random number by the request id",
		Example: "iriscli rand query-rand --request-id=<request id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			reqID := viper.GetString(FlagReqID)
			if err := rand.CheckReqID(reqID); err != nil {
				return err
			}

			params := rand.QueryRandParams{
				ReqID: reqID,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRand), bz)
			if err != nil {
				return err
			}

			var rawRand rand.Rand
			err = cdc.UnmarshalJSON(res, &rawRand)
			if err != nil {
				return err
			}

			readableRand := types.ReadableRand{
				RequestTxHash: hex.EncodeToString(rawRand.RequestTxHash),
				Height:        rawRand.Height,
				Value:         rawRand.Value.Rat.FloatString(rand.RandPrec),
			}

			return cliCtx.PrintOutput(readableRand)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryRand)
	cmd.MarkFlagRequired(FlagReqID)

	return cmd
}

// GetCmdQueryRandRequestQueue implements the query-queue command.
func GetCmdQueryRandRequestQueue(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-queue",
		Short:   "Query the random number request queue with an optional height",
		Example: "iriscli rand query-queue [--queue-height=<queue height>]",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			height := viper.GetInt64(FlagQueueHeight)
			if height < 0 {
				return fmt.Errorf("the height must not be less than 0: %d", height)
			}

			params := rand.QueryRandRequestQueueParams{
				Height: height,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.RandRoute, rand.QueryRandRequestQueue), bz)
			if err != nil {
				return err
			}

			var requests rand.Requests
			err = cdc.UnmarshalJSON(res, &requests)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryQueue)

	return cmd
}
