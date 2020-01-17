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
	app := createTestApp(nil, 0)

	cases := []struct {
		params types.Params
	}{
		{types.DefaultParams()},
		{types.NewParams(sdk.NewRat(5, 10))},
	}

	for _, tc := range cases {
		app.csk.SetParams(app.ctx, tc.params)

		feeParam := app.csk.GetParams(app.ctx)
		require.Equal(t, tc.params.Fee, feeParam.Fee)
	}
}

func TestAddAndRemoveLiquidity(t *testing.T) {
	total, _ := sdk.NewIntFromString("10000000000000000000")
	denom1 := "btc-min"
	denom2 := sdk.IrisAtto

	initCoins := sdk.NewCoins(sdk.NewCoin(denom1, total), sdk.NewCoin(denom2, total))
	app := createTestApp(initCoins, 1)
	sender := app.accounts[0].GetAddress()

	uniID, _ := types.GetUniID(denom1, denom2)

	btcAmt, _ := sdk.NewIntFromString("1")
	depositCoin := sdk.NewCoin(denom1, btcAmt)

	irisAmt, _ := sdk.NewIntFromString("10000000000000000000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, irisAmt, minReward, deadline.Unix(), sender)
	_, err := app.csk.HandleAddLiquidity(app.ctx, msg)
	//assert
	require.Nil(t, err)
	pool, existed := app.csk.GetPool(app.ctx, uniID)
	require.True(t, existed)
	require.Equal(t, "1btc-min,10000000000000000000iris-atto,10000000000000000000uni:btc-min", pool.Balance().String())
	senderBalances := app.ak.GetAccount(app.ctx, sender).GetCoins()
	require.Equal(t, "9999999999999999999btc-min,10000000000000000000uni:btc-min", senderBalances.String())

	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(sdk.NewInt(1), sdk.NewCoin("uni:btc-min", withdraw),
		sdk.NewInt(1), app.ctx.BlockHeader().Time.Unix(),
		sender)

	_, err = app.csk.HandleRemoveLiquidity(app.ctx, msgRemove)
	require.Nil(t, err)

	pools := app.csk.GetPools(app.ctx)
	require.Len(t, pools, 1)
	acc := app.ak.GetAccount(app.ctx, sender)
	require.True(t, existed)
	require.Equal(t, "", pools[0].Balance().String())
	require.Equal(t, "10000000000000000000btc-min,10000000000000000000iris-atto", acc.GetCoins().String())
}
