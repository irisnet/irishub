package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// GetQueryCmd returns the cli query commands for the HTLC module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the HTLC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	mintingQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryHTLC(cdc),
		)...,
	)
	return mintingQueryCmd
}

// GetCmdQueryHTLC implements the query HTLC command.
func GetCmdQueryHTLC(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "htlc",
		Short:   "Query details of an HTLC",
		Example: "iriscli query htlc htlc <hash-lock>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			hashLock, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			params := types.QueryHTLCParams{
				HashLock: hashLock,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryHTLC)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var htlc types.HTLC
			if err := cdc.UnmarshalJSON(res, &htlc); err != nil {
				return err
			}

			return cliCtx.PrintOutput(htlc)
		},
	}

	return cmd
}
