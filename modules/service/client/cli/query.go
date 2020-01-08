package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	serviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the service module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryServiceDefinition(queryRoute, cdc),
		GetCmdQuerySvcBind(queryRoute, cdc),
		GetCmdQuerySvcBinds(queryRoute, cdc),
		GetCmdQuerySvcRequests(queryRoute, cdc),
		GetCmdQuerySvcResponse(queryRoute, cdc),
		GetCmdQuerySvcFees(queryRoute, cdc),
	)...)

	return serviceQueryCmd
}

// GetCmdQueryServiceDefinition implements the query service definition command.
func GetCmdQueryServiceDefinition(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "Query service definition",
		Example: "iriscli query service definition <service name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryDefinitionParams{
				ServiceName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDefinition)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var svcDef types.ServiceDefinition
			if err := cdc.UnmarshalJSON(res, &svcDef); err != nil {
				return err
			}

			return cliCtx.PrintOutput(svcDef)
		},
	}

	return cmd
}

// GetCmdQuerySvcBind implements the query service binding command.
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

			params := types.QueryBindingParams{
				DefChainID:  args[0],
				ServiceName: args[1],
				BindChainID: args[2],
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBinding)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var binding types.SvcBinding
			if err := cdc.UnmarshalJSON(res, &binding); err != nil {
				return err
			}

			return cliCtx.PrintOutput(binding)
		},
	}

	return cmd
}

// GetCmdQuerySvcBinds implements the query service bindings command.
func GetCmdQuerySvcBinds(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iriscli query service bindings <def-chain-id> <service name>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryBindingsParams{
				DefChainID:  args[0],
				ServiceName: args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBindings)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var bindings []types.SvcBinding
			if err := cdc.UnmarshalJSON(res, &bindings); err != nil {
				return err
			}

			return cliCtx.PrintOutput(bindings)
		},
	}

	return cmd
}

// GetCmdQuerySvcRequests implements the query service requests command.
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

			params := types.QueryBindingParams{
				DefChainID:  args[0],
				ServiceName: args[1],
				BindChainID: args[2],
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequests)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var requests []types.SvcRequest
			if err := cdc.UnmarshalJSON(res, &requests); err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	return cmd
}

// GetCmdQuerySvcResponse implements the query service response command.
func GetCmdQuerySvcResponse(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query a service response",
		Example: "iriscli query service response <req-chain-id> <request-id>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryResponseParams{
				ReqChainID: args[0],
				RequestID:  args[1],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryResponse)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var response types.SvcResponse
			if err := cdc.UnmarshalJSON(res, &response); err != nil {
				return err
			}

			return cliCtx.PrintOutput(response)
		},
	}

	return cmd
}

// GetCmdQuerySvcFees implements the query returned and incoming fee command.
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

			params := types.QueryFeesParams{
				Address: addr,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryFees)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var fees types.FeesOutput
			if err := cdc.UnmarshalJSON(res, &fees); err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	return cmd
}
