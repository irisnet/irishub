package cli

import (
	"os"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	cmn "github.com/irisnet/irishub/client/service"
)

const NULL = "null"

func GetCmdQuerySvcDef(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "Query service definition",
		Example: "iriscli service definition --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			res, err := cliCtx.QueryStore(service.GetServiceDefinitionKey(defChainId, name), storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("chain-id [%s] service [%s] is not existed", defChainId, name)
			}

			var svcDef service.SvcDef
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &svcDef)

			res1, err := cliCtx.QuerySubspace(service.GetMethodsSubspaceKey(defChainId, name), storeName)
			if err != nil {
				return err
			}

			var methods []service.MethodProperty
			for _, re := range res1 {
				var method service.MethodProperty
				cdc.MustUnmarshalBinaryLengthPrefixed(re.Value, &method)
				methods = append(methods, method)
			}

			output, err := codec.MarshalJSONIndent(cdc, cmn.DefOutput{Definition: svcDef, Methods: methods})
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)

	return cmd
}

func GetCmdQuerySvcBind(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "Query service binding",
		Example: "iriscli service binding --def-chain-id=<chain-id> --service-name=<service name> --bind-chain-id=<chain-id> --provider=<provider>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)
			bindChainId := viper.GetString(FlagBindChainID)
			providerStr := viper.GetString(FlagProvider)

			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryStore(service.GetServiceBindingKey(defChainId, name, bindChainId, provider), storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("def-chain-id [%s] service [%s] bind-chain-id [%s] provider [%s] is not existed", defChainId, name, bindChainId, provider)
			}

			var svcBinding service.SvcBinding
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &svcBinding)
			output, err := codec.MarshalJSONIndent(cdc, svcBinding)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsBindChainID)
	cmd.Flags().AddFlagSet(FsProvider)

	return cmd
}

func GetCmdQuerySvcBinds(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "Query service bindings",
		Example: "iriscli service bindings --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			res, err := cliCtx.QuerySubspace(service.GetBindingsSubspaceKey(defChainId, name), storeName)
			if err != nil {
				return err
			}

			var bindings []service.SvcBinding
			for _, re := range res {
				var binding service.SvcBinding
				cdc.MustUnmarshalBinaryLengthPrefixed(re.Value, &binding)
				bindings = append(bindings, binding)
			}

			output, err := codec.MarshalJSONIndent(cdc, bindings)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)

	return cmd
}

func GetCmdQuerySvcRequests(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "requests",
		Short: "Query service requests",
		Example: "iriscli service requests --def-chain-id=<service-def-chain-id> --service-name=test " +
			"--bind-chain-id=<bind-chain-id> --provider=<provider>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)
			bindChainId := viper.GetString(FlagBindChainID)
			providerStr := viper.GetString(FlagProvider)

			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}

			res, err := cliCtx.QuerySubspace(service.GetSubActiveRequestKey(defChainId, name, bindChainId, provider), storeName)
			if err != nil {
				return err
			}

			var reqs []service.SvcRequest
			for _, re := range res {
				var req service.SvcRequest
				cdc.MustUnmarshalBinaryLengthPrefixed(re.Value, &req)
				reqs = append(reqs, req)
			}

			output, err := codec.MarshalJSONIndent(cdc, reqs)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsBindChainID)
	cmd.Flags().AddFlagSet(FsProvider)

	return cmd
}

func GetCmdQuerySvcResponse(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Short:   "Query a service response",
		Example: "iriscli service response --request-chain-id=<req-chain-id> --request-id=<request-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			reqChainId := viper.GetString(FlagReqChainId)
			reqId := viper.GetString(FlagReqId)

			eHeight, rHeight, counter, err := service.ConvertRequestID(reqId)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(service.GetResponseKey(reqChainId, eHeight, rHeight, counter), storeName)
			var resp service.SvcResponse
			if err != nil {
				return err
			}
			if len(res) > 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res, &resp)
			} else {
				fmt.Println(NULL)
				return nil
			}
			output, err := codec.MarshalJSONIndent(cdc, resp)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsReqChainId)
	cmd.Flags().AddFlagSet(FsReqId)

	return cmd
}

func GetCmdQuerySvcFees(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fees",
		Short:   "Query return and incoming fee of a particular address",
		Example: "iriscli service fees <account address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			addrString := args[0]

			delAddr, err := sdk.AccAddressFromBech32(addrString)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryStore(service.GetReturnedFeeKey(delAddr), storeName)
			var returnedFee service.ReturnedFee
			if err != nil {
				return err
			}
			if len(res) > 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res, &returnedFee)
			}

			res1, err := cliCtx.QueryStore(service.GetIncomingFeeKey(delAddr), storeName)
			var incomingFee service.IncomingFee
			if err != nil {
				return err
			}
			if len(res1) > 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(res1, &incomingFee)
			}

			output, err := codec.MarshalJSONIndent(cdc, cmn.FeesOutput{ReturnedFee: returnedFee.Coins, IncomingFee: incomingFee.Coins})
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}
