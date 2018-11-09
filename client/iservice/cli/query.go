package cli

import (
	"os"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/iservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client/context"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	cmn "github.com/irisnet/irishub/client/iservice"
)

func GetCmdQueryScvDef(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "definition",
		Short:   "query service definition",
		Example: "iriscli iservice definition --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			res, err := cliCtx.QueryStore(iservice.GetServiceDefinitionKey(defChainId, name), storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("chain-id [%s] service [%s] is not existed", defChainId, name)
			}

			var svcDef iservice.SvcDef
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &svcDef)

			res2, err := cliCtx.QuerySubspace(iservice.GetMethodsSubspaceKey(defChainId, name), storeName)
			if err != nil {
				return err
			}

			var methods []iservice.MethodProperty
			for i := 0; i < len(res2); i++ {
				var method iservice.MethodProperty
				cdc.MustUnmarshalBinaryLengthPrefixed(res2[i].Value, &method)
				methods = append(methods, method)
			}

			service := cmn.ServiceOutput{SvcDef: svcDef, Methods: methods}
			output, err := codec.MarshalJSONIndent(cdc, service)
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

func GetCmdQueryScvBind(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "binding",
		Short:   "query service binding",
		Example: "iriscli iservice binding --def-chain-id=<chain-id> --service-name=<service name> --bind-chain-id=<chain-id> --provider=<provider>",
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
			res, err := cliCtx.QueryStore(iservice.GetServiceBindingKey(defChainId, name, bindChainId, provider), storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("def-chain-id [%s] service [%s] bind-chain-id [%s] provider [%s] is not existed", defChainId, name, bindChainId, provider)
			}

			var svcBinding iservice.SvcBinding
			cdc.MustUnmarshalBinary(res, &svcBinding)
			output, err := codec.MarshalJSONIndent(cdc, svcBinding)
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

func GetCmdQueryScvBinds(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bindings",
		Short:   "query service bindings",
		Example: "iriscli iservice bindings --def-chain-id=<chain-id> --service-name=<service name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			res, err := cliCtx.QuerySubspace(iservice.GetBindingsSubspaceKey(defChainId, name), storeName)
			if err != nil {
				return err
			}

			var bindings []iservice.SvcBinding
			for i := 0; i < len(res); i++ {
				var binding iservice.SvcBinding
				cdc.MustUnmarshalBinary(res[i].Value, &binding)
				bindings = append(bindings, binding)
			}

			output, err := cdc.MarshalJSONIndent(bindings, "", "")
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsServiceName)

	return cmd
}
