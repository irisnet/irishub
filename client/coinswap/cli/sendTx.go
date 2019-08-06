package cli

import (
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

// SendTxCmd will create a send tx and sign it with the given key.
func GetCmdAddLiquidity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-liquidity",
		Short:   "Create and sign a send tx",
		Example: "iriscli swap add-liquidity --deposit=1eth --amount=1500iris --period=10s --from=<key name> --fee=0.4iris --chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			// parse coins trying to be sent
			depositTokenStr := viper.GetString(flagDeposit)
			depositToken, err := cliCtx.ParseCoin(depositTokenStr)
			if err != nil {
				return err
			}

			nativeTokenStr := viper.GetString(flagAmount)
			nativeToken, err := cliCtx.ParseCoin(nativeTokenStr)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			periodStr := viper.GetString(flagPeriod)
			period, err := time.ParseDuration(periodStr)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(period)

			minReward := sdk.ZeroInt()
			minRewardStr := viper.GetString(flagMinReward)
			if len(minRewardStr) > 0 {
				reward, ok := sdk.NewIntFromString(minRewardStr)
				if !ok {
					return fmt.Errorf("invalid min reward:%s", minRewardStr)
				}
				minReward = reward
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			coins := sdk.NewCoins(depositToken, nativeToken)
			if !account.GetCoins().IsAllGTE(coins) {
				return fmt.Errorf("Address %s doesn't have enough coins to pay for this transaction", from.String())
			}

			msg := coinswap.NewMsgAddLiquidity(depositToken, nativeToken.Amount, minReward, deadline, from)
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsAddLiquidity)
	cmd.MarkFlagRequired(flagDeposit)
	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagPeriod)

	return cmd
}

func GetCmdRemoveLiquidity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove-liquidity",
		Short:   "Create and sign a send tx",
		Example: "iriscli swap remove-liquidity --min-token=1eth --min-native=1000 --amount=1000uni --from <key name> --fee=0.4iris --chain-id=<chain-id> --amount=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			var minToken, liquidity sdk.Coin
			var minNative = sdk.ZeroInt()
			// parse coins trying to be sent
			minTokenStr := viper.GetString(flagMinToken)
			if len(minTokenStr) > 0 {
				coin, err := cliCtx.ParseCoin(minTokenStr)
				if err != nil {
					return err
				} else {
					minToken = coin
				}
			}

			minNativeStr := viper.GetString(flagMinNative)
			if len(minNativeStr) > 0 {
				amt, ok := sdk.NewIntFromString(minNativeStr)
				if !ok {
					return fmt.Errorf("invalid amount:%s", minNativeStr)
				} else {
					minNative = amt
				}
			}

			liquidityStr := viper.GetString(flagAmount)
			if len(liquidityStr) > 0 {
				coin, err := cliCtx.ParseCoin(liquidityStr)
				if err != nil {
					return err
				} else {
					liquidity = coin
				}
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}
			// ensure account has enough coins
			coins := sdk.NewCoins(liquidity)
			if !account.GetCoins().IsAllGTE(coins) {
				return fmt.Errorf("Address %s doesn't have enough coins to pay for this transaction", from.String())
			}

			periodStr := viper.GetString(flagPeriod)
			period, err := time.ParseDuration(periodStr)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(period)

			msg := coinswap.NewMsgRemoveLiquidity(minToken, liquidity.Amount, minNative, deadline, from)
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsRemoveLiquidity)
	cmd.MarkFlagRequired(flagAmount)
	cmd.MarkFlagRequired(flagPeriod)
	return cmd
}

func GetCmdPlaceOrder(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "place-order",
		Short:   "swap tokens with other token",
		Example: "iriscli swap place-order --input-token=<token> --output-token=<token> --deadline=<token> --recipient=<addr> --is-buy-order=<true|false> --from <key name> --fee=0.4iris --chain-id=<chain-id> --amount=10iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).WithCliCtx(cliCtx)

			// parse coins trying to be sent
			inputStr := viper.GetString(flagInput)
			input, err := cliCtx.ParseCoin(inputStr)
			if err != nil {
				return err
			}
			outputStr := viper.GetString(flagOutput)
			output, err := cliCtx.ParseCoin(outputStr)
			if err != nil {
				return err
			}
			periodStr := viper.GetString(flagDeadline)
			period, err := time.ParseDuration(periodStr)
			if err != nil {
				return err
			}
			deadline := time.Now().Add(period)

			recipientStr := viper.GetString(FlagRecipient)
			recipient, err := sdk.AccAddressFromBech32(recipientStr)
			if err != nil {
				return err
			}

			isBuyOrder := viper.GetBool(FlagIsBuyOrder)

			sender, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := coinswap.NewMsgSwapOrder(input, output, deadline, sender, recipient, isBuyOrder)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsSwapTokens)
	return cmd
}
