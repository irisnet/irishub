package cli

import (
	"fmt"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/client/context"
)

// queryTokenFees retrieves the fees of issuance and minting for the specified symbol
func queryTokenFees(cliCtx context.CLIContext, symbol string) (asset.TokenFeesOutput, error) {
	params := asset.QueryTokenFeesParams{
		Symbol: symbol,
	}

	bz := cliCtx.Codec.MustMarshalJSON(params)

	route := fmt.Sprintf("custom/%s/fees/tokens", protocol.AssetRoute)
	res, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return asset.TokenFeesOutput{}, err
	}

	var out asset.TokenFeesOutput
	err = cliCtx.Codec.UnmarshalJSON(res, &out)
	return out, err
}
