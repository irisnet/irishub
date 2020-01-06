package cli

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	serviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Service transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceTxCmd.AddCommand(flags.PostCommands(
		GetCmdDefineService(cdc),
		GetCmdSvcBind(cdc),
		GetCmdSvcBindUpdate(cdc),
		GetCmdSvcDisable(cdc),
		GetCmdSvcEnable(cdc),
		GetCmdSvcRefundDeposit(cdc),
		GetCmdSvcCall(cdc),
		GetCmdSvcRespond(cdc),
		GetCmdSvcRefundFees(cdc),
		GetCmdSvcWithdrawFees(cdc),
		GetCmdSvcWithdrawTax(cdc),
	)...)

	return serviceTxCmd
}

// GetCmdDefineService implements defining a service command.
func GetCmdDefineService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Define a new service",
		Example: "iriscli tx service define --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--name=<service name> --description=<service description> --author-description=<author description> " +
			"--tags=tag1,tag2 --schemas=<schemas content or path/to/schemas.json>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			author := cliCtx.GetFromAddress()
			name := viper.GetString(FlagName)
			description := viper.GetString(FlagDescription)
			authorDescription := viper.GetString(FlagAuthorDescription)
			tags := viper.GetStringSlice(FlagTags)
			schemas := viper.GetString(FlagSchemas)

			if !json.Valid([]byte(schemas)) {
				schemasContent, err := ioutil.ReadFile(schemas)
				if err != nil {
					return errors.New("neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(schemasContent) {
					return errors.New(".json file content is invalid JSON")
				}

				schemas = string(schemasContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(schemas)); err != nil {
				return errors.New("failed to compact the schema")
			}

			fmt.Printf("schemas content: \n%s\n", buf.String())

			msg := types.NewMsgDefineService(name, description, tags, author, authorDescription, schemas)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefine)
	cmd.MarkFlagRequired(FlagName)
	cmd.MarkFlagRequired(FlagSchemas)

	return cmd
}

// GetCmdSvcBind implements the create service bind command.
func GetCmdSvcBind(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Create a new service binding",
		Example: "iriscli tx service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTime := viper.GetInt64(FlagAvgRspTime)
			usableTime := viper.GetInt64(FlagUsableTime)
			bindingTypeStr := viper.GetString(FlagBindType)

			bindingType, err := types.BindingTypeFromString(bindingTypeStr)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			var prices []sdk.Coin
			for _, ip := range initialPrices {
				price, err := sdk.ParseCoin(ip)
				if err != nil {
					return err
				}
				prices = append(prices, price)
			}

			level := types.Level{
				AvgRspTime: avgRspTime,
				UsableTime: usableTime,
			}

			msg := types.NewMsgSvcBind(defChainID, name, chainID, fromAddr, bindingType, deposit, prices, level)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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

// GetCmdSvcBindUpdate implements the update service bind command.
func GetCmdSvcBindUpdate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding",
		Short: "Update a service binding",
		Example: "iriscli tx service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTime := viper.GetInt64(FlagAvgRspTime)
			usableTime := viper.GetInt64(FlagUsableTime)
			bindingTypeStr := viper.GetString(FlagBindType)

			var err error

			var bindingType types.BindingType
			if bindingTypeStr != "" {
				bindingType, err = types.BindingTypeFromString(bindingTypeStr)
				if err != nil {
					return err
				}
			}

			var deposit sdk.Coins
			if initialDeposit != "" {
				deposit, err = sdk.ParseCoins(initialDeposit)
				if err != nil {
					return err
				}
			}

			var prices []sdk.Coin
			for _, ip := range initialPrices {
				price, err := sdk.ParseCoin(ip)
				if err != nil {
					return err
				}
				prices = append(prices, price)
			}

			level := types.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}

			msg := types.NewMsgSvcBindingUpdate(defChainID, name, chainID, fromAddr, bindingType, deposit, prices, level)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().AddFlagSet(FsServiceBindingUpdate)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

// GetCmdSvcDisable implements the disable service binding command.
func GetCmdSvcDisable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable an available service binding",
		Example: "iriscli tx service disable --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)

			msg := types.NewMsgSvcDisable(defChainID, name, chainID, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

// GetCmdSvcEnable implements the enable service binding command.
func GetCmdSvcEnable(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable an unavailable service binding",
		Example: "iriscli tx service enable --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --deposit=1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)

			initialDeposit := viper.GetString(FlagDeposit)
			deposit, err := sdk.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}

			msg := types.NewMsgSvcEnable(defChainID, name, chainID, fromAddr, deposit)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.Flags().String(FlagDeposit, "", "additional deposit of binding")
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

// GetCmdSvcRefundDeposit implements the refund all deposit command.
func GetCmdSvcRefundDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit",
		Short: "Refund all deposit from a service binding",
		Example: "iriscli tx service refund-deposit --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--service-name=<service name> --def-chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
			name := viper.GetString(FlagServiceName)
			defChainID := viper.GetString(FlagDefChainID)

			msg := types.NewMsgSvcRefundDeposit(defChainID, name, chainID, fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceDefinition)
	cmd.MarkFlagRequired(FlagDefChainID)
	cmd.MarkFlagRequired(FlagServiceName)

	return cmd
}

// GetCmdSvcCall implements the call service method command.
func GetCmdSvcCall(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Call a service method",
		Example: "iriscli tx service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --def-chain-id=<bind-chain-id> " +
			"--service-name=<service name> --method-id=<method-id> --bind-chain-id=<chain-id> --provider=<provider> --service-fee=1iris --request-data=<req>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			chainID := viper.GetString(flags.FlagChainID)
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
			serviceFee, err := sdk.ParseCoins(serviceFeeStr)
			if err != nil {
				return err
			}

			inputString := viper.GetString(FlagReqData)
			input, err := hex.DecodeString(inputString)
			if err != nil {
				return err
			}

			profiling := viper.GetBool(FlagProfiling)

			msg := types.NewMsgSvcRequest(defChainID, name, bindChainID, chainID, fromAddr, provider, methodID, input, serviceFee, profiling)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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

// GetCmdSvcRespond implements the respond service method invocation command.
func GetCmdSvcRespond(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond",
		Short: "Respond a service method invocation",
		Example: "iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --request-chain-id=<call-chain-id> " +
			"--request-id=<request-id> --response-data=<resp>",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

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

			msg := types.NewMsgSvcResponse(reqChainID, reqID, fromAddr, output, errMsg)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceResponse)
	cmd.MarkFlagRequired(FlagReqChainID)
	cmd.MarkFlagRequired(FlagReqID)

	return cmd
}

// GetCmdSvcRefundFees implements the refund all fees command.
func GetCmdSvcRefundFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "refund-fees",
		Short:   "Refund all fees from service call timeout",
		Example: "iriscli tx service refund-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --dest-address=<account address> --withdraw-amount 1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgSvcRefundFees(fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdSvcWithdrawFees implements the withdraw all fees command.
func GetCmdSvcWithdrawFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-fees",
		Short:   "Withdraw all fees from service call reward",
		Example: "iriscli tx service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgSvcWithdrawFees(fromAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdSvcWithdrawTax implements the withdraw service fee tax command.
func GetCmdSvcWithdrawTax(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-tax",
		Short:   "Withdraw service fee tax to an account",
		Example: "iriscli tx service withdraw-tax --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --dest-address=<account address> --withdraw-amount=1iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			fromAddr := cliCtx.GetFromAddress()

			destAddressStr := viper.GetString(FlagDestAddress)
			destAddress, err := sdk.AccAddressFromBech32(destAddressStr)
			if err != nil {
				return err
			}

			withdrawAmountStr := viper.GetString(FlagWithdrawAmount)
			withdrawAmount, err := sdk.ParseCoins(withdrawAmountStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgSvcWithdrawTax(fromAddr, destAddress, withdrawAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsServiceWithdrawTax)
	cmd.MarkFlagRequired(FlagDestAddress)
	cmd.MarkFlagRequired(FlagWithdrawAmount)

	return cmd
}
