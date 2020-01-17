package cli

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// GetCmdDefineService implements defining a service command
func GetCmdDefineService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Define a new service",
		Example: "iriscli service define --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--name=<service name> --description=<service description> --author-description=<author description> " +
			"--tags=tag1,tag2 --schemas=<schemas content or path/to/schemas.json>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			author, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			name := viper.GetString(FlagName)
			description := viper.GetString(FlagDescription)
			authorDescription := viper.GetString(FlagAuthorDescription)
			tags := viper.GetStringSlice(FlagTags)
			schemas := viper.GetString(FlagSchemas)

			if !json.Valid([]byte(schemas)) {
				schemasContent, err := ioutil.ReadFile(schemas)
				if err != nil {
					return fmt.Errorf("neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(schemasContent) {
					return fmt.Errorf(".json file content is invalid JSON")
				}

				schemas = string(schemasContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(schemas)); err != nil {
				return fmt.Errorf("failed to compact the schema")
			}

			schemas = buf.String()
			fmt.Printf("schemas content: \n%s\n", schemas)

			msg := service.NewMsgDefineService(name, description, tags, author, authorDescription, schemas)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefine)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagSchemas)

	return cmd
}

// GetCmdBindService implements binding a service command
func GetCmdBindService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Bind a service",
		Example: "iriscli service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --deposit=1iris --pricing=<service pricing> --withdraw-addr=<withdrawal address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			provider, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			serviceName := viper.GetString(FlagServiceName)

			depositStr := viper.GetString(FlagDeposit)
			deposit, err := cliCtx.ParseCoins(depositStr)
			if err != nil {
				return err
			}

			var withdrawAddr sdk.AccAddress
			withdrawAddrStr := viper.GetString(FlagWithdrawAddr)

			if len(withdrawAddrStr) != 0 {
				withdrawAddr, err = sdk.AccAddressFromBech32(withdrawAddrStr)
				if err != nil {
					return err
				}
			}

			pricing := viper.GetString(FlagPricing)

			if !json.Valid([]byte(pricing)) {
				pricingContent, err := ioutil.ReadFile(pricing)
				if err != nil {
					return fmt.Errorf("neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(pricingContent) {
					return fmt.Errorf(".json file content is invalid JSON")
				}

				pricing = string(pricingContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(pricing)); err != nil {
				return fmt.Errorf("failed to compact the pricing")
			}

			pricing = buf.String()

			msg := service.NewMsgBindService(serviceName, provider, deposit, pricing, withdrawAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceBind)
	_ = cmd.MarkFlagRequired(FlagServiceName)
	_ = cmd.MarkFlagRequired(FlagDeposit)
	_ = cmd.MarkFlagRequired(FlagPricing)

	return cmd
}

func GetCmdUpdateServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding",
		Short: "Update a service binding",
		Example: "iriscli service update-binding <service name> --chain-id=<chain-id> --from=<key-name> " +
			"--fee=0.3iris --deposit=1iris --pricing=<pricing>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			provider, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var deposit sdk.Coins
			depositStr := viper.GetString(FlagDeposit)

			if len(depositStr) != 0 {
				deposit, err = cliCtx.ParseCoins(depositStr)
				if err != nil {
					return err
				}
			}

			pricing := viper.GetString(FlagPricing)

			if len(pricing) != 0 {
				if !json.Valid([]byte(pricing)) {
					pricingContent, err := ioutil.ReadFile(pricing)
					if err != nil {
						return fmt.Errorf("neither JSON input nor path to .json file were provided")
					}

					if !json.Valid(pricingContent) {
						return fmt.Errorf(".json file content is invalid JSON")
					}

					pricing = string(pricingContent)
				}

				buf := bytes.NewBuffer([]byte{})
				if err := json.Compact(buf, []byte(pricing)); err != nil {
					return fmt.Errorf("failed to compact the pricing")
				}

				pricing = buf.String()
			}

			msg := service.NewMsgUpdateServiceBinding(args[0], provider, deposit, pricing)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceUpdateBinding)

	return cmd
}

func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-withdraw-addr",
		Short: "Set a new withdrawal address for a service binding",
		Example: "iriscli service set-withdraw-addr <service name> --chain-id=<chain-id> --from=<key-name> " +
			"--fee=0.3iris --withdraw-addr=<withdrawal address>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			provider, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			withdrawAddrStr := viper.GetString(FlagWithdrawAddr)
			withdrawAddr, err := sdk.AccAddressFromBech32(withdrawAddrStr)
			if err != nil {
				return err
			}

			msg := service.NewMsgSetWithdrawAddress(args[0], provider, withdrawAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceSetWithdrawAddr)
	_ = cmd.MarkFlagRequired(FlagWithdrawAddr)

	return cmd
}

func GetCmdDisableService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable",
		Short:   "Disable an available service binding",
		Example: "iriscli service disable <service name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			provider, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := service.NewMsgDisableService(args[0], provider)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdEnableService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable an unavailable service binding",
		Example: "iriscli service enable <service name> --chain-id=<chain-id> --from=<key-name> " +
			"--fee=0.3iris --deposit=1iris",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			provider, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			depositStr := viper.GetString(FlagDeposit)

			var deposit sdk.Coins
			if len(depositStr) != 0 {
				deposit, err = cliCtx.ParseCoins(depositStr)
				if err != nil {
					return err
				}
			}

			msg := service.NewMsgEnableService(args[0], provider, deposit)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceEnable)

	return cmd
}

func GetCmdRefundServiceDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit",
		Short: "Refund the deposit from a service binding",
		Example: "iriscli service refund-deposit <service name> --chain-id=<chain-id> --from=<key-name> " +
			"--fee=0.3iris",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			provider, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := service.NewMsgRefundServiceDeposit(args[0], provider)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

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

			chainID := viper.GetString(client.FlagChainID)

			defChainID := viper.GetString(FlagDefChainID)
			name := viper.GetString(FlagServiceName)
			bindChainID := viper.GetString(FlagBindChainID)
			methodID := int16(viper.GetInt(FlagMethodID))

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

			msg := service.NewMsgSvcRequest(defChainID, name, bindChainID, chainID, fromAddr, provider, methodID, input, serviceFee, profiling)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

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

			reqChainID := viper.GetString(FlagReqChainID)
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

			reqID := viper.GetString(FlagReqID)

			msg := service.NewMsgSvcResponse(reqChainID, reqID, fromAddr, output, errMsg)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceResponse)
	_ = cmd.MarkFlagRequired(FlagReqChainID)
	_ = cmd.MarkFlagRequired(FlagReqID)
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
