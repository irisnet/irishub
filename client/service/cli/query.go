package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

func GetCmdQuerySvcBind(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "Query service binding",
		Example: "iriscli service binding --def-chain-id=<chain-id> --service-name=<service name> --bind-chain-id=<chain-id> --provider=<provider>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

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
				BindChainID: bindChainId,
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

			fmt.Println(string(res))
			return nil
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

func GetCmdQuerySvcBinds(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iriscli service bindings --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)

			params := service.QueryBindingsParams{
				DefChainID:  defChainID,
				ServiceName: name,
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

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	return cmd
}

func GetCmdQuerySvcRequests(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "requests",
		Short: "Query service requests",
		Example: "iriscli service requests --def-chain-id=<service-def-chain-id> --service-name=test " +
			"--bind-chain-id=<bind-chain-id> --provider=<provider>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)
			bindChainID := viper.GetString(FlagBindChainID)
			providerStr := viper.GetString(FlagProvider)

			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}

			params := service.QueryBindingParams{
				DefChainID:  defChainID,
				ServiceName: name,
				BindChainID: bindChainID,
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

			fmt.Println(string(res))
			return nil
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

func GetCmdQuerySvcResponse(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query a service response",
		Example: "iriscli service response --request-chain-id=<req-chain-id> --request-id=<request-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			reqChainID := viper.GetString(FlagReqChainID)
			reqID := viper.GetString(FlagReqID)

			params := service.QueryResponseParams{
				ReqChainID: reqChainID,
				RequestID:  reqID,
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

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().String(FlagReqChainID, "", "the ID of the blockchain that the service invocation initiated")
	cmd.Flags().String(FlagReqID, "", "the ID of the service invocation")
	_ = cmd.MarkFlagRequired(FlagReqChainID)
	_ = cmd.MarkFlagRequired(FlagReqID)
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
