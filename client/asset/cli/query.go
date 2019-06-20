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

// GetCmdQueryAsset implements the query asset command.
func GetCmdQueryAsset(cdc *codec.Codec) *cobra.Command {
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

// GetCmdQueryGateway implements the query gateway command.
func GetCmdQueryGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-gateway",
		Short:   "Query details of a gateway of the given moniker",
		Example: "iriscli asset query-gateway --moniker=<gateway moniker>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			moniker := viper.GetString(FlagMoniker)
			if len(moniker) < asset.MinimumGatewayMonikerSize || len(moniker) > asset.MaximumGatewayMonikerSize {
				return asset.ErrInvalidMoniker(asset.DefaultCodespace, fmt.Sprintf("the length of the moniker must be [%d,%d]", asset.MinimumGatewayMonikerSize, asset.MaximumGatewayMonikerSize))
			}

			if !asset.IsAlpha(moniker) {
				return asset.ErrInvalidMoniker(asset.DefaultCodespace, fmt.Sprintf("the moniker must contain only letters"))
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

// GetCmdQueryFee implements the query asset-related fees command.
func GetCmdQueryFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query-fee",
		Short:   "Query the asset-related fees",
		Example: "iriscli asset query-fee --subject=<gateway|token> --moniker=<gateway moniker> --id=<token id>",
		PreRunE: preQueryFeeCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// subject validity is checked in PreRunE
			subject := viper.GetString(FlagSubject)

			if subject == "gateway" {
				moniker := viper.GetString(FlagMoniker)
				if err := asset.ValidateMoniker(moniker); err != nil {
					return err
				}

				fee, err := queryGatewayFee(cliCtx, moniker)
				if err != nil {
					return err
				}

				fee.Fee = ConvertToNativeToken(fee.Fee)

				return cliCtx.PrintOutput(fee)

			} else {
				id := viper.GetString(FlagID)
				if err := asset.CheckAssetID(id); err != nil {
					return err
				}

				fees, err := queryTokenFees(cliCtx, id)
				if err != nil {
					return err
				}

				fees.IssueFee = ConvertToNativeToken(fees.IssueFee)
				fees.MintFee = ConvertToNativeToken(fees.MintFee)

				return cliCtx.PrintOutput(fees)
			}
		},
	}

	cmd.Flags().AddFlagSet(FsFeeQuery)
	cmd.MarkFlagRequired(FlagSubject)

	return cmd
}

// preQueryFeeCmd is used to check if the subject is valid and the corresponding flag to the subject is provided
func preQueryFeeCmd(cmd *cobra.Command, args []string) error {
	subject := viper.GetString(FlagSubject)

	if subject != "gateway" && subject != "token" {
		return fmt.Errorf("the subject must be gateway or token")
	}

	if subject == "gateway" {
		cmd.MarkFlagRequired(FlagMoniker)
	} else if subject == "token" {
		cmd.MarkFlagRequired(FlagID)
	}

	return nil
}
