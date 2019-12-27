package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

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

func TestKeeper_UpdateLiquidity(t *testing.T) {
	total, _ := sdk.NewIntFromString("10000000000000000000")
	ctx, keeper, accs := createTestInput(t, total, 1)
	sender := accs[0].GetAddress()
	denom1 := "btc-min"
	denom2 := sdk.IrisAtto
	uniId, _ := types.GetUniId(denom1, denom2)
	poolAddr := getReservePoolAddr(uniId)

	btcAmt, _ := sdk.NewIntFromString("1")
	depositCoin := sdk.NewCoin("btc-min", btcAmt)

	irisAmt, _ := sdk.NewIntFromString("10000000000000000000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, irisAmt, minReward, deadline.Unix(), sender)
	_, err := keeper.HandleAddLiquidity(ctx, msg)
	//assert
	require.Nil(t, err)
	reservePoolBalances := keeper.ak.GetAccount(ctx, poolAddr).GetCoins()
	require.Equal(t, "1btc-min,10000000000000000000iris-atto,10000000000000000000uni:btc-min", reservePoolBalances.String())
	senderBlances := keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "9999999999999999999btc-min,10000000000000000000uni:btc-min", senderBlances.String())

	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(sdk.NewInt(1), sdk.NewCoin("uni:btc-min", withdraw),
		sdk.NewInt(1), ctx.BlockHeader().Time.Unix(),
		sender)

	_, err = keeper.HandleRemoveLiquidity(ctx, msgRemove)
	require.Nil(t, err)

	poolAccout := keeper.ak.GetAccount(ctx, poolAddr)
	acc := keeper.ak.GetAccount(ctx, sender)
	require.Equal(t, "", poolAccout.GetCoins().String())
	require.Equal(t, "10000000000000000000btc-min,10000000000000000000iris-atto", acc.GetCoins().String())
}
