package cli

import (
	"os"
	"fmt"
	"strings"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	iservice1 "github.com/irisnet/irishub/modules/iservice1"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/client"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/irisnet/irishub/modules/iservice"
)

func GetCmdScvDef(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "create new service definition",
		Example: "iriscli iservice define --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--service-name=<service name> --name=<service name> --service-description=<service description> --author-description=<author description> " +
			"--tags=tag1 --messaging=Unicast --broadcast=Broadcast --idl-content=<interface description content> --file=test.proto",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			var newVersion bool
			messagingStr := viper.GetString(FlagMessaging)
			var messaging iservice1.MessagingType
			var broadcast iservice.BroadcastEnum
			var err error
			if messagingStr != "" {
				messaging, err = iservice1.MessagingTypeFromString(messagingStr)
				if err != nil {
					return err
				}
				newVersion = true
			} else {
				broadcastStr := viper.GetString(FlagBroadcast)
				broadcast, err = iservice.BroadcastEnumFromString(broadcastStr)
				if err != nil {
					return err
				}
			}

			name := viper.GetString(FlagServiceName)
			if newVersion {
				name = viper.GetString(FlagService1Name)
			}

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

			if newVersion {
				msg := iservice1.NewMsgSvcDef(name, chainId, description, tags, fromAddr, authorDescription, content, messaging)
				cliCtx.PrintResponse = true
				return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
			} else {
				msg := iservice.NewMsgSvcDef(name, chainId, description, tags, fromAddr, authorDescription, content, broadcast)
				cliCtx.PrintResponse = true
				return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
			}
		},
	}
	cmd.Flags().AddFlagSet(FsServiceName)
	cmd.Flags().AddFlagSet(FsService1Name)
	cmd.Flags().AddFlagSet(FsServiceDescription)
	cmd.Flags().AddFlagSet(FsTags)
	cmd.Flags().AddFlagSet(FsAuthorDescription)
	cmd.Flags().AddFlagSet(FsIdlContent)
	cmd.Flags().AddFlagSet(FsMessaging)
	cmd.Flags().AddFlagSet(FsBroadcast)
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

			name := viper.GetString(FlagService1Name)
			defChainId := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTimeStr := viper.GetString(FlagAvgRspTime)
			usableTimeStr := viper.GetString(FlagUsableTime)
			expirationStr := viper.GetString(FlagExpiration)
			bindingTypeStr := viper.GetString(FlagBindType)

			bindingType, err := iservice1.BindingTypeFromString(bindingTypeStr)
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
			level := iservice1.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}
			msg := iservice1.NewMsgSvcBind(defChainId, name, chainId, fromAddr, bindingType, deposit, prices, level, expiration)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsService1Name)
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
			name := viper.GetString(FlagService1Name)
			defChainId := viper.GetString(FlagDefChainID)
			initialDeposit := viper.GetString(FlagDeposit)
			initialPrices := viper.GetStringSlice(FlagPrices)
			avgRspTimeStr := viper.GetString(FlagAvgRspTime)
			usableTimeStr := viper.GetString(FlagUsableTime)
			expirationStr := viper.GetString(FlagExpiration)
			bindingTypeStr := viper.GetString(FlagBindType)

			var bindingType iservice1.BindingType
			if bindingTypeStr != "" {
				bindingType, err = iservice1.BindingTypeFromString(bindingTypeStr)
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

			level := iservice1.Level{AvgRspTime: avgRspTime, UsableTime: usableTime}
			msg := iservice1.NewMsgSvcBindingUpdate(defChainId, name, chainId, fromAddr, bindingType, deposit, prices, level, expiration)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsService1Name)
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

			name := viper.GetString(FlagService1Name)
			defChainId := viper.GetString(FlagDefChainID)

			msg := iservice1.NewMsgSvcRefundDeposit(defChainId, name, chainId, fromAddr)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsService1Name)
	cmd.Flags().AddFlagSet(FsDefChainID)

	return cmd
}
