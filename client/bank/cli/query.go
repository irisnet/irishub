package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
)

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
// nolint: unparam
func GetAccountCmd(cdc *codec.Codec, decoder auth.AccountDecoder) *cobra.Command {
	return &cobra.Command{
		Use:     "account [address]",
		Short:   "Query account balance",
		Example: "iriscli bank account <account address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// find the key to look up the account
			addrString := args[0]

			addr, err := sdk.AccAddressFromBech32(addrString)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(decoder)

			if err := cliCtx.EnsureAccountExistsFromAddr(addr); err != nil {
				return err
			}

			acc, err := cliCtx.GetAccount(addr)
			if err != nil {
				return err
			}

			accountRes, err := bank.ConvertAccountCoin(cliCtx, acc)
			if err != nil {
				return err
			}

			output, err := codec.MarshalJSONIndent(cdc, accountRes)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}

// GetCmdQueryCoinType performs coin type query
func GetCmdQueryCoinType(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "coin-type [coin_name]",
		Short:   "query coin type",
		Example: "iriscli bank coin-type iris",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.GetCoinType(args[0])
			if err != nil {
				return err
			}
			output, err := codec.MarshalJSONIndent(cdc, res)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	return cmd
}

// GetCmdQueryTokenStats performs token statistic query
func GetCmdQueryTokenStats(cdc *codec.Codec, queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-stats",
		Short:   "query token statistics",
		Example: "iriscli bank token-stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resToken, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, auth.QueryTokenStats), nil)
			if err != nil {
				return err
			}

			var tokenStats bank.TokenStats
			err = cdc.UnmarshalJSON(resToken, &tokenStats)
			if err != nil {
				return err
			}

			resPool, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.StakeRoute, stake.QueryPool), nil)
			if err != nil {
				return err
			}
			var poolStatus types.PoolStatus
			err = cdc.UnmarshalJSON(resPool, &poolStatus)
			if err != nil {
				return err
			}

			tokenStats.BondedToken = poolStatus.BondedTokens
			res, err := codec.MarshalJSONIndent(cdc, tokenStats)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}
