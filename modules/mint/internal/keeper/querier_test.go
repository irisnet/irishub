package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keep "github.com/irisnet/irishub/modules/mint/internal/keeper"
	"github.com/irisnet/irishub/modules/mint/internal/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	app, ctx := createTestApp(true)
	querier := keep.NewQuerier(app.MintKeeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)
}

func TestQueryParams(t *testing.T) {
	app, ctx := createTestApp(true)
	querier := keep.NewQuerier(app.MintKeeper)

	var params types.Params

	res, sdkErr := querier(ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, sdkErr)

	err := app.Codec().UnmarshalJSON(res, &params)
	require.NoError(t, err)

	require.Equal(t, app.MintKeeper.GetParamSet(ctx), params)
}
