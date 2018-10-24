package cli

import (
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
)

// GetCmdUnrevoke implements the create unrevoke validator command.
func GetCmdUnrevoke(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unrevoke",
		Args:  cobra.ExactArgs(0),
		Short: "unrevoke validator previously revoked for downtime",
		Example: "iriscli stake unrevoke --to=<account address> --from <key name> --fee=0.004iris --chain-id=<chain-id>",
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

			msg := slashing.NewMsgUnrevoke(validatorAddr)

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
