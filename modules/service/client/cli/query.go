package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Example: "iriscli service definition --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			params := service.QueryServiceParams{
				DefChainID:  defChainId,
				ServiceName: name,
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

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

func GetCmdQuerySvcBind(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "Query service binding",
		Example: "iriscli service binding --def-chain-id=<chain-id> --service-name=<service name> --bind-chain-id=<chain-id> --provider=<provider>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)
			bindChainId := viper.GetString(FlagBindChainID)
			providerStr := viper.GetString(FlagProvider)

			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}

			params := service.QueryBindingParams{
				DefChainID:  defChainId,
				ServiceName: name,
				BindChainId: bindChainId,
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

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBinding)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	cmd.MarkFlagRequired(FlagBindChainID)
	cmd.MarkFlagRequired(FlagProvider)

	return cmd
}

func GetCmdQuerySvcBinds(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iriscli service bindings --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			params := service.QueryServiceParams{
				DefChainID:  defChainId,
				ServiceName: name,
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

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

func GetCmdQuerySvcRequests(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "requests",
		Short: "Query service requests",
		Example: "iriscli service requests --def-chain-id=<service-def-chain-id> --service-name=test " +
			"--bind-chain-id=<bind-chain-id> --provider=<provider>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)
			bindChainId := viper.GetString(FlagBindChainID)
			providerStr := viper.GetString(FlagProvider)

			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}

			params := service.QueryBindingParams{
				DefChainID:  defChainId,
				ServiceName: name,
				BindChainId: bindChainId,
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

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBinding)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	cmd.MarkFlagRequired(FlagBindChainID)
	cmd.MarkFlagRequired(FlagProvider)

	return cmd
}

func GetCmdQuerySvcResponse(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query a service response",
		Example: "iriscli service response --request-chain-id=<req-chain-id> --request-id=<request-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			reqChainId := viper.GetString(FlagReqChainId)
			reqId := viper.GetString(FlagReqId)

			params := service.QueryResponseParams{
				ReqChainId: reqChainId,
				RequestId:  reqId,
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

	cmd.Flags().String(FlagReqChainId, "", "the ID of the blockchain that the service invocation initiated")
	cmd.Flags().String(FlagReqId, "", "the ID of the service invocation")
	cmd.MarkFlagRequired(FlagReqChainId)
	cmd.MarkFlagRequired(FlagReqId)

	return cmd
}

func GetCmdQuerySvcFees(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees",
		Short:   "Query return and incoming fee of a particular address",
		Example: "iriscli service fees <account address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addrStr := args[0]
			addr, err := sdk.AccAddressFromBech32(addrStr)
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
