package cli

import (
	"os"
	"fmt"
	"strings"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/modules/iservice"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/client"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func GetCmdScvDef(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "create new service definition",
		Example: "iriscli iservice define --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --service-description=<service description> --author-description=<author description> " +
			"--tags=tag1,tag2 --messaging=Unicast --idl-content=<interface description content> --file=test.proto",
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
			messagingStr := viper.GetString(FlagMessaging)
			chainId := viper.GetString(client.FlagChainID)

			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			messaging, err := iservice.MessagingTypeFromString(messagingStr)
			if err != nil {
				return err
			}

			msg := iservice.NewMsgSvcDef(name, chainId, description, tags, fromAddr, authorDescription, content, messaging)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsServiceDescription)
	cmd.Flags().AddFlagSet(FsTags)
	cmd.Flags().AddFlagSet(FsAuthorDescription)
	cmd.Flags().AddFlagSet(FsIdlContent)
	cmd.Flags().AddFlagSet(FsMessaging)
	cmd.Flags().AddFlagSet(FsFile)

	return cmd
}

func GetCmdScvBind(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "create new service binding",
		Example: "iriscli iservice bind --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100 --expiration=-1",
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
			avgRspTimeStr := viper.GetString(FlagAvgRspTime)
			usableTimeStr := viper.GetString(FlagUsableTime)
			expirationStr := viper.GetString(FlagExpiration)
			bindingTypeStr := viper.GetString(FlagBindType)

			bindingType, err := iservice.BindingTypeFromString(bindingTypeStr)
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

			avgRspTime, err := strconv.ParseInt(avgRspTimeStr, 10, 64)
			if err != nil {
				return err
			}
			usableTime, err := strconv.ParseInt(usableTimeStr, 10, 64)
			if err != nil {
				return err
			}
			expiration, err := strconv.ParseInt(expirationStr, 10, 64)
			if err != nil {
				return err
			}
			level := iservice.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}
			msg := iservice.NewMsgSvcBind(defChainId, name, chainId, fromAddr, bindingType, deposit, prices, level, expiration)
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
	cmd.Flags().AddFlagSet(FsExpiration)

	return cmd
}

func GetCmdScvBindUpdate(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding",
		Short: "update a service binding",
		Example: "iriscli iservice update-binding --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --def-chain-id=<chain-id> --bind-type=Local " +
			"--deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100 --expiration=-1",
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
			avgRspTimeStr := viper.GetString(FlagAvgRspTime)
			usableTimeStr := viper.GetString(FlagUsableTime)
			expirationStr := viper.GetString(FlagExpiration)
			bindingTypeStr := viper.GetString(FlagBindType)

			var bindingType iservice.BindingType
			if bindingTypeStr != "" {
				bindingType, err = iservice.BindingTypeFromString(bindingTypeStr)
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

			var avgRspTime int64
			if avgRspTimeStr != "" {
				avgRspTime, err = strconv.ParseInt(avgRspTimeStr, 10, 64)
				if err != nil {
					return err
				}
			}

			var usableTime int64
			if usableTimeStr != "" {
				usableTime, err = strconv.ParseInt(usableTimeStr, 10, 64)
				if err != nil {
					return err
				}
			}

			var expiration int64
			if expirationStr != "" {
				expiration, err = strconv.ParseInt(expirationStr, 10, 64)
				if err != nil {
					return err
				}
			}

			level := iservice.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}
			msg := iservice.NewMsgSvcBindingUpdate(defChainId, name, chainId, fromAddr, bindingType, deposit, prices, level, expiration)
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
	cmd.Flags().AddFlagSet(FsExpiration)

	return cmd
}

func GetCmdScvRefundDeposit(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit",
		Short: "refund all deposit from a service binding",
		Example: "iriscli iservice refund-deposit --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
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

			msg := iservice.NewMsgSvcRefundDeposit(defChainId, name, chainId, fromAddr)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsDefChainID)

	return cmd
}
