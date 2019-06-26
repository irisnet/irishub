package cli

import (
	"fmt"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/auth"
	bankv1 "github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/stake"
	bankcli "github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
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

			if cliCtx.OutputFormat == "text" {
				coins, err := bankcli.ConvertToMainUnit(cliCtx, acc.GetCoins())
				if err != nil {
					return err
				}

				acc1, err := bankcli.ConvertAccountCoin(cliCtx, acc)
				if err != nil {
					return err
				}

				acc1.Coins = coins
				return cliCtx.PrintOutput(acc1)
			}
			return cliCtx.PrintOutput(acc)
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

			return cliCtx.PrintOutput(res)
		},
	}

	return cmd
}

func isIris(assetId string) bool {
	for _, ir := range sdk.IRIS.Units {
		if assetId == ir.Denom {
			return true
		}
	}
	return false
}

// GetCmdQueryTokenStats performs token statistic query
func GetCmdQueryTokenStats(cdc *codec.Codec, decoder auth.AccountDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-stats [id]",
		Short:   "query token statistics",
		Example: "iriscli bank token-stats iris",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(decoder)
			tokenId := ""
			if len(args) > 0 {
				tokenId = args[0]
			}
			params := asset.QueryTokenParams{
				TokenId: tokenId,
			}
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, bankv1.QueryTokenStats), bz)
			if err != nil {
				return err
			}

			var tokenStats bankv1.TokenStats
			err = cdc.UnmarshalJSON(res, &tokenStats)
			if err != nil {
				return err
			}

			// query bonded tokens for iris
			if tokenId == "" || tokenId == sdk.NativeTokenName {
				resPool, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.StakeRoute, stake.QueryPool), nil)
				if err != nil {
					return err
				}
				var poolStatus stake.PoolStatus
				err = cdc.UnmarshalJSON(resPool, &poolStatus)
				if err != nil {
					return err
				}

				tokenStats.BondedTokens = sdk.Coins{sdk.Coin{Denom: stake.BondDenom, Amount: poolStatus.BondedTokens.TruncateInt()}}
				tokenStats.TotalSupply = tokenStats.TotalSupply.Plus(tokenStats.LooseTokens.Plus(tokenStats.BondedTokens))
			}

			if cliCtx.OutputFormat == "text" {
				var tokenStats1 bankcli.TokenStats
				tokenStats1.LooseTokens, err = cliCtx.ConvertCoinToMainUnit(tokenStats.LooseTokens.String())
				if err != nil {
					return err
				}
				tokenStats1.BondedTokens, err = cliCtx.ConvertCoinToMainUnit(tokenStats.BondedTokens.String())
				if err != nil {
					return err
				}
				tokenStats1.BurnedTokens, err = cliCtx.ConvertCoinToMainUnit(tokenStats.BurnedTokens.String())
				if err != nil {
					return err
				}
				tokenStats1.TotalSupply, err = cliCtx.ConvertCoinToMainUnit(tokenStats.TotalSupply.String())
				if err != nil {
					return err
				}

				return cliCtx.PrintOutput(tokenStats1)
			}

			return cliCtx.PrintOutput(tokenStats)
		},
	}

	return cmd
}
