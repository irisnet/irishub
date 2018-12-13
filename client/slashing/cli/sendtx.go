package cli

import (
	"os"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	authcmd "github.com/irisnet/irishub/client/auth/cli"
	"github.com/irisnet/irishub/modules/slashing"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/spf13/cobra"
)

// GetCmdUnrevoke implements the create unrevoke validator command.
func GetCmdUnrevoke(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unjail",
		Args:    cobra.ExactArgs(0),
		Short:   "unjail validator previously jailed for downtime",
		Example: "iriscli stake unjail --from <key name> --fee=0.004iris --chain-id=<chain-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			validatorAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(validatorAddr))

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
