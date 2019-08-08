package keeper

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/irisnet/irishub/types"
)

var (
	native = sdk.IrisAtto
)

func TestIsDoubleSwap(t *testing.T) {
	_, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		name         string
		denom1       string
		denom2       string
		isDoubleSwap bool
	}{
		{"denom1 is native", native, "btc", false},
		{"denom2 is native", "btc", native, false},
		{"neither denom is native", "eth", "btc", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doubleSwap := keeper.IsDoubleSwap(tc.denom1, tc.denom2)
			require.Equal(t, tc.isDoubleSwap, doubleSwap)
		})
	}
}

func TestGetReservePoolName(t *testing.T) {
	_, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		name         string
		denom1       string
		denom2       string
		expectResult string
		expectPass   bool
	}{
		{"denom1 is native", native, "btc", "s-btc", true},
		{"denom2 is native", "btc", native, "s-btc", true},
		{"denom1 equals denom2", "btc", "btc", "s-btc", false},
		{"neither denom is native", "eth", "btc", "s-btc", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reservePoolName, err := keeper.GetReservePoolName(tc.denom1, tc.denom2)
			if tc.expectPass {
				require.Equal(t, tc.expectResult, reservePoolName)
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
		actual := GetInputPrice(data.delta, data.x, data.y, data.fee)
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
		actual := GetOutputPrice(data.delta, data.x, data.y, data.fee)
		fmt.Println(fmt.Sprintf("expect:%s,actual:%s", tcase.expect.String(), actual.String()))
		require.Equal(t, tcase.expect, actual)
	}
}

//func TestSwapByInput(t *testing.T) {
//	ctx, keeper, accs := createTestInput(t, sdk.NewInt(100000000), 1)
//	sender := accs[0].GetAddress()
//
//	depositCoin := sdk.NewCoin("btc-min",sdk.NewInt(1000))
//	depositAmount := sdk.NewInt(1000)
//	minReward := sdk.NewInt(1)
//	deadline := time.Now().Add(1 * time.Minute)
//	msg := types.NewMsgAddLiquidity(depositCoin,depositAmount,minReward,deadline,sender)
//	err := keeper.AddLiquidity(ctx,msg)
//	require.Nil(t,err)
//
//	exactSoldCoin := sdk.NewCoin("btc-min",sdk.NewInt(100))
//	minExpect := sdk.NewCoin(sdk.IrisAtto,sdk.NewInt(100))
//	reward,err := keeper.SwapByInput(ctx,exactSoldCoin,minExpect,accs[0].GetAddress(),nil)
//	require.Nil(t,err)
//	require.Equal(t,sdk.NewInt(90),reward)
//}
//
//func TestSwapByOutput(t *testing.T) {
//
//}
