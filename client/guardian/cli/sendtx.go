package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCmdCreateProfiler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-profiler",
		Short: "Add a new profiler",
		Example: "iriscli guardian add-profiler --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--profiler-address=<added address> --profiler-name=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)
			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			paStr := viper.GetString(FlagProfilerAddress)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --profiler-address flag")
			}
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			name := viper.GetString(FlagProfilerName)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --profiler-name flag")
			}
			msg := guardian.NewMsgAddProfiler(pAddr, fromAddr, name)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsProfiler)
	cmd.MarkFlagRequired(FlagProfilerAddress)
	cmd.MarkFlagRequired(FlagProfilerName)
	return cmd
}
