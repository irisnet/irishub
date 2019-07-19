package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryToken implements the query token command.
func GetCmdQueryToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-token",
		Short:   "Query details of a token",
		Example: "iriscli asset query-token <token-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := asset.QueryTokenParams{
				TokenId: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryToken), bz)
			if err != nil {
				return err
			}

			var token asset.FungibleToken
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
func GetCmdQueryTokens(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-tokens",
		Short:   "Query details of a group of tokens",
		Example: "iriscli asset query-tokens --source=<native|gateway|external> --gateway=<gateway_moniker> --owner=<address>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := asset.QueryTokensParams{
				Source:  viper.GetString(FlagSource),
				Gateway: viper.GetString(FlagGateway),
				Owner:   viper.GetString(FlagOwner),
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", protocol.AssetRoute, asset.QueryTokens), bz)
			if err != nil {
				return err
			}

			var tokens asset.Tokens
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

// GetCmdQueryGateway implements the query gateway command.
func GetCmdQueryGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-gateway",
		Short:   "Query details of a gateway of the given moniker",
		Example: "iriscli asset query-gateway --moniker=<gateway moniker>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			moniker := viper.GetString(FlagMoniker)
			if err := asset.ValidateMoniker(moniker); err != nil {
				return err
			}

			params := asset.QueryGatewayParams{
				Moniker: moniker,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/gateway", protocol.AssetRoute), bz)
			if err != nil {
				return err
			}

			var gateway asset.Gateway
			err = cdc.UnmarshalJSON(res, &gateway)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(gateway)
		},
	}

	cmd.Flags().String(FlagMoniker, "", "the unique name of the destination gateway")
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

// GetCmdQueryGateways implements the query gateways command.
func GetCmdQueryGateways(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-gateways",
		Short:   "Query all gateways with an optional owner",
		Example: "iriscli asset query-gateways --owner=<gateway owner>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				owner sdk.AccAddress
				err   error
			)

			ownerStr := viper.GetString(FlagOwner)
			if ownerStr != "" {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			params := asset.QueryGatewaysParams{
				Owner: owner,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/gateways", protocol.AssetRoute), bz)
			if err != nil {
				return err
			}

			var gateways asset.Gateways
			err = cdc.UnmarshalJSON(res, &gateways)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(gateways)
		},
	}

	cmd.Flags().String(FlagOwner, "", "the owner address to be queried")

	return cmd
}

// GetCmdQueryFee implements the query asset related fees command.
func GetCmdQueryFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-fee",
		Short:   "Query the asset related fees",
		Example: "iriscli asset query-fee --gateway=<gateway moniker>|--token=<token id>",
		PreRunE: preQueryFeeCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			flags := cmd.Flags()
			if flags.Changed(FlagGateway) {
				// query gateway fee
				moniker := viper.GetString(FlagGateway)
				if err := asset.ValidateMoniker(moniker); err != nil {
					return err
				}

				fee, err := queryGatewayFee(cliCtx, moniker)
				if err != nil {
					return err
				}

				return cliCtx.PrintOutput(fee)

			} else {
				// query token fees
				tokenID := viper.GetString(FlagToken)
				if err := asset.CheckTokenID(tokenID); err != nil {
					return err
				}

				fees, err := queryTokenFees(cliCtx, tokenID)
				if err != nil {
					return err
				}

				return cliCtx.PrintOutput(fees)
			}
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
