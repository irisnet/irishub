package cli


import (
	"os"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	authcmd "github.com/irisnet/irishub/client/auth/cli"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/modules/profiling"
)

func GetCmdCreateProfiler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-profiler",
		Short: "Create a new profiler",
		Example: "iriscli profiling create-profiler --chain-id=<chain-id> --from=<key name> --fee=0.004iris " +
			"--profiler-address=<added address> --profiler-name=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
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
			msg := profiling.NewMsgAddProfiler(pAddr, fromAddr, name)
			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsProfilerAddress)
	cmd.Flags().AddFlagSet(FsProfilerName)
	return cmd
}