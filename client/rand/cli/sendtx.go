package cli

import (
	"os"

	"github.com/irisnet/irishub/app/v1/rand"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
)

// GetCmdRequestRand implements the request-rand command
func GetCmdRequestRand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request-rand",
		Short:   "request a random number",
		Example: "iriscli rand request-rand --from=<key-name> --chain-id=<chain-id> --fee=0.6iris",
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

			msg := rand.MsgRequestRand{
				Consumer: consumer,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
