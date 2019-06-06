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
		Short:   "Query all the gateways of the specified owner",
		Example: "iriscli asset query-gateways --owner=<gateway owner>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner, err := sdk.AccAddressFromBech32(viper.GetString(FlagOwner))
			if err != nil {
				return err
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

			var gateways []asset.Gateway
			err = cdc.UnmarshalJSON(res, &gateways)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(gateways)
		},
	}

	cmd.Flags().String(FlagOwner, "", "the owner address to be queried")
	cmd.MarkFlagRequired(FlagOwner)

	return cmd
}
