package cli

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/v1/service"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func GetCmdSvcDef(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Create a new service definition",
		Example: "iriscli service define --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --service-description=<service description> --author-description=<author description> " +
			"--tags=tag1,tag2 --idl-content=<interface description content> --file=test.proto",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
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
	cmd.Flags().AddFlagSet(FsServiceDefinitionCreate)
	cmd.MarkFlagRequired(FlagServiceName)
	return cmd
}

func GetCmdSvcBind(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Create a new service binding",
		Example: "iriscli service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
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
	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBindingCreate)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	cmd.MarkFlagRequired(FlagBindType)
	cmd.MarkFlagRequired(FlagPrices)
	cmd.MarkFlagRequired(FlagAvgRspTime)
	cmd.MarkFlagRequired(FlagUsableTime)
	return cmd
}

func GetCmdSvcBindUpdate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding",
		Short: "Update a service binding",
		Example: "iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
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
	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBindingUpdate)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	return cmd
}

func GetCmdSvcDisable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable a available service binding",
		Example: "iriscli service disable --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
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
	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	return cmd
}

func GetCmdSvcEnable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable an unavailable service binding",
		Example: "iriscli service enable --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --deposit=1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
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
	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().String(FlagDeposit, "", "additional deposit of binding")
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	return cmd
}

func GetCmdSvcRefundDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit",
		Short: "Refund all deposit from a service binding",
		Example: "iriscli service refund-deposit --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
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
	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	return cmd
}

func GetCmdSvcCall(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Call a service method",
		Example: "iriscli service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --def-chain-id=<bind-chain-id> " +
			"--service-name=<service name> --method-id=<method-id> --bind-chain-id=<chain-id> --provider=<provider> --service-fee=1iris --request-data=<req>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			chainId := viper.GetString(client.FlagChainID)

			defChainId := viper.GetString(FlagDefChainID)
			name := viper.GetString(FlagServiceName)
			bindChainId := viper.GetString(FlagBindChainID)
			methodId := int16(viper.GetInt(FlagMethodID))

			providerStr := viper.GetString(FlagProvider)
			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				return err
			}

			serviceFeeStr := viper.GetString(FlagServiceFee)
			serviceFee, err := cliCtx.ParseCoins(serviceFeeStr)
			if err != nil {
				return err
			}

			inputString := viper.GetString(FlagReqData)
			input, err := hex.DecodeString(inputString)
			if err != nil {
				return err
			}

			profiling := viper.GetBool(FlagProfiling)

			msg := service.NewMsgSvcRequest(defChainId, name, bindChainId, chainId, fromAddr, provider, methodId, input, serviceFee, profiling)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBinding)
	cmd.Flags().AddFlagSet(FsServiceRequest)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)
	cmd.MarkFlagRequired(FlagBindChainID)
	cmd.MarkFlagRequired(FlagProvider)
	cmd.MarkFlagRequired(FlagMethodID)
	return cmd
}

func GetCmdSvcRespond(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond",
		Short: "Respond a service method invocation",
		Example: "iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --request-chain-id=<call-chain-id> " +
			"--request-id=<request-id> --response-data=<resp>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			reqChainId := viper.GetString(FlagReqChainId)
			outputString := viper.GetString(FlagRespData)
			output, err := hex.DecodeString(outputString)
			if err != nil {
				return err
			}

			errMsgString := viper.GetString(FlagErrMsg)
			errMsg, err := hex.DecodeString(errMsgString)
			if err != nil {
				return err
			}

			reqId := viper.GetString(FlagReqId)

			msg := service.NewMsgSvcResponse(reqChainId, reqId, fromAddr, output, errMsg)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceResponse)
	cmd.MarkFlagRequired(FlagReqChainId)
	cmd.MarkFlagRequired(FlagReqId)
	return cmd
}

func GetCmdSvcRefundFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "refund-fees",
		Short:   "Refund all fees from service call timeout",
		Example: "iriscli service refund-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --dest-address=<account address> --withdraw-amount 1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			msg := service.NewMsgSvcRefundFees(fromAddr)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdSvcWithdrawFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-fees",
		Short:   "withdraw all fees from service call reward",
		Example: "iriscli service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			msg := service.NewMsgSvcWithdrawFees(fromAddr)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdSvcWithdrawTax(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-tax",
		Short:   "withdraw service fee tax to a account",
		Example: "iriscli service withdraw-tax --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --dest-address=<account address> --withdraw-amount=1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			destAddressStr := viper.GetString(FlagDestAddress)
			destAddress, err := sdk.AccAddressFromBech32(destAddressStr)
			if err != nil {
				return err
			}

			withdrawAmountStr := viper.GetString(FlagWithdrawAmount)
			withdrawAmount, err := cliCtx.ParseCoins(withdrawAmountStr)
			if err != nil {
				return err
			}
			msg := service.NewMsgSvcWithdrawTax(fromAddr, destAddress, withdrawAmount)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceWithdrawTax)
	cmd.MarkFlagRequired(FlagDestAddress)
	cmd.MarkFlagRequired(FlagWithdrawAmount)
	return cmd
}
