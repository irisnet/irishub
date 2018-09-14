package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
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

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
func GetAccountCmd(storeName string, cdc *wire.Codec, decoder auth.AccountDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			addr := args[0]

			key, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				return err
			}

			ctx := context.NewCLIContext()
			cliCtx := ctx.WithCodec(cdc).
				WithAccountDecoder(decoder)
			if err := cliCtx.EnsureAccountExistsFromAddr(key); err != nil {
				return err
			}

			acc, err := cliCtx.GetAccount(key)
			if err != nil {
				return err
			}

			coins := acc.GetCoins()
			var coins_str []string
			for _, coin := range coins {
				coinName, _ := types.GetCoinName(coin.String())
				ct, err := ctx.GetCoinType(coinName)
				if err != nil {
					return err
				}
				mainCoin, err := ct.Convert(coin.String(), ct.Name)
				if err != nil {
					return err
				}
				coins_str = append(coins_str, mainCoin)
			}
			acct := account{
				Address:       acc.GetAddress(),
				Coins:         coins_str,
				PubKey:        acc.GetPubKey(),
				AccountNumber: acc.GetAccountNumber(),
				Sequence:      acc.GetSequence(),
			}

			output, err := wire.MarshalJSONIndent(cdc, acct)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}

type account struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber int64          `json:"account_number"`
	Sequence      int64          `json:"sequence"`
}
