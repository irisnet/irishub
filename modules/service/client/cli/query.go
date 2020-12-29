package cli

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/service/client/utils"
	"github.com/irisnet/irismod/modules/service/types"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd() *cobra.Command {
	serviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the service module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceQueryCmd.AddCommand(
		GetCmdQueryServiceDefinition(),
		GetCmdQueryServiceBinding(),
		GetCmdQueryServiceBindings(),
		GetCmdQueryWithdrawAddr(),
		GetCmdQueryServiceRequest(),
		GetCmdQueryServiceRequests(),
		GetCmdQueryServiceResponse(),
		GetCmdQueryRequestContext(),
		GetCmdQueryServiceResponses(),
		GetCmdQueryEarnedFees(),
		GetCmdQuerySchema(),
		GetCmdQueryParams(),
	)

	return serviceQueryCmd
}

// GetCmdQueryServiceDefinition implements the query service definition command.
func GetCmdQueryServiceDefinition() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition [service-name]",
		Short:   "Query a service definition",
		Long:    "Query details of a service definition.",
		Example: fmt.Sprintf("$ %s query service definition <service-name>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateServiceName(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Definition(
				context.Background(),
				&types.QueryDefinitionRequest{
					ServiceName: args[0],
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.ServiceDefinition)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryServiceBinding implements the query service binding command
func GetCmdQueryServiceBinding() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding [service-name] [provider-address]",
		Short:   "Query a service binding",
		Long:    "Query details of a service binding.",
		Example: fmt.Sprintf("$ %s query service binding <service-name> <provider-address>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateServiceName(args[0]); err != nil {
				return err
			}

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Binding(
				context.Background(),
				&types.QueryBindingRequest{
					ServiceName: args[0],
					Provider:    provider.String(),
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.ServiceBinding)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryServiceBindings implements the query service bindings command
func GetCmdQueryServiceBindings() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings [service-name]",
		Short:   "Query all bindings of a service definition with an optional owner",
		Long:    "Query all bindings of a service definition with an optional owner.",
		Example: fmt.Sprintf("$ %s query service bindings <service-name> --owner=<address>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateServiceName(args[0]); err != nil {
				return err
			}

			var owner sdk.AccAddress
			ownerStr, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}
			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Bindings(
				context.Background(),
				&types.QueryBindingsRequest{
					ServiceName: args[0],
					Owner:       owner.String(),
					Pagination:  pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(FsQueryServiceBindings)
	flags.AddPaginationFlagsToCmd(cmd, "all bindings")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryWithdrawAddr implements the query withdraw address command
func GetCmdQueryWithdrawAddr() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-addr [address]",
		Short:   "Query the withdrawal address of an owner",
		Long:    "Query the withdrawal address of an owner.",
		Example: fmt.Sprintf("$ %s query service withdraw-addr <address>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.WithdrawAddress(
				context.Background(),
				&types.QueryWithdrawAddressRequest{
					Owner: owner.String(),
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryServiceRequest implements the query service request command
func GetCmdQueryServiceRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request [request-id]",
		Short:   "Query a request by the request ID",
		Long:    "Query details of a service request.",
		Example: fmt.Sprintf("$ %s query service request <request-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			requestID, err := types.ConvertRequestID(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Request(
				context.Background(),
				&types.QueryRequestRequest{
					RequestId: requestID.String(),
				},
			)
			if err != nil {
				return err
			}

			if res.Request.Empty() {
				request, err := utils.QueryRequestByTxQuery(clientCtx, types.QuerierRoute, requestID)
				if err != nil {
					return err
				}
				res.Request = &request
			}

			if res.Request.Empty() {
				return fmt.Errorf("unknown request: %s", requestID)
			}

			return clientCtx.PrintProto(res.Request)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryServiceRequests implements the query service requests command
func GetCmdQueryServiceRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "requests [service-name] [provider] | [request-context-id] [batch-counter]",
		Short:   "Query active requests by the service binding or request context ID",
		Long:    "Query active requests by the service binding or request context ID.",
		Example: fmt.Sprintf("$ %s query service requests <service-name> <provider> | <request-context-id> <batch-counter>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryByBinding := true
			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				queryByBinding = false
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			if queryByBinding {
				res, err := queryClient.Requests(context.Background(), &types.QueryRequestsRequest{
					ServiceName: args[0],
					Provider:    provider.String(),
					Pagination:  pageReq,
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			batchCounter, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.RequestsByReqCtx(
				context.Background(),
				&types.QueryRequestsByReqCtxRequest{
					RequestContextId: tmbytes.HexBytes(requestContextID).String(),
					BatchCounter:     batchCounter,
					Pagination:       pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all requests")
	return cmd
}

// GetCmdQueryServiceResponse implements the query service response command
func GetCmdQueryServiceResponse() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response [request-id]",
		Short:   "Query a response by the request ID",
		Long:    "Query details of a service response.",
		Example: fmt.Sprintf("$ %s query service response <request-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			requestID, err := types.ConvertRequestID(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Response(
				context.Background(),
				&types.QueryResponseRequest{
					RequestId: requestID.String(),
				},
			)
			if err != nil {
				return err
			}

			if res.Response.Empty() {
				response, err := utils.QueryResponseByTxQuery(clientCtx, types.QuerierRoute, requestID)
				if err != nil {
					return err
				}
				res.Response = &response
			}

			if res.Response.Empty() {
				return fmt.Errorf("unknown response: %s", requestID)
			}

			return clientCtx.PrintProto(res.Response)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryServiceResponses implements the query service responses command
func GetCmdQueryServiceResponses() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "responses [request-context-id] [batch-counter]",
		Short:   "Query active responses by the request context ID and batch counter",
		Long:    "Query active responses by the request context ID and batch counter.",
		Example: fmt.Sprintf("$ %s query service responses <request-context-id> <batch-counter>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			batchCounter, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Responses(
				context.Background(),
				&types.QueryResponsesRequest{
					RequestContextId: tmbytes.HexBytes(requestContextID).String(),
					BatchCounter:     batchCounter,
					Pagination:       pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all responses")

	return cmd
}

// GetCmdQueryRequestContext implements the query request context command
func GetCmdQueryRequestContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request-context [request-context-id]",
		Short:   "Query a request context",
		Long:    "Query a request context.",
		Example: fmt.Sprintf("$ %s query service request-context <request-context-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			res, err := utils.QueryRequestContext(
				clientCtx,
				types.QuerierRoute,
				types.QueryRequestContextRequest{
					RequestContextId: tmbytes.HexBytes(requestContextID).String(),
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryEarnedFees implements the query earned fees command
func GetCmdQueryEarnedFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees [provider-address]",
		Short:   "Query the earned fees of a provider",
		Long:    "Query the earned fees of a provider.",
		Example: fmt.Sprintf("$ %s query service fees <provider-address>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.EarnedFees(
				context.Background(),
				&types.QueryEarnedFeesRequest{
					Provider: provider.String(),
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQuerySchema implements the query schema command
func GetCmdQuerySchema() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "schema [schema-name]",
		Short:   "Query the system schema by the schema name",
		Long:    "Query the system schema by the schema name, only pricing and result allowed.",
		Example: fmt.Sprintf("$ %s query service schema <schema-name>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Schema(
				context.Background(),
				&types.QuerySchemaRequest{
					SchemaName: args[0],
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "params",
		Short:   "Query the current service parameter values",
		Long:    "Query values set as service parameters.",
		Example: fmt.Sprintf("$ %s query service params", version.AppName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
