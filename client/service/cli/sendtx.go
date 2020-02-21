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

func GetCmdRequestService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Call a service",
		Example: "iriscli service call --chain-id=<chain-id> --from=<key name> --fee=0.3iris --service-name=<service name> " +
			"--providers=<provider list> --service-fee-cap=1iris --data=<request data> --repeated --frequency=10 --total=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			consumer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			serviceName := viper.GetString(FlagServiceName)

			var providers []sdk.AccAddress
			providerList := viper.GetStringSlice(FlagProviders)

			for _, p := range providerList {
				provider, err := sdk.AccAddressFromBech32(p)
				if err != nil {
					return err
				}

				providers = append(providers, provider)
			}

			serviceFeeCap, err := cliCtx.ParseCoins(viper.GetString(FlagServiceFeeCap))
			if err != nil {
				return err
			}

			input := viper.GetString(FlagData)
			timeout := viper.GetInt64(FlagTimeout)
			repeated := viper.GetBool(FlagRepeated)

			frequency := uint64(0)
			total := int64(0)

			if repeated {
				frequency = uint64(viper.GetInt64(FlagFrequency))
				total = viper.GetInt64(FlagTotal)
			}

			msg := service.NewMsgRequestService(serviceName, providers, consumer, input, serviceFeeCap, timeout, repeated, frequency, total)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceRequest)
	_ = cmd.MarkFlagRequired(FlagServiceName)
	_ = cmd.MarkFlagRequired(FlagProviders)
	_ = cmd.MarkFlagRequired(FlagServiceFeeCap)
	_ = cmd.MarkFlagRequired(FlagData)

	return cmd
}

func GetCmdRespondService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond",
		Short: "Respond a service request",
		Example: "iriscli service respond --chain-id=<chain-id> --from=<key name> --fee=0.3iris " +
			"--request-id=<request-id> --data=<response data>",
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

			requestID := viper.GetString(FlagRequestID)
			output := viper.GetString(FlagData)
			errMsg := viper.GetString(FlagError)

			msg := service.NewMsgRespondService(requestID, provider, output, errMsg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceRespond)
	_ = cmd.MarkFlagRequired(FlagRequestID)

	return cmd
}

func GetCmdPauseRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause",
		Short:   "Pause a request context",
		Example: "iriscli service pause <request-context-id> --chain-id=<chain-id> --from=<key name> --fee=0.3iris",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			consumer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := service.NewMsgPauseRequestContext(requestContextID, consumer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdStartRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Resume a paused request context",
		Example: "iriscli service start <request-context-id> --chain-id=<chain-id> --from=<key name> --fee=0.3iris",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			consumer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := service.NewMsgStartRequestContext(requestContextID, consumer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdKillRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kill",
		Short:   "Terminate a request context",
		Example: "iriscli service kill <request-context-id> --chain-id=<chain-id> --from=<key name> --fee=0.3iris",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			consumer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := service.NewMsgKillRequestContext(requestContextID, consumer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdWithdrawEarnedFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-fees",
		Short:   "Withdraw the earned fees",
		Example: "iriscli service withdraw-fees --chain-id=<chain-id> --from=<key name> --fee=0.3iris",
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

			msg := service.NewMsgWithdrawEarnedFees(provider)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdWithdrawTax(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-tax",
		Short: "Withdraw service tax",
		Example: "iriscli service withdraw-tax <destination address> <withdrawal amount> --chain-id=<chain-id> " +
			"--from=<key name> --fee=0.3iris",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			trustee, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			destAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := service.NewMsgWithdrawTax(trustee, destAddr, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
