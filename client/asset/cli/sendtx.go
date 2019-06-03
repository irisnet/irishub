package cli

import (
	"os"

	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdCreateGateway implements the create gateway command
func GetCmdCreateGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gateway",
		Short: "create a gateway",
		Example: "iriscli asset create-gateway --moniker=<moniker> --identity=<identity> --details=<details>" +
			"--website=<website> --redeem-address=<redeeming address> --operators=<comma-seperated operator addresses>",
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

			var redeemAddr sdk.AccAddress
			redeemAddrStr := viper.GetString(FlagRedeemAddress)
			if redeemAddrStr == "" {
				redeemAddr = owner
			} else {
				redeemAddr, err = sdk.AccAddressFromBech32(redeemAddrStr)
				if err != nil {
					return err
				}
			}

			var operators []sdk.AccAddress
			operatorSlice := viper.GetStringSlice(FlagOperators)
			for _, operator := range operatorSlice {
				if len(operator) > 0 {
					operatorAddr, err := sdk.AccAddressFromBech32(operator)
					if err != nil {
						return err
					} else {
						operators = append(operators, operatorAddr)
					}
				}
			}

			moniker := viper.GetString(FlagMoniker)
			identity := viper.GetString(FlagIdentity)
			details := viper.GetString(FlagDetails)
			website := viper.GetString(FlagWebsite)

			var msg sdk.Msg
			msg = asset.NewMsgCreateGateway(
				identity, moniker, details, website, redeemAddr, owner, operators,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsGatewayCreate)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}
