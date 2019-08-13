package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

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
//	require.Nil(t, keeper.HandleAddLiquidity(ctx, msgAdd))
//
//	poolAccout := keeper.ak.GetAccount(ctx, poolAddr)
//	acc := keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "1btc,10iris-atto,10swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "999btc,990iris-atto,10swap:btc:iris-atto", acc.GetCoins().String())
//
//	msgAdd1 := types.NewMsgAddLiquidity(sdk.Coin{Denom: "btc", Amount: sdk.NewInt(1)},
//		sdk.NewInt(3), sdk.NewInt(3), ctx.BlockHeader().Time,
//		accs[0].GetAddress())
//	require.Nil(t, keeper.HandleAddLiquidity(ctx, msgAdd1))
//
//	poolAccout = keeper.ak.GetAccount(ctx, poolAddr)
//	acc = keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "2btc,13iris-atto,13swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "998btc,987iris-atto,13swap:btc:iris-atto", acc.GetCoins().String())
//
//	require.Equal(t, "100btc,10iris-atto,10swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "900btc,990iris-atto,10swap:btc:iris-atto", acc.GetCoins().String())
//
//	require.Nil(t, keeper.HandleAddLiquidity(ctx, msgAdd))
//
//	poolAccout = keeper.ak.GetAccount(ctx, poolAddr)
//	acc = keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "200btc,20iris-atto,20swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "800btc,980iris-atto,20swap:btc:iris-atto", acc.GetCoins().String())
//
//	msgRemove := types.NewMsgRemoveLiquidity(sdk.Coin{Denom: "btc", Amount: sdk.NewInt(1)},
//		sdk.NewInt(3), sdk.NewInt(3), ctx.BlockHeader().Time,
//		accs[0].GetAddress())
//	require.Nil(t, keeper.HandleRemoveLiquidity(ctx, msgRemove))
//
//	poolAccout = keeper.ak.GetAccount(ctx, poolAddr)
//	acc = keeper.ak.GetAccount(ctx, accs[0].GetAddress())
//	require.Equal(t, "2btc,10iris-atto,10swap:btc:iris-atto", poolAccout.GetCoins().String())
//	require.Equal(t, "998btc,990iris-atto,10swap:btc:iris-atto", acc.GetCoins().String())
//}
