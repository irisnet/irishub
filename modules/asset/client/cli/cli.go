package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	token "github.com/irisnet/irishub/modules/asset/01-token"
	"github.com/irisnet/irishub/modules/asset/types"
)

// GetTxCmd returns the transaction commands for asset module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Asset transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(flags.PostCommands(
		token.GetTxCmd(cdc, storeKey),
	)...)
	return txCmd
}

// GetQueryCmd returns the query commands for asset module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the asset module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(flags.GetCommands(
		token.GetQueryCmd(cdc, queryRoute),
	)...)

	return queryCmd
}
