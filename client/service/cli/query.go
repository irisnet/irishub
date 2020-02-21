package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// GetCmdQueryServiceDefinition implements the query service definition command
func GetCmdQueryServiceDefinition(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "Query service definition",
		Example: "iriscli service definition <service name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QueryDefinitionParams{
				ServiceName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryDefinition)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var svcDef service.ServiceDefinition
			if err := cdc.UnmarshalJSON(res, &svcDef); err != nil {
				return err
			}

			return cliCtx.PrintOutput(svcDef)
		},
	}

	return cmd
}

// GetCmdQueryServiceBinding implements the query service binding command
func GetCmdQueryServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "Query service binding",
		Example: "iriscli service binding <service name> <provider>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := service.QueryBindingParams{
				ServiceName: args[0],
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryBinding)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var svcBinding service.ServiceBinding
			if err := cdc.UnmarshalJSON(res, &svcBinding); err != nil {
				return err
			}

			return cliCtx.PrintOutput(svcBinding)
		},
	}

	return cmd
}

// GetCmdQueryServiceBindings implements the query service bindings command
func GetCmdQueryServiceBindings(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iriscli service bindings <service name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QueryBindingsParams{
				ServiceName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryBindings)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var bindings service.ServiceBindings
			if err := cdc.UnmarshalJSON(res, &bindings); err != nil {
				return err
			}

			return cliCtx.PrintOutput(bindings)
		},
	}

	return cmd
}

func GetCmdQueryServiceRequests(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "requests",
		Short:   "Query service requests",
		Example: "iriscli service requests <service-name> <provider>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			serviceName := args[0]

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := service.QueryRequestsParams{
				ServiceName: serviceName,
				Provider:    provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequests)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var requests service.Requests
			if err := cdc.UnmarshalJSON(res, &requests); err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	return cmd
}

func GetCmdQueryServiceResponse(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query a service response",
		Example: "iriscli service response <request-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QueryResponseParams{
				RequestID: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryResponse)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var response service.Response
			if err := cdc.UnmarshalJSON(res, &response); err != nil {
				return err
			}

			return cliCtx.PrintOutput(response)
		},
	}

	return cmd
}

func GetCmdQuerySvcFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees",
		Short:   "Query return and incoming fee of a particular address",
		Example: "iriscli service fees <account address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			addrString := args[0]

			addr, err := sdk.AccAddressFromBech32(addrString)
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

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryFees)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}
