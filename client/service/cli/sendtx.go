package cli

import (
	"os"
	"fmt"
	"strings"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	authcmd "github.com/irisnet/irishub/modules/auth/client/cli"
	"github.com/irisnet/irishub/modules/service"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/client"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func GetCmdScvDef(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Create a new service definition",
		Example: "iriscli service define --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --service-description=<service description> --author-description=<author description> " +
			"--tags=tag1,tag2 --idl-content=<interface description content> --file=test.proto",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			name := viper.GetString(FlagServiceName)
			description := viper.GetString(FlagServiceDescription)
			authorDescription := viper.GetString(FlagAuthorDescription)
			tags := viper.GetStringSlice(FlagTags)
			content := viper.GetString(FlagIdlContent)
			if len(content) > 0 {
				content = strings.Replace(content, `\n`, "\n", -1)
			}
			filePath := viper.GetString(FlagFile)
			if len(filePath) > 0 {
				contentBytes, err := cmn.ReadFile(filePath)
				if err != nil {
					return err
				}
				content = string(contentBytes)
			}
			fmt.Printf("idl condent: \n%s\n", content)
			chainId := viper.GetString(client.FlagChainID)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			msg := service.NewMsgSvcDef(name, chainId, description, tags, fromAddr, authorDescription, content)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsServiceDescription)
	cmd.Flags().AddFlagSet(FsTags)
	cmd.Flags().AddFlagSet(FsAuthorDescription)
	cmd.Flags().AddFlagSet(FsIdlContent)
	cmd.Flags().AddFlagSet(FsFile)

	return cmd
}

func GetCmdScvBind(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Create a new service binding",
		Example: "iriscli service bind --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			chainId := viper.GetString(client.FlagChainID)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTime := viper.GetInt64(FlagAvgRspTime)
			usableTime := viper.GetInt64(FlagUsableTime)
			bindingTypeStr := viper.GetString(FlagBindType)

			bindingType, err := service.BindingTypeFromString(bindingTypeStr)
			if err != nil {
				return err
			}

			deposit, err := cliCtx.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			var prices []sdk.Coin
			for _, ip := range initialPrices {
				price, err := cliCtx.ParseCoin(ip)
				if err != nil {
					return err
				}
				prices = append(prices, price)
			}

			level := service.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}
			msg := service.NewMsgSvcBind(defChainId, name, chainId, fromAddr, bindingType, deposit, prices, level)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsDeposit)
	cmd.Flags().AddFlagSet(FsPrices)
	cmd.Flags().AddFlagSet(FsBindType)
	cmd.Flags().AddFlagSet(FsAvgRspTime)
	cmd.Flags().AddFlagSet(FsUsableTime)

	return cmd
}

func GetCmdScvBindUpdate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding",
		Short: "Update a service binding",
		Example: "iriscli service update-binding --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			chainId := viper.GetString(client.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTime := viper.GetInt64(FlagAvgRspTime)
			usableTime := viper.GetInt64(FlagUsableTime)
			bindingTypeStr := viper.GetString(FlagBindType)

			var bindingType service.BindingType
			if bindingTypeStr != "" {
				bindingType, err = service.BindingTypeFromString(bindingTypeStr)
				if err != nil {
					return err
				}
			}

			var deposit sdk.Coins
			if initialDeposit != "" {
				deposit, err = cliCtx.ParseCoins(initialDeposit)
				if err != nil {
					return err
				}
			}

			var prices []sdk.Coin
			for _, ip := range initialPrices {
				price, err := cliCtx.ParseCoin(ip)
				if err != nil {
					return err
				}
				prices = append(prices, price)
			}

			level := service.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}
			msg := service.NewMsgSvcBindingUpdate(defChainId, name, chainId, fromAddr, bindingType, deposit, prices, level)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsDeposit)
	cmd.Flags().AddFlagSet(FsPrices)
	cmd.Flags().AddFlagSet(FsBindType)
	cmd.Flags().AddFlagSet(FsAvgRspTime)
	cmd.Flags().AddFlagSet(FsUsableTime)

	return cmd
}

func GetCmdScvDisable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable a available service binding",
		Example: "iriscli service disable --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			chainId := viper.GetString(client.FlagChainID)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			msg := service.NewMsgSvcDisable(defChainId, name, chainId, fromAddr)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsDefChainID)

	return cmd
}

func GetCmdScvEnable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable a unavailable service binding",
		Example: "iriscli service enable --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --deposit=1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			chainId := viper.GetString(client.FlagChainID)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			initialDeposit := viper.GetString(FlagDeposit)
			deposit, err := cliCtx.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			msg := service.NewMsgSvcEnable(defChainId, name, chainId, fromAddr, deposit)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsDefChainID)
	cmd.Flags().AddFlagSet(FsDeposit)

	return cmd
}

func GetCmdScvRefundDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit",
		Short: "Refund all deposit from a service binding",
		Example: "iriscli service refund-deposit --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			chainId := viper.GetString(client.FlagChainID)

			name := viper.GetString(FlagServiceName)
			defChainId := viper.GetString(FlagDefChainID)

			msg := service.NewMsgSvcRefundDeposit(defChainId, name, chainId, fromAddr)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsDefChainID)

	return cmd
}
