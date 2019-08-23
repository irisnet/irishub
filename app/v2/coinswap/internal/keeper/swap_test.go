package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	sdk "github.com/irisnet/irishub/types"
)

var (
	native = sdk.IrisAtto
)

func TestGetUniId(t *testing.T) {

	cases := []struct {
		name         string
		denom1       string
		denom2       string
		expectResult string
		expectPass   bool
	}{
		{"denom1 is native", native, "btc-min", "u-btc", true},
		{"denom2 is native", "btc-min", native, "u-btc", true},
		{"denom1 equals denom2", "btc-min", "btc-min", "u-btc", false},
		{"neither denom is native", "eth-min", "btc-min", "u-btc", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uniId, err := types.GetUniId(tc.denom1, tc.denom2)
			if tc.expectPass {
				require.Equal(t, tc.expectResult, uniId)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

type Data struct {
	delta sdk.Int
	x     sdk.Int
	y     sdk.Int
	fee   sdk.Rat
}
type SwapCase struct {
	data   Data
	expect sdk.Int
}

func TestGetInputPrice(t *testing.T) {
	var datas = []SwapCase{
		{
			data:   Data{delta: sdk.NewInt(100), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewRat(3, 1000)},
			expect: sdk.NewInt(90),
		},
		{
			data:   Data{delta: sdk.NewInt(200), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewRat(3, 1000)},
			expect: sdk.NewInt(166),
		},
		{
			data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewRat(3, 1000)},
			expect: sdk.NewInt(230),
		},
		{
			data:   Data{delta: sdk.NewInt(1000), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewRat(3, 1000)},
			expect: sdk.NewInt(499),
		},
		{
			data:   Data{delta: sdk.NewInt(1000), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.ZeroRat()},
			expect: sdk.NewInt(500),
		},
	}
	for _, tcase := range datas {
		data := tcase.data
		actual := getInputPrice(data.delta, data.x, data.y, data.fee)
		fmt.Println(fmt.Sprintf("expect:%s,actual:%s", tcase.expect.String(), actual.String()))
		require.Equal(t, tcase.expect, actual)
	}
}

func TestGetOutputPrice(t *testing.T) {
	var datas = []SwapCase{
		{
			data:   Data{delta: sdk.NewInt(100), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewRat(3, 1000)},
			expect: sdk.NewInt(112),
		},
		{
			data:   Data{delta: sdk.NewInt(200), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewRat(3, 1000)},
			expect: sdk.NewInt(251),
		},
		{
			data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewRat(3, 1000)},
			expect: sdk.NewInt(430),
		},
		{
			data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.ZeroRat()},
			expect: sdk.NewInt(429),
		},
	}
	for _, tcase := range datas {
		data := tcase.data
		actual := getOutputPrice(data.delta, data.x, data.y, data.fee)
		fmt.Println(fmt.Sprintf("expect:%s,actual:%s", tcase.expect.String(), actual.String()))
		require.Equal(t, tcase.expect, actual)
	}
}

func TestKeeperSwap(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewInt(100000000), 1)
	sender := accs[0].GetAddress()
	denom1 := "btc-min"
	denom2 := sdk.IrisAtto
	uniId, _ := types.GetUniId(denom1, denom2)
	reservePoolAddr := getReservePoolAddr(uniId)

	depositCoin := sdk.NewCoin("btc-min", sdk.NewInt(1000))
	depositAmount := sdk.NewInt(1000)
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, depositAmount, minReward, deadline.Unix(), sender)
	_, err := keeper.HandleAddLiquidity(ctx, msg)

	//assert
	require.Nil(t, err)
	reservePoolBalances := keeper.ak.GetAccount(ctx, reservePoolAddr).GetCoins()
	require.Equal(t, "1000btc-min,1000iris-atto,1000u-btc-min", reservePoolBalances.String())
	senderBlances := keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "99999000btc-min,99999000iris-atto,1000u-btc-min", senderBlances.String())

	inputCoin := sdk.NewCoin("btc-min", sdk.NewInt(100))
	outputCoin := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1))

	input := types.Input{
		Address: sender,
		Coin:    inputCoin,
	}

	output := types.Output{
		Coin: outputCoin,
	}

	deadline1 := time.Now().Add(1 * time.Minute)
	msg1 := types.NewMsgSwapOrder(input, output, deadline1.Unix(), false)
	_, err = keeper.HandleSwap(ctx, msg1)
	require.Nil(t, err)

	reservePoolBalances = keeper.ak.GetAccount(ctx, reservePoolAddr).GetCoins()
	require.Equal(t, "1100btc-min,910iris-atto,1000u-btc-min", reservePoolBalances.String())
	senderBlances = keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "99998900btc-min,99999090iris-atto,1000u-btc-min", senderBlances.String())
}
