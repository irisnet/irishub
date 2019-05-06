package cli

import (
	"os"

	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

// GetSignCommand returns the sign command
func GetBroadcastCommand(codec *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcast <file>",
		Short: "Broadcast transactions generated offline",
		Long: `Broadcast transactions created with the --generate-only flag and signed with the sign command.
Read a transaction from <file> and broadcast it to a node.`,
		Example: "iriscli bank broadcast <file>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithLogger(os.Stdout).WithCodec(codec)
			stdTx, err := readAndUnmarshalStdTx(cliCtx.Codec, args[0])
			if err != nil {
				return
			}

			txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(stdTx)
			if err != nil {
				return
			}
			cliCtx.PrintResponse = true
			_, err = cliCtx.BroadcastTx(txBytes)
			return err
		},
	}

	return cmd
}
