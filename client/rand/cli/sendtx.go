package cli

import (
	"os"

	"github.com/irisnet/irishub/app/v1/rand"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdRequestRand implements the request-rand command
func GetCmdRequestRand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request-rand",
		Short:   "Request a random number",
		Example: "iriscli rand request-rand --block-interval=10",
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
				Consumer:      consumer,
				BlockInterval: uint64(viper.GetInt64(FlagBlockInterval)),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsRequestRand)

	return cmd
}
