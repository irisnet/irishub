package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// submit switch msg
func GetCmdSubmitSwitch(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-switch",
		Short: "Submit a switch msg for a upgrade propsal",
		RunE: func(cmd *cobra.Command, args []string) error {
			title := viper.GetString(flagTitle)
			proposalID := viper.GetInt64(flagProposalID)

			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from/to address
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := upgrade.NewMsgSwitch(title, proposalID, from)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txCtx, cliCtx, []sdk.Msg{msg})
			}

			cliCtx.PrintResponse = true
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of switch")
	cmd.Flags().String(flagProposalID, "", "proposalID of upgrade proposal")

	return cmd
}
