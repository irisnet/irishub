package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

func TestParams(t *testing.T) {
	ctx, _, keeper := CreateTestInput(t, false, sdk.ZeroInt())
	expParams := types.DefaultParams()

	//check that the empty keeper loads the default
	resParams := keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))

	//modify a params, save, and retrieve
	expParams.MaxValidators = 777
	keeper.SetParams(ctx, expParams)
	resParams = keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))
}

func TestPool(t *testing.T) {
	ctx, _, keeper := CreateTestInput(t, false, sdk.ZeroInt())
	expPool := types.InitialBondedPool()

	//check that the empty keeper loads the default
	resPool := keeper.GetPool(ctx)
	require.True(t, expPool.Equal(resPool.BondedPool))

	//modify a params, save, and retrieve
	expPool.BondedTokens = sdk.NewDec(777)
	resPool.BondedPool = expPool
	keeper.SetPool(ctx, resPool)
	resPool = keeper.GetPool(ctx)
	require.True(t, expPool.Equal(resPool.BondedPool))
}
