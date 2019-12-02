package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/modules/asset/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetQueryCmd returns the query commands for this channels
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the asset module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(client.GetCommands(
		GetCmdQueryToken(queryRoute, cdc),
		GetCmdQueryTokens(queryRoute, cdc),
		GetCmdQueryFee(queryRoute, cdc),
	)...)

	return queryCmd
}

// GetCmdQueryToken implements the query token command.
func GetCmdQueryToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-token",
		Short:   "Query details of a token",
		Example: "iriscli asset query-token <token-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryTokenParams{
				TokenId: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryToken), bz)
			if err != nil {
				return err
			}

			var token types.FungibleToken
			err = cdc.UnmarshalJSON(res, &token)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(token)
		},
	}

	return cmd
}

// GetCmdQueryTokens implements the query tokens command.
func GetCmdQueryTokens(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-tokens",
		Short:   "Query details of a group of tokens",
		Example: "iriscli asset query-tokens --source=<native|gateway|external> --gateway=<gateway_moniker> --owner=<address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QueryTokensParams{
				Source:  viper.GetString(FlagSource),
				Gateway: viper.GetString(FlagGateway),
				Owner:   viper.GetString(FlagOwner),
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTokens), bz)
			if err != nil {
				return err
			}

			var tokens types.Tokens
			err = cdc.UnmarshalJSON(res, &tokens)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}

	cmd.Flags().AddFlagSet(FsTokensQuery)

	return cmd
}

// GetCmdQueryFee implements the query asset related fees command.
func GetCmdQueryFee(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-fee",
		Short:   "Query the asset related fees",
		Example: "iriscli asset query-fee --token=<token id>",
		PreRunE: preQueryFeeCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query token fees
			tokenID := viper.GetString(FlagToken)
			if err := types.CheckTokenID(tokenID); err != nil {
				return err
			}

			fees, err := queryTokenFees(cliCtx, queryRoute, tokenID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	cmd.Flags().AddFlagSet(FsFeeQuery)

	return cmd
}

// preQueryFeeCmd is used to check if the specified flags are valid
func preQueryFeeCmd(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	if flags.Changed(FlagGateway) && flags.Changed(FlagToken) {
		return fmt.Errorf("only one flag is allowed among the gateway and token")
	} else if !flags.Changed(FlagGateway) && !flags.Changed(FlagToken) {
		return fmt.Errorf("must specify the gateway or token to be queried")
	}

	return nil
}
