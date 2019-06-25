package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
)

// queryGatewayFee retrieves the gateway creation fee for the specified moniker
func queryGatewayFee(cliCtx context.CLIContext, moniker string) (asset.GatewayFeeOutput, error) {
	params := asset.QueryGatewayFeeParams{
		Moniker: moniker,
	}

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return asset.GatewayFeeOutput{}, err
	}

	path := fmt.Sprintf("custom/%s/fees/gateways", protocol.AssetRoute)

	res, err := cliCtx.QueryWithData(path, bz)
	if err != nil {
		return asset.GatewayFeeOutput{}, err
	}

	var out asset.GatewayFeeOutput
	err = cliCtx.Codec.UnmarshalJSON(res, &out)
	if err != nil {
		return asset.GatewayFeeOutput{}, err
	}

	return out, nil
}

// queryTokenFees retrieves the fees of token issuance and minting for the specified id
func queryTokenFees(cliCtx context.CLIContext, tokenID string) (asset.TokenFeesOutput, error) {
	params := asset.QueryTokenFeesParams{
		ID: tokenID,
	}

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return asset.TokenFeesOutput{}, err
	}

	path := fmt.Sprintf("custom/%s/fees/tokens", protocol.AssetRoute)

	res, err := cliCtx.QueryWithData(path, bz)
	if err != nil {
		return asset.TokenFeesOutput{}, err
	}

	var out asset.TokenFeesOutput
	err = cliCtx.Codec.UnmarshalJSON(res, &out)
	if err != nil {
		return asset.TokenFeesOutput{}, err
	}

	return out, nil
}
