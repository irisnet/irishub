package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client/context"
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

// GetCmdQueryWithdrawAddr implements the query withdraw address command
func GetCmdQueryWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-addr",
		Short:   "Query the withdrawal address of a provider",
		Example: "iriscli service withdraw-addr <provider>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := service.QueryWithdrawAddressParams{
				Provider: provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryWithdrawAddress)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var withdrawAddr sdk.AccAddress
			if err := cdc.UnmarshalJSON(res, &withdrawAddr); err != nil {
				return err
			}

			return cliCtx.PrintOutput(withdrawAddr)
		},
	}

	return cmd
}

func GetCmdQueryServiceRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request",
		Short:   "Query a request by the request ID",
		Example: "iriscli service request <request-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QueryRequestParams{
				RequestID: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequest)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var request service.Request
			if err := cdc.UnmarshalJSON(res, &request); err != nil {
				return err
			}

			return cliCtx.PrintOutput(request)
		},
	}

	return cmd
}

func GetCmdQueryServiceRequests(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "requests",
		Short:   "Query service requests by the service binding or request context ID",
		Example: "iriscli service requests <service-name> <provider> | <request-context-id> <batch-counter>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryByBinding := true

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				queryByBinding = false
			}

			var params interface{}
			var route string

			if queryByBinding {
				params = service.QueryRequestsParams{
					ServiceName: args[0],
					Provider:    provider,
				}

				route = fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequests)
			} else {
				requestContextID, err := hex.DecodeString(args[0])
				if err != nil {
					return err
				}

				batchCounter, err := strconv.ParseUint(args[1], 10, 64)
				if err != nil {
					return err
				}

				params = service.QueryRequestsByReqCtxParams{
					RequestContextID: requestContextID,
					BatchCounter:     batchCounter,
				}

				route = fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequestsByReqCtx)
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

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
		Short:   "Query a response by the request ID",
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

func GetCmdQueryServiceResponses(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query responses by the request context ID and batch counter",
		Example: "iriscli service responses <request-context-id> <batch-counter>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			batchCounter, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			params := service.QueryResponsesParams{
				RequestContextID: requestContextID,
				BatchCounter:     batchCounter,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryResponses)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var responses service.Responses
			if err := cdc.UnmarshalJSON(res, &responses); err != nil {
				return err
			}

			return cliCtx.PrintOutput(responses)
		},
	}

	return cmd
}

func GetCmdQueryRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request-context",
		Short:   "Query a request context",
		Example: "iriscli service request-context <request-context-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			params := service.QueryRequestContextParams{
				RequestContextID: requestContextID,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequestContext)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var requestContext service.RequestContext
			if err := cdc.UnmarshalJSON(res, &requestContext); err != nil {
				return err
			}

			return cliCtx.PrintOutput(requestContext)
		},
	}

	return cmd
}

func GetCmdQueryEarnedFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees",
		Short:   "Query the earned fees",
		Example: "iriscli service fees <address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := service.QueryFeesParams{
				Address: provider,
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

			var fees service.EarnedFees
			if err := cdc.UnmarshalJSON(res, &fees); err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	return cmd
}
