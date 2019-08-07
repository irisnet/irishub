package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

const (
	reservePoolName = "swap:iris:btc"
)

// test that the module account gets created with an initial
// balance of zero coins.
func TestCreateReservePool(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	poolAcc := getReservePoolAddr(reservePoolName)
	moduleAcc := keeper.ak.GetAccount(ctx, poolAcc)
	require.Nil(t, moduleAcc)

	keeper.CreateReservePool(ctx, reservePoolName)
	moduleAcc = keeper.ak.GetAccount(ctx, poolAcc)
	require.NotNil(t, moduleAcc)
	require.Equal(t, true, moduleAcc.GetCoins().Empty(), "module account has non zero balance after creation")

	// attempt to recreate existing ModuleAccount
	require.Panics(t, func() { keeper.CreateReservePool(ctx, reservePoolName) })
}

// test that the params can be properly set and retrieved
func TestParams(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		params types.Params
	}{
		{types.DefaultParams()},
		{types.NewParams(sdk.NewRat(5, 10))},
	}

	for _, tc := range cases {
		keeper.SetParams(ctx, tc.params)

		feeParam := keeper.GetParams(ctx)
		require.Equal(t, tc.params.Fee, feeParam.Fee)
	}
}

// test that non existent reserve pool returns false and
// that balance is updated.
func TestGetReservePool(t *testing.T) {
	amt := sdk.NewInt(100)
	ctx, keeper, accs := createTestInput(t, amt, 1)

	poolAcc := getReservePoolAddr(reservePoolName)
	reservePool, found := keeper.GetReservePool(ctx, reservePoolName)
	require.False(t, found)

	keeper.CreateReservePool(ctx, reservePoolName)
	reservePool, found = keeper.GetReservePool(ctx, reservePoolName)
	require.True(t, found)

	keeper.bk.SendCoins(ctx, accs[0].GetAddress(), poolAcc, sdk.Coins{sdk.NewCoin(sdk.IrisAtto, amt)})
	reservePool, found = keeper.GetReservePool(ctx, reservePoolName)
	reservePool, found = keeper.GetReservePool(ctx, reservePoolName)
	require.True(t, found)
	require.Equal(t, amt, reservePool.AmountOf(sdk.IrisAtto))
}

//func TestKeeper_UpdateLiquidity(t *testing.T) {
//	ctx, keeper, accs := createTestInput(t, sdk.NewInt(1000), 1)
//
//	liquidityName := "swap:btc:iris-atto"
//	poolAddr := getReservePoolAddr(liquidityName)
//
//	// init liquidity
//	msgAdd := types.NewMsgAddLiquidity(sdk.Coin{Denom: "btc", Amount: sdk.NewInt(1)},
//		sdk.NewInt(10), sdk.NewInt(10), ctx.BlockHeader().Time,
//		accs[0].GetAddress())
//
//	require.Nil(t, keeper.AddLiquidity(ctx, msgAdd))
//
//	poolAccout := keeper.ak.GetAccount(ctx, poolAddr)
//	acc := keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "1btc,10iris-atto,10swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "999btc,990iris-atto,10swap:btc:iris-atto", acc.GetCoins().String())
//
//	msgAdd1 := types.NewMsgAddLiquidity(sdk.Coin{Denom: "btc", Amount: sdk.NewInt(1)},
//		sdk.NewInt(3), sdk.NewInt(3), ctx.BlockHeader().Time,
//		accs[0].GetAddress())
//	require.Nil(t, keeper.AddLiquidity(ctx, msgAdd1))
//
//	poolAccout = keeper.ak.GetAccount(ctx, poolAddr)
//	acc = keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "2btc,13iris-atto,13swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "998btc,987iris-atto,13swap:btc:iris-atto", acc.GetCoins().String())
//
//	require.Equal(t, "100btc,10iris-atto,10swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "900btc,990iris-atto,10swap:btc:iris-atto", acc.GetCoins().String())
//
//	require.Nil(t, keeper.AddLiquidity(ctx, msgAdd))
//
//	poolAccout = keeper.ak.GetAccount(ctx, poolAddr)
//	acc = keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "200btc,20iris-atto,20swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "800btc,980iris-atto,20swap:btc:iris-atto", acc.GetCoins().String())
//
//	msgRemove := types.NewMsgRemoveLiquidity(sdk.Coin{Denom: "btc", Amount: sdk.NewInt(1)},
//		sdk.NewInt(3), sdk.NewInt(3), ctx.BlockHeader().Time,
//		accs[0].GetAddress())
//	require.Nil(t, keeper.RemoveLiquidity(ctx, msgRemove))
//
//	poolAccout = keeper.ak.GetAccount(ctx, poolAddr)
//	acc = keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "2btc,10iris-atto,10swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "998btc,990iris-atto,10swap:btc:iris-atto", acc.GetCoins().String())
//}
