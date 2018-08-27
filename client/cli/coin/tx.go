package coin

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/irisnet/irishub/app"
)

func GetCmdQueryCoinType(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "types [coin_name]",
		Short: "query coin_type",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := app.NewContext().WithCodeC(cdc)
			res , err  := ctx.GetCoinType(args[0])
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