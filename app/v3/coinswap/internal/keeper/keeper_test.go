package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// test that the params can be properly set and retrieved
func TestParams(t *testing.T) {
	ctx, keeper, _, _ := createTestInput(t, sdk.NewInt(0), 0)

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

func TestKeeper_UpdateLiquidity(t *testing.T) {
	total, _ := sdk.NewIntFromString("10000000000000000000")
	ctx, keeper, ak, accs := createTestInput(t, total, 1)
	sender := accs[0].GetAddress()
	denom1 := "btc-min"
	denom2 := sdk.IrisAtto

	uniID, _ := types.GetUniID(denom1, denom2)

	btcAmt, _ := sdk.NewIntFromString("1")
	depositCoin := sdk.NewCoin(denom1, btcAmt)

	irisAmt, _ := sdk.NewIntFromString("10000000000000000000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, irisAmt, minReward, deadline.Unix(), sender)
	_, err := keeper.HandleAddLiquidity(ctx, msg)
	//assert
	require.Nil(t, err)
	reservePoolBalances, existed := keeper.GetPool(ctx, uniID)
	require.True(t, existed)
	require.Equal(t, "1btc-min,10000000000000000000iris-atto,10000000000000000000uni:btc-min", reservePoolBalances.String())
	senderBalances := ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "9999999999999999999btc-min,10000000000000000000uni:btc-min", senderBalances.String())

	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(sdk.NewInt(1), sdk.NewCoin("uni:btc-min", withdraw),
		sdk.NewInt(1), ctx.BlockHeader().Time.Unix(),
		sender)

	_, err = keeper.HandleRemoveLiquidity(ctx, msgRemove)
	require.Nil(t, err)

	reservePoolBalances, existed = keeper.GetPool(ctx, uniID)
	acc := ak.GetAccount(ctx, sender)
	require.True(t, existed)
	require.Equal(t, "", reservePoolBalances.String())
	require.Equal(t, "10000000000000000000btc-min,10000000000000000000iris-atto", acc.GetCoins().String())
}
