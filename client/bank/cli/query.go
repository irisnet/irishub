package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/client/bank"
)

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
// nolint: unparam
func GetAccountCmd(storeName string, cdc *wire.Codec, decoder auth.AccountDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Example: "iriscli bank account <account address>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			addr := args[0]

			key, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(decoder)

			if err := cliCtx.EnsureAccountExistsFromAddr(key); err != nil {
				return err
			}

			acc, err := cliCtx.GetAccount(key)
			if err != nil {
				return err
			}

			accountRes, err := bank.ConvertAccountCoin(cliCtx, acc)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, accountRes)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}

// GetCmdQueryCoinType performs coin type query
func GetCmdQueryCoinType(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coin-type [coin_name]",
		Short: "query coin type",
		Example: "iriscli bank coin-type iris",
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
