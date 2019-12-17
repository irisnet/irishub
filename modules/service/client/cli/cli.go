package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group service queries under a subcommand
	serviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the service module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceQueryCmd.AddCommand(client.GetCommands(
		GetCmdQuerySvcDef(queryRoute, cdc),
		GetCmdQuerySvcBind(queryRoute, cdc),
		GetCmdQuerySvcBinds(queryRoute, cdc),
		GetCmdQuerySvcRequests(queryRoute, cdc),
		GetCmdQuerySvcResponse(queryRoute, cdc),
		GetCmdQuerySvcFees(queryRoute, cdc))...)

	return serviceQueryCmd
}


// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	serviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Service transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceTxCmd.AddCommand(client.PostCommands(
		GetCmdSvcDef(cdc),
		GetCmdSvcBind(cdc),
		GetCmdSvcBindUpdate(cdc),
		GetCmdSvcDisable(cdc),
		GetCmdSvcEnable(cdc),
		GetCmdSvcRefundDeposit(cdc),
		GetCmdSvcCall(cdc),
		GetCmdSvcRespond(cdc),
		GetCmdSvcRefundFees(cdc),
		GetCmdSvcWithdrawFees(cdc),
		GetCmdSvcWithdrawTax(cdc))...)

	return serviceTxCmd
}