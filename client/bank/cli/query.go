package cli

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/auth"
	bankv1 "github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/stake"
	stakeTypes "github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/client/bank"
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

// GetCmdQueryTokenStats performs token statistic query
func GetCmdQueryTokenStats(cdc *codec.Codec, decoder auth.AccountDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token-stats [id]",
		Short:   "query token statistics",
		Example: "iriscli bank token-stats --id",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(decoder)
			if len(args) == 0 {
				//get the token-stats of iris
				resToken, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, bankv1.QueryTokenStats), nil)
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
				var poolStatus stakeTypes.PoolStatus
				err = cdc.UnmarshalJSON(resPool, &poolStatus)
				if err != nil {
					return err
				}

				tokenStats.BondedTokens = sdk.Coins{sdk.Coin{Denom: stakeTypes.StakeDenom, Amount: poolStatus.BondedTokens.TruncateInt()}}

				return cliCtx.PrintOutput(tokenStats)
			} else {
				//get the token-stats of other assets
				assetId := args[0]
				params := asset.QueryAssetParams{
					Asset: assetId,
				}
				bz, err := cdc.MarshalJSON(params)
				if err != nil {
					return err
				}

				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryAsset), bz)
				if err != nil {
					return err
				}

				var nAsset asset.Asset
				err = cdc.UnmarshalJSON(res, &nAsset)
				if err != nil {
					return err
				}
				var tokenStats bank.TokenStatsOfAsset

				//get loose token from asset
				looseToken := sdk.Coin{}
				looseToken.Denom = nAsset.GetDenom()
				looseToken.Amount = nAsset.GetTotalSupply()
				tokenStats.LooseToken = looseToken

				//get burned token from burnAddress
				burnedAcc, err := cliCtx.GetAccount(bankv1.BurnedCoinsAccAddr)
				if err != nil {
					return err
				}
				burnToken := sdk.Coin{nAsset.GetDenom(), burnedAcc.Coins.AmountOf(nAsset.GetDenom())}
				tokenStats.BurnedToken = burnToken
				return cliCtx.PrintOutput(tokenStats)
			}

		},
	}

	return cmd
}
