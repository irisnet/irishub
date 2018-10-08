package cli

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

const (
	FlagFileName   = "file-name"
	FlagTargetPath = "target-path"
	FlagTxHash     = "tx-hash"
)

func GetCmdSubmit(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "submit specified file",
		RunE: func(cmd *cobra.Command, args []string) error {

			// cliCtx := context.NewCLIContext().
			// 	WithCodec(cdc).
			// 	WithLogger(os.Stdout).
			// 	WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			// txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
			// 	WithCliCtx(cliCtx)

			// from, err := cliCtx.GetFromAddress()
			// if err != nil {
			// 	return err
			// }

			// // TO DO
			// fileNameStr := viper.GetString(FlagFileName)

			// var msg sdk.Msg

			//msg = record.NewMsgRecord("this should be ipfs hash", fileNameStr, from)

			//return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})

			return nil
		},
	}

	cmd.Flags().String(FlagFileName, "", "")
	cmd.Flags().String(FlagTargetPath, "", "tx hash")

	return cmd
}
