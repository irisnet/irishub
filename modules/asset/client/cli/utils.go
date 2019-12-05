package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/irisnet/irishub/modules/asset/internal/types"
)

// queryTokenFees retrieves the fees of token issuance and minting for the specified id
func queryTokenFees(cliCtx context.CLIContext, queryRoute string, tokenID string) (types.TokenFeesOutput, error) {
	params := types.QueryTokenFeesParams{
		ID: tokenID,
	}

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return types.TokenFeesOutput{}, err
	}

	path := fmt.Sprintf("custom/%s/fees/tokens", queryRoute)

	res, _, err := cliCtx.QueryWithData(path, bz)
	if err != nil {
		return types.TokenFeesOutput{}, err
	}

	var out types.TokenFeesOutput
	err = cliCtx.Codec.UnmarshalJSON(res, &out)
	if err != nil {
		return types.TokenFeesOutput{}, err
	}

	return out, nil
}
