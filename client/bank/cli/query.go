package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
)

func GetCmdQueryCoinType(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coin-type [coin_name]",
		Short: "query coin type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.GetCoinType(args[0])
			if err != nil {
				return err
			}
			output, err := wire.MarshalJSONIndent(cdc, res)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	return cmd
}
