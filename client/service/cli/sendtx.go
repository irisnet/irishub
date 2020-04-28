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
					return fmt.Errorf("invalid schemas: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(schemasContent) {
					return fmt.Errorf("invalid schemas: .json file content is invalid JSON")
				}

				schemas = string(schemasContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(schemas)); err != nil {
				return fmt.Errorf("failed to compact the schemas")
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
			"--service-name=<service-name> --deposit=1iris --pricing=<pricing content or path/to/pricing.json> --qos=50",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var provider sdk.AccAddress

			providerStr := viper.GetString(FlagProvider)
			if len(providerStr) > 0 {
				provider, err = sdk.AccAddressFromBech32(providerStr)
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			serviceName := viper.GetString(FlagServiceName)
			qos := uint64(viper.GetInt64(FlagQoS))

			depositStr := viper.GetString(FlagDeposit)
			deposit, err := cliCtx.ParseCoins(depositStr)
			if err != nil {
				return err
			}

			pricing := viper.GetString(FlagPricing)

			if !json.Valid([]byte(pricing)) {
				pricingContent, err := ioutil.ReadFile(pricing)
				if err != nil {
					return fmt.Errorf("invalid pricing: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(pricingContent) {
					return fmt.Errorf("invalid pricing: .json file content is invalid JSON")
				}

				pricing = string(pricingContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(pricing)); err != nil {
				return fmt.Errorf("failed to compact the pricing")
			}

			pricing = buf.String()

			msg := service.NewMsgBindService(serviceName, provider, deposit, pricing, qos, owner)
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
	_ = cmd.MarkFlagRequired(FlagQoS)

	return cmd
}

// GetCmdUpdateServiceBinding implements updating a service binding command
func GetCmdUpdateServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding [service-name] [provider-address]",
		Short: "Update an existing service binding",
		Example: "iriscli service update-binding <service-name> <provider-address> --chain-id=<chain-id> --from=<key-name> " +
			"--fee=0.3iris --deposit=1iris --pricing=<pricing content or path/to/pricing.json> --qos=50",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
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
						return fmt.Errorf("invalid pricing: neither JSON input nor path to .json file were provided")
					}

					if !json.Valid(pricingContent) {
						return fmt.Errorf("invalid pricing: .json file content is invalid JSON")
					}

					pricing = string(pricingContent)
				}

				buf := bytes.NewBuffer([]byte{})
				if err := json.Compact(buf, []byte(pricing)); err != nil {
					return fmt.Errorf("failed to compact the pricing")
				}

				pricing = buf.String()
			}

			qos := uint64(viper.GetInt64(FlagQoS))

			msg := service.NewMsgUpdateServiceBinding(args[0], provider, deposit, pricing, qos, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceUpdateBinding)

	return cmd
}

// GetCmdSetWithdrawAddr implements setting a withdrawal address command
func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-withdraw-addr [withdrawal-address]",
		Short: "Set a withdrawal address for an owner",
		Example: "iriscli service set-withdraw-addr <withdrawal-address> --chain-id=<chain-id> " +
			"--from=<key-name> --fee=0.3iris",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			withdrawAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := service.NewMsgSetWithdrawAddress(owner, withdrawAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdDisableServiceBinding implements disabling a service binding command
func GetCmdDisableServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable [service-name] [provider-address]",
		Short:   "Disable an available service binding",
		Example: "iriscli service disable <service-name> <provider-address> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			msg := service.NewMsgDisableServiceBinding(args[0], provider, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdEnableServiceBinding implements enabling a service binding command
func GetCmdEnableServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable [service-name] [provider-address]",
		Short: "Enable an unavailable service binding",
		Example: "iriscli service enable <service-name> <provider-address> --chain-id=<chain-id> --from=<key-name> " +
			"--fee=0.3iris --deposit=1iris",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			var deposit sdk.Coins

			depositStr := viper.GetString(FlagDeposit)
			if len(depositStr) != 0 {
				deposit, err = cliCtx.ParseCoins(depositStr)
				if err != nil {
					return err
				}
			}

			msg := service.NewMsgEnableServiceBinding(args[0], provider, deposit, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceEnableBinding)

	return cmd
}

// GetCmdRefundServiceDeposit implements refunding deposit command
func GetCmdRefundServiceDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit [service-name] [provider-address]",
		Short: "Refund all deposit from a service binding",
		Example: "iriscli service refund-deposit <service-name> <provider-address> --chain-id=<chain-id> --from=<key-name> " +
			"--fee=0.3iris",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			msg := service.NewMsgRefundServiceDeposit(args[0], provider, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdCallService implements initiating a service call command
func GetCmdCallService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Initiate a service call",
		Example: "iriscli service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --service-name=<service-name> " +
			"--providers=<provider list> --service-fee-cap=1iris --data=<input content or path/to/input.json> --timeout=100 --repeated --frequency=150 --total=100",
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

			if !json.Valid([]byte(input)) {
				inputContent, err := ioutil.ReadFile(input)
				if err != nil {
					return fmt.Errorf("invalid input data: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(inputContent) {
					return fmt.Errorf("invalid input data: .json file content is invalid JSON")
				}

				input = string(inputContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(input)); err != nil {
				return fmt.Errorf("failed to compact the input data")
			}

			input = buf.String()
			timeout := viper.GetInt64(FlagTimeout)
			superMode := viper.GetBool(FlagSuperMode)
			repeated := viper.GetBool(FlagRepeated)

			frequency := uint64(0)
			total := int64(0)

			if repeated {
				frequency = uint64(viper.GetInt64(FlagFrequency))
				total = viper.GetInt64(FlagTotal)
			}

			msg := service.NewMsgCallService(
				serviceName, providers, consumer, input, serviceFeeCap,
				timeout, superMode, repeated, frequency, total,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceCall)
	_ = cmd.MarkFlagRequired(FlagServiceName)
	_ = cmd.MarkFlagRequired(FlagProviders)
	_ = cmd.MarkFlagRequired(FlagServiceFeeCap)
	_ = cmd.MarkFlagRequired(FlagData)
	_ = cmd.MarkFlagRequired(FlagTimeout)

	return cmd
}

// GetCmdRespondService implements responding to a service request command
func GetCmdRespondService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond",
		Short: "Respond to a service request",
		Example: "iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--request-id=<request-id> --result=<result content or path/to/result.json> --data=<output content or path/to/output.json>",
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

			requestIDStr := viper.GetString(FlagRequestID)
			requestID, err := service.ConvertRequestID(requestIDStr)
			if err != nil {
				return err
			}

			result := viper.GetString(FlagResult)
			output := viper.GetString(FlagData)

			if len(result) > 0 {
				if !json.Valid([]byte(result)) {
					resultContent, err := ioutil.ReadFile(result)
					if err != nil {
						return fmt.Errorf("invalid result: neither JSON input nor path to .json file were provided")
					}

					if !json.Valid(resultContent) {
						return fmt.Errorf("invalid result: .json file content is invalid JSON")
					}

					result = string(resultContent)
				}

				buf := bytes.NewBuffer([]byte{})
				if err := json.Compact(buf, []byte(result)); err != nil {
					return fmt.Errorf("failed to compact the result")
				}

				result = buf.String()
			}

			if len(output) > 0 {
				if !json.Valid([]byte(output)) {
					outputContent, err := ioutil.ReadFile(output)
					if err != nil {
						return fmt.Errorf("invalid output data: neither JSON input nor path to .json file were provided")
					}

					if !json.Valid(outputContent) {
						return fmt.Errorf("invalid output data: .json file content is invalid JSON")
					}

					output = string(outputContent)
				}

				buf := bytes.NewBuffer([]byte{})
				if err := json.Compact(buf, []byte(output)); err != nil {
					return fmt.Errorf("failed to compact the output data")
				}

				output = buf.String()
			}

			msg := service.NewMsgRespondService(requestID, provider, result, output)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceRespond)
	_ = cmd.MarkFlagRequired(FlagRequestID)
	_ = cmd.MarkFlagRequired(FlagResult)

	return cmd
}

// GetCmdPauseRequestContext implements pausing a request context command
func GetCmdPauseRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause [request-context-id]",
		Short:   "Pause a running request context",
		Example: "iriscli service pause <request-context-id> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
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

// GetCmdStartRequestContext implements restarting a request context command
func GetCmdStartRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start [request-context-id]",
		Short:   "Start a paused request context",
		Example: "iriscli service start <request-context-id> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
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

// GetCmdKillRequestContext implements terminating a request context command
func GetCmdKillRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kill [request-context-id]",
		Short:   "Terminate a request context",
		Example: "iriscli service kill <request-context-id> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
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

// GetCmdUpdateRequestContext implements updating a request context command
func GetCmdUpdateRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [request-context-id]",
		Short: "Update a request context",
		Example: "iriscli service update <request-context-id> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--providers=<new providers> --service-fee-cap=2iris --timeout=0 --frequency=200 --total=200",
		Args: cobra.ExactArgs(1),
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

			var providers []sdk.AccAddress
			providerList := viper.GetStringSlice(FlagProviders)

			for _, p := range providerList {
				provider, err := sdk.AccAddressFromBech32(p)
				if err != nil {
					return err
				}

				providers = append(providers, provider)
			}

			var serviceFeeCap sdk.Coins

			serviceFeeCapStr := viper.GetString(FlagServiceFeeCap)
			if len(serviceFeeCapStr) != 0 {
				serviceFeeCap, err = cliCtx.ParseCoins(serviceFeeCapStr)
				if err != nil {
					return err
				}
			}

			timeout := viper.GetInt64(FlagTimeout)
			frequency := uint64(viper.GetInt64(FlagFrequency))
			total := viper.GetInt64(FlagTotal)

			msg := service.NewMsgUpdateRequestContext(
				requestContextID, providers, serviceFeeCap,
				timeout, frequency, total, consumer,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceUpdateRequestContext)

	return cmd
}

// GetCmdWithdrawEarnedFees implements withdrawing earned fees command
func GetCmdWithdrawEarnedFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-fees [provider-address]",
		Short:   "Withdraw the earned fees of the specified provider or all providers if not given",
		Example: "iriscli service withdraw-fees <provider-address> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			var provider sdk.AccAddress

			if len(args) > 0 {
				provider, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			msg := service.NewMsgWithdrawEarnedFees(owner, provider)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdWithdrawTax implements withdrawing tax command
func GetCmdWithdrawTax(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-tax [destination-address] [withdrawal-amount]",
		Short: "Withdraw the service tax",
		Example: "iriscli service withdraw-tax <destination-address> <withdrawal-amount> --chain-id=<chain-id> " +
			"--from=<key-name> --fee=0.3iris",
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

			amount, err := cliCtx.ParseCoins(args[1])
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
