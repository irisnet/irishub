package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
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

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	randTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Rand transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	randTxCmd.AddCommand(client.PostCommands(
		GetCmdRequestRand(cdc),
	)...)

	return randTxCmd
}
