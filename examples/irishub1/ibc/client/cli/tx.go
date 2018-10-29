package cli

import (
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/examples/irishub1/ibc"
	"os"
)

// IBC transfer command
func IBCGetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// build the message
			msg := ibc.NewIBCGetMsg(from)

			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}


// IBC transfer command
func IBCSetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "set",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := context.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// build the message
			msg := ibc.NewIBCSetMsg(from)

			cliCtx.PrintResponse = true
			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}

