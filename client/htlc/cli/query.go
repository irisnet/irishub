package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v2/htlc"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/htlc/types"
	"github.com/irisnet/irishub/codec"
)

// GetCmdQueryHTLC implements the query HTLC command.
func GetCmdQueryHTLC(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-htlc [hash-lock]",
		Short:   "Query details of an HTLC",
		Example: "iriscli htlc query-htlc <hash-lock>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			hashLock, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			params := htlc.QueryHTLCParams{
				HashLock: hashLock,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.HtlcRoute, htlc.QueryHTLC), bz)
			if err != nil {
				return err
			}

			var htlc htlc.HTLC
			if err = cdc.UnmarshalJSON(res, &htlc); err != nil {
				return err
			}

			oh := types.NewOutputHTLC(htlc)

			return cliCtx.PrintOutput(oh)
		},
	}

	return cmd
}
