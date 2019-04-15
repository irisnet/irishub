package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/service"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/app/protocol"
)

func GetCmdQuerySvcDef(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "Query service definition",
		Example: "iriscli service definition --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

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

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryDefinition)
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
				BindChainId: bindChainId,
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
			defChainId := viper.GetString(FlagDefChainID)

			params := service.QueryServiceParams{
				DefChainID:  defChainId,
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

			route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryResponse)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().String(FlagReqChainId, "", "the ID of the blockchain that the service invocation initiated")
	cmd.Flags().String(FlagReqId, "", "the ID of the service invocation")
	cmd.MarkFlagRequired(FlagReqChainId)
	cmd.MarkFlagRequired(FlagReqId)
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
