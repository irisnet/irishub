package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/client/context"
	stakeClient "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdCreateValidator implements the create validator command handler.
func GetCmdCreateValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-validator",
		Short:   "create new validator initialized with a self-delegation to it",
		Example: "iriscli stake create-validator --chain-id=<chain-id> --from=<key name> --fee=0.004iris --pubkey=<validator public key> --amount=10iris --moniker=<validator name> --commission-max-change-rate=0.1 --commission-max-rate=0.5 --commission-rate=0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			amounstStr := viper.GetString(FlagAmount)
			if amounstStr == "" {
				return fmt.Errorf("Must specify amount to stake using --amount")
			}
			amount, err := cliCtx.ParseCoin(amounstStr)
			if err != nil {
				return err
			}

			validatorAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			pkStr := viper.GetString(FlagPubKey)
			if len(pkStr) == 0 {
				return fmt.Errorf("must use --pubkey flag")
			}

			pk, err := sdk.GetConsPubKeyBech32(pkStr)
			if err != nil {
				return err
			}

			if viper.GetString(FlagMoniker) == "" {
				return fmt.Errorf("please enter a moniker for the validator using --moniker")
			}

			description := stake.Description{
				Moniker:  viper.GetString(FlagMoniker),
				Identity: viper.GetString(FlagIdentity),
				Website:  viper.GetString(FlagWebsite),
				Details:  viper.GetString(FlagDetails),
			}

			// get the initial validator commission parameters
			rateStr := viper.GetString(FlagCommissionRate)
			maxRateStr := viper.GetString(FlagCommissionMaxRate)
			maxChangeRateStr := viper.GetString(FlagCommissionMaxChangeRate)
			commissionMsg, err := stakeClient.BuildCommissionMsg(rateStr, maxRateStr, maxChangeRateStr)
			if err != nil {
				return err
			}

			var msg sdk.Msg
			if viper.GetString(FlagAddressDelegator) != "" {
				delAddr, err := sdk.AccAddressFromBech32(viper.GetString(FlagAddressDelegator))
				if err != nil {
					return err
				}

				msg = stake.NewMsgCreateValidatorOnBehalfOf(
					delAddr, sdk.ValAddress(validatorAddr), pk, amount, description, commissionMsg,
				)
			} else {
				msg = stake.NewMsgCreateValidator(
					sdk.ValAddress(validatorAddr), pk, amount, description, commissionMsg,
				)
			}

			if viper.GetBool(FlagGenesisFormat) {
				ip := viper.GetString(FlagIP)
				nodeID := viper.GetString(FlagNodeID)
				if nodeID != "" && ip != "" {
					txCtx = txCtx.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
				}
			}

			if viper.GetBool(FlagGenesisFormat) || cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg}, true)
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(fsDescriptionCreate)
	cmd.Flags().AddFlagSet(FsCommissionCreate)
	cmd.Flags().AddFlagSet(fsDelegator)
	cmd.Flags().Bool(FlagGenesisFormat, false, "Export the transaction in gen-tx format; it implies --generate-only")
	cmd.Flags().String(FlagIP, "", fmt.Sprintf("Node's public IP. It takes effect only when used in combination with --%s", FlagGenesisFormat))
	cmd.Flags().String(FlagNodeID, "", "Node's ID")
	cmd.MarkFlagRequired(FlagMoniker)
	cmd.MarkFlagRequired(FlagPubKey)
	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagCommissionRate)
	cmd.MarkFlagRequired(FlagCommissionMaxRate)
	cmd.MarkFlagRequired(FlagCommissionMaxChangeRate)
	return cmd
}

// GetCmdEditValidator implements the create edit validator command.
func GetCmdEditValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit-validator",
		Short:   "edit and existing validator account",
		Example: "iriscli stake edit-validator --chain-id=<chain-id> --from=<key name> --fee=0.004iris --moniker=<validator name> --details=<optional details> --website=<optional website>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			valAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			description := stake.Description{
				Moniker:  viper.GetString(FlagMoniker),
				Identity: viper.GetString(FlagIdentity),
				Website:  viper.GetString(FlagWebsite),
				Details:  viper.GetString(FlagDetails),
			}

			var newRate *sdk.Dec

			commissionRate := viper.GetString(FlagCommissionRate)
			if commissionRate != "" {
				rate, err := sdk.NewDecFromStr(commissionRate)
				if err != nil {
					return fmt.Errorf("invalid new commission rate: %v", err)
				}

				newRate = &rate
			}

			msg := stake.NewMsgEditValidator(sdk.ValAddress(valAddr), description, newRate)

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg}, false)
			}

			// build and sign the transaction, then broadcast to Tendermint
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsDescriptionEdit)
	cmd.Flags().AddFlagSet(fsCommissionUpdate)

	return cmd
}

// GetCmdDelegate implements the delegate command.
func GetCmdDelegate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegate",
		Short:   "delegate liquid tokens to an validator",
		Example: "iriscli stake delegate --chain-id=<chain-id> --from=<key name> --fee=0.004iris --amount=10iris --address-validator=<validator owner address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			amount, err := cliCtx.ParseCoin(viper.GetString(FlagAmount))
			if err != nil {
				return err
			}

			delegatorAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			validatorAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidator))
			if err != nil {
				return err
			}

			msg := stake.NewMsgDelegate(delegatorAddr, validatorAddr, amount)

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(fsValidator)
	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagAddressValidator)
	return cmd
}

// GetCmdRedelegate implements the redelegate validator command.
func GetCmdRedelegate(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "redelegate",
		Short:   "redelegate illiquid tokens from one validator to another",
		Example: "iriscli stake redelegate --chain-id=<chain-id> --from=<key name> --fee=0.004iris --address-validator-source=<source validator address> --address-validator-dest=<destination validator address> --shares-percent=0.5",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			var err error
			delegatorAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			validatorSrcAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidatorSrc))
			if err != nil {
				return err
			}

			validatorDstAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidatorDst))
			if err != nil {
				return err
			}

			// get the shares amount
			sharesAmountStr := viper.GetString(FlagSharesAmount)
			sharesPercentStr := viper.GetString(FlagSharesPercent)
			sharesAmount, err := getShares(
				storeName, cliCtx, cdc, sharesAmountStr, sharesPercentStr,
				delegatorAddr, validatorSrcAddr,
			)
			if err != nil {
				return err
			}

			msg := stake.NewMsgBeginRedelegate(delegatorAddr, validatorSrcAddr, validatorDstAddr, sharesAmount)

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsShares)
	cmd.Flags().AddFlagSet(fsRedelegation)
	cmd.MarkFlagRequired(FlagAddressValidatorSrc)
	cmd.MarkFlagRequired(FlagAddressValidatorDst)
	return cmd
}

// nolint: gocyclo
// TODO: Make this pass gocyclo linting
func getShares(
	storeName string, cliCtx context.CLIContext, cdc *codec.Codec, sharesAmountStr,
	sharesPercentStr string, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
) (sharesAmount sdk.Dec, err error) {
	switch {
	case sharesAmountStr != "" && sharesPercentStr != "":
		return sharesAmount, errors.Errorf("can either specify the amount OR the percent of the shares, not both")

	case sharesAmountStr == "" && sharesPercentStr == "":
		return sharesAmount, errors.Errorf("can either specify the amount OR the percent of the shares, not both")

	case sharesAmountStr != "":
		sharesAmount, err = sdk.NewDecFromStr(sharesAmountStr)
		if err != nil {
			return sharesAmount, err
		}
		if !sharesAmount.GT(sdk.ZeroDec()) {
			return sharesAmount, errors.Errorf("shares amount must be positive number (ex. 123, 1.23456789)")
		}

		stakeTokenDenom, err := cliCtx.GetCoinType(types.StakeDenomName)
		if err != nil {
			panic(err)
		}
		decimalDiff := stakeTokenDenom.MinUnit.Decimal - stakeTokenDenom.GetMainUnit().Decimal
		exRate := sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, decimalDiff))
		sharesAmount = sharesAmount.Mul(exRate)
	case sharesPercentStr != "":
		var sharesPercent sdk.Dec
		sharesPercent, err = sdk.NewDecFromStr(sharesPercentStr)
		if err != nil {
			return sharesAmount, err
		}
		if !sharesPercent.GT(sdk.ZeroDec()) || !sharesPercent.LTE(sdk.OneDec()) {
			return sharesAmount, errors.Errorf("shares percent must be >0 and <=1 (ex. 0.01, 0.75, 1)")
		}

		// make a query to get the existing delegation shares
		key := stake.GetDelegationKey(delegatorAddr, validatorAddr)
		cliCtx := context.NewCLIContext().
			WithCodec(cdc).
			WithAccountDecoder(utils.GetAccountDecoder(cdc))

		resQuery, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			return sharesAmount, errors.Errorf("cannot find delegation to determine percent Error: %v", err)
		}

		delegation, err := types.UnmarshalDelegation(cdc, key, resQuery)
		if err != nil {
			return sdk.ZeroDec(), err
		}

		sharesAmount = sharesPercent.Mul(delegation.Shares)
	}
	return
}

// GetCmdBeginUnbonding implements the begin unbonding validator command.
func GetCmdUnbond(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unbond",
		Short:   "unbond shares from a validator",
		Example: "iriscli stake unbond --chain-id=<chain-id> --from=<key name> --fee=0.004iris --address-validator=<validator address> --shares-percent=0.5",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			delegatorAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			validatorAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidator))
			if err != nil {
				return err
			}

			// get the shares amount
			sharesAmountStr := viper.GetString(FlagSharesAmount)
			sharesPercentStr := viper.GetString(FlagSharesPercent)
			sharesAmount, err := getShares(
				storeName, cliCtx, cdc, sharesAmountStr, sharesPercentStr,
				delegatorAddr, validatorAddr,
			)
			if err != nil {
				return err
			}

			msg := stake.NewMsgBeginUnbonding(delegatorAddr, validatorAddr, sharesAmount)

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsShares)
	cmd.Flags().AddFlagSet(fsValidator)
	cmd.MarkFlagRequired(FlagAddressValidator)

	return cmd
}
