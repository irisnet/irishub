package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/service/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// GetCmdQueryServiceDefinition implements the query service definition command
func GetCmdQueryServiceDefinition(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition [service-name]",
		Short:   "Query a service definition",
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
		Use:     "binding [service-name] [provider]",
		Short:   "Query a service binding",
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
		Use:     "bindings [service-name]",
		Short:   "Query all bindings of a service definition",
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
		Use:     "withdraw-addr [provider]",
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

// GetCmdQueryServiceRequest implements the query service request command
func GetCmdQueryServiceRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request [request-id]",
		Short:   "Query a request by the request ID",
		Example: "iriscli service request <request-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestID, err := service.ConvertRequestID(args[0])
			if err != nil {
				return err
			}
			params := service.QueryRequestParams{
				RequestID: requestID,
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
			_ = cdc.UnmarshalJSON(res, &request)
			if request.Empty() {
				request, err = utils.QueryRequestByTxQuery(cliCtx, params)
				if err != nil {
					return err
				}
			}

			if request.Empty() {
				return fmt.Errorf("unknown request: %s", params.RequestID)
			}

			return cliCtx.PrintOutput(request)
		},
	}

	return cmd
}

// GetCmdQueryServiceRequests implements the query service requests command
func GetCmdQueryServiceRequests(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "requests [service-name] [provider] | [request-context-id] [batch-counter]",
		Short:   "Query active requests by the service binding or request context ID",
		Example: "iriscli service requests <service-name> <provider> | <request-context-id> <batch-counter>",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryByBinding := true

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				queryByBinding = false
			}

			var requests service.Requests

			if queryByBinding {
				requests, err = utils.QueryRequestsByBinding(cliCtx, args[0], provider)
			} else {
				requests, err = utils.QueryRequestsByReqCtx(cliCtx, args[0], args[1])
			}

			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	return cmd
}

// GetCmdQueryServiceResponse implements the query service response command
func GetCmdQueryServiceResponse(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response [request-id]",
		Short:   "Query a response by the request ID",
		Example: "iriscli service response <request-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestID, err := service.ConvertRequestID(args[0])
			if err != nil {
				return err
			}
			params := service.QueryResponseParams{
				RequestID: requestID,
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
			_ = cdc.UnmarshalJSON(res, &response)
			if response.Empty() {
				response, err = utils.QueryResponseByTxQuery(cliCtx, params)
				if err != nil {
					return err
				}
			}

			if response.Empty() {
				return fmt.Errorf("unknown response: %s", params.RequestID)
			}

			return cliCtx.PrintOutput(response)
		},
	}

	return cmd
}

// GetCmdQueryServiceResponses implements the query service responses command
func GetCmdQueryServiceResponses(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "responses [request-context-id] [batch-counter]",
		Short:   "Query active responses by the request context ID and batch counter",
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

// GetCmdQueryRequestContext implements the query request context command
func GetCmdQueryRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request-context [request-context-id]",
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

			requestContext, err := utils.QueryRequestContext(cliCtx, params)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requestContext)
		},
	}

	return cmd
}

// GetCmdQueryEarnedFees implements the query earned fees command
func GetCmdQueryEarnedFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees [provider]",
		Short:   "Query the earned fees of a provider",
		Example: "iriscli service fees <provider>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := service.QueryEarnedFeesParams{
				Provider: provider,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryEarnedFees)
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

// GetCmdQuerySchema implements the query schema command
func GetCmdQuerySchema(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema [schema-name]",
		Short:   "Query the system schema by the schema name, only pricing and result allowed",
		Example: "iriscli service schema <schema-name>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := service.QuerySchemaParams{
				SchemaName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QuerySchema)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var schema service.SchemaType
			if err := cdc.UnmarshalJSON(res, &schema); err != nil {
				return err
			}

			return cliCtx.PrintOutput(schema)
		},
	}

	return cmd
}
