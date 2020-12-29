package cli

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/htlc/types"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd() *cobra.Command {
	htlcQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the HTLC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	htlcQueryCmd.AddCommand(
		GetCmdQueryHTLC(),
	)

	return htlcQueryCmd
}

// GetCmdQueryHTLC implements the query HTLC command.
func GetCmdQueryHTLC() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "htlc [hash-lock]",
		Short:   "Query an HTLC",
		Long:    "Query details of an HTLC with the specified hash lock.",
		Example: fmt.Sprintf("$ %s query htlc htlc <hash-lock>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			hashLock, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			param := types.QueryHTLCRequest{HashLock: tmbytes.HexBytes(hashLock).String()}
			response, err := queryClient.HTLC(context.Background(), &param)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(response.Htlc)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
