package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group service queries under a subcommand
	serviceQueryCmd := &cobra.Command{
		Use:                        service.ModuleName,
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

func GetCmdQuerySvcDef(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "Query service definition",
		Example: "iriscli query service definition <def-chain-id> <service name>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QueryDefinitionParams{
				DefChainID:  args[0],
				ServiceName: args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, service.QueryDefinition)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return cmd
}

func GetCmdQuerySvcBind(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "Query service binding",
		Example: "iriscli query service binding <def-chain-id> <service name> <bind-chain-id> <provider>",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			params := service.QueryBindingParams{
				DefChainID:  args[0],
				ServiceName: args[1],
				BindChainId: args[2],
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, service.QueryBinding)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return cmd
}

func GetCmdQuerySvcBinds(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iriscli query service bindings <def-chain-id> <service name>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QueryBindingsParams{
				DefChainID:  args[0],
				ServiceName: args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, service.QueryBindings)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return cmd
}

func GetCmdQuerySvcRequests(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "requests",
		Short:   "Query service requests",
		Example: "iriscli query service requests <def-chain-id> <service-name> <bind-chain-id> <provider>",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			params := service.QueryBindingParams{
				DefChainID:  args[0],
				ServiceName: args[1],
				BindChainId: args[2],
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, service.QueryRequests)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return cmd
}

func GetCmdQuerySvcResponse(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query a service response",
		Example: "iriscli query service response <req-chain-id> <request-id>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QueryResponseParams{
				ReqChainId: args[0],
				RequestId:  args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, service.QueryResponse)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return cmd
}

func GetCmdQuerySvcFees(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees",
		Short:   "Query returned and incoming fee of a particular address",
		Example: "iriscli query service fees <account address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := service.QueryFeesParams{
				Address: addr,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, service.QueryFees)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(res)
		},
	}

	return cmd
}
