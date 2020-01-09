package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/context"
)

// queryTokenFees retrieves the fees of token issuance and minting for the specified id
func queryTokenFees(cliCtx context.CLIContext, symbol string) (asset.TokenFeesOutput, error) {
	params := asset.QueryTokenFeesParams{
		Symbol: symbol,
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
