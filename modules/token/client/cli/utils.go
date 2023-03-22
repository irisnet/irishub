package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/token/types"
)

// queryTokenFees retrieves the fees of issuance and minting for the specified symbol
func queryTokenFees(cliCtx client.Context, symbol string) (types.QueryFeesResponse, error) {
	queryClient := types.NewQueryClient(cliCtx)

	resp, err := queryClient.Fees(context.Background(), &types.QueryFeesRequest{Symbol: symbol})
	if err != nil {
		return types.QueryFeesResponse{}, err
	}

	return *resp, err
}

// queryToken query token information
func queryToken(cliCtx client.Context, denom string) (types.TokenI, error) {
	queryClient := types.NewQueryClient(cliCtx)

	resp, err := queryClient.Token(context.Background(), &types.QueryTokenRequest{
		Denom: denom,
	})
	if err != nil {
		return nil, err
	}

	var evi types.TokenI
	err = cliCtx.InterfaceRegistry.UnpackAny(resp.Token, &evi)
	if err != nil {
		return nil, err
	}

	return evi, err
}

func parseCoin(cliCtx client.Context, denom string) (sdk.Coin, types.TokenI, error) {
	decCoin, err := sdk.ParseDecCoin(denom)
	if err != nil {
		return sdk.Coin{}, nil, err
	}

	token, err := queryToken(cliCtx, decCoin.Denom)
	if err != nil {
		return sdk.Coin{}, nil, err
	}

	coin, err := token.ToMinCoin(decCoin)
	if err != nil {
		return sdk.Coin{}, nil, err
	}
	return coin, token, err
}
