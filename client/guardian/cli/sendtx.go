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
		Example: "iriscli guardian add-profiler --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--address=<added address> --description=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)
			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			paStr := viper.GetString(FlagAddress)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			description := viper.GetString(FlagDescription)
			if len(description) == 0 {
				return fmt.Errorf("must use --description flag")
			}
			msg := guardian.NewMsgAddProfiler(description, pAddr, fromAddr)
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	cmd.MarkFlagRequired(FlagAddress)
	cmd.MarkFlagRequired(FlagDescription)
	return cmd
}

func GetCmdDeleteProfiler(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-profiler",
		Short: "Delete a profiler",
		Example: "iriscli guardian delete-profiler --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--address=<deleted address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)
			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			paStr := viper.GetString(FlagAddress)
			if len(paStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			pAddr, err := sdk.AccAddressFromBech32(paStr)
			if err != nil {
				return err
			}
			msg := guardian.NewMsgDeleteProfiler(pAddr, fromAddr)
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	cmd.MarkFlagRequired(FlagAddress)
	return cmd
}

func GetCmdCreateTrustee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-trustee",
		Short: "Add a new trustee",
		Example: "iriscli guardian add-trustee --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--address=<added address> --description=<name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)
			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			taStr := viper.GetString(FlagAddress)
			if len(taStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			tAddr, err := sdk.AccAddressFromBech32(taStr)
			if err != nil {
				return err
			}
			description := viper.GetString(FlagDescription)
			if len(description) == 0 {
				return fmt.Errorf("must use --description flag")
			}
			msg := guardian.NewMsgAddTrustee(description, tAddr, fromAddr)
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAddGuardian)
	return cmd
}

func GetCmdDeleteTrustee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-trustee",
		Short: "Delete a trustee",
		Example: "iriscli guardian delete-trustee --chain-id=<chain-id> --from=<key-name> --fee=0.3iris " +
			"--address=<deleted address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)
			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			taStr := viper.GetString(FlagAddress)
			if len(taStr) == 0 {
				return fmt.Errorf("must use --address flag")
			}
			tAddr, err := sdk.AccAddressFromBech32(taStr)
			if err != nil {
				return err
			}
			msg := guardian.NewMsgDeleteTrustee(tAddr, fromAddr)
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsDeleteGuardian)
	return cmd
}
