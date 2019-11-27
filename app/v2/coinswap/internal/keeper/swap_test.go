package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
		{"denom1 is native", native, "btc-min", "uni:btc", true},
		{"denom2 is native", "btc-min", native, "uni:btc", true},
		{"denom1 equals denom2", "btc-min", "btc-min", "uni:btc", false},
		{"neither denom is native", "eth-min", "btc-min", "uni:btc", false},
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
	ctx, keeper, sender, reservePoolAddr, err, reservePoolBalances, senderBlances := createReservePool(t)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(100))
	inputCoin := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1000))

	input := types.Input{
		Address: sender,
		Coin:    inputCoin,
	}

	output := types.Output{
		Coin: outputCoin,
	}

	deadline1 := time.Now().Add(1 * time.Minute)
	msg1 := types.NewMsgSwapOrder(input, output, deadline1.Unix(), true)

	// first swap
	_, err = keeper.HandleSwap(ctx, msg1)
	require.Nil(t, err)
	reservePoolBalances = keeper.ak.GetAccount(ctx, reservePoolAddr).GetCoins()
	require.Equal(t, "900btc-min,1112iris-atto,1000uni:btc-min", reservePoolBalances.String())
	senderBlances = keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "99999100btc-min,99998888iris-atto,1000uni:btc-min", senderBlances.String())

	// second swap
	_, err = keeper.HandleSwap(ctx, msg1)
	require.Nil(t, err)
	reservePoolBalances = keeper.ak.GetAccount(ctx, reservePoolAddr).GetCoins()
	require.Equal(t, "800btc-min,1252iris-atto,1000uni:btc-min", reservePoolBalances.String())
	senderBlances = keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "99999200btc-min,99998748iris-atto,1000uni:btc-min", senderBlances.String())

	// third swap
	_, err = keeper.HandleSwap(ctx, msg1)
	require.Nil(t, err)
	reservePoolBalances = keeper.ak.GetAccount(ctx, reservePoolAddr).GetCoins()
	require.Equal(t, "700btc-min,1432iris-atto,1000uni:btc-min", reservePoolBalances.String())
}

func createReservePool(t *testing.T) (sdk.Context, Keeper, sdk.AccAddress, sdk.AccAddress, sdk.Error, sdk.Coins, sdk.Coins) {
	ctx, keeper, accs := createTestInput(t, sdk.NewInt(100000000), 1)
	sender := accs[0].GetAddress()
	denom1 := "btc-min"
	denom2 := sdk.IrisAtto
	uniId, _ := types.GetUniId(denom1, denom2)
	reservePoolAddr := getReservePoolAddr(uniId)

	btcAmt, _ := sdk.NewIntFromString("1000")
	depositCoin := sdk.NewCoin("btc-min", btcAmt)

	irisAmt, _ := sdk.NewIntFromString("1000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, irisAmt, minReward, deadline.Unix(), sender)
	_, err := keeper.HandleAddLiquidity(ctx, msg)
	//assert
	require.Nil(t, err)
	reservePoolBalances := keeper.ak.GetAccount(ctx, reservePoolAddr).GetCoins()
	require.Equal(t, "1000btc-min,1000iris-atto,1000uni:btc-min", reservePoolBalances.String())
	senderBlances := keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "99999000btc-min,99999000iris-atto,1000uni:btc-min", senderBlances.String())
	return ctx, keeper, sender, reservePoolAddr, err, reservePoolBalances, senderBlances
}

func TestTradeInputForExactOutput(t *testing.T) {
	ctx, keeper, sender, poolAddr, _, poolBalances, senderBlances := createReservePool(t)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(100))
	inputCoin := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100000))
	input := types.Input{
		Address: sender,
		Coin:    inputCoin,
	}
	output := types.Output{
		Coin: outputCoin,
	}

	initSupplyOutput := poolBalances.AmountOf(outputCoin.Denom)
	maxCnt := int(initSupplyOutput.Div(outputCoin.Amount).Int64())

	for i := 1; i < 100; i++ {
		amt, err := keeper.tradeInputForExactOutput(ctx, input, output)
		if i == maxCnt {
			require.NotNil(t, err)
			break
		}
		ifNil(t, err)

		bought := sdk.NewCoins(outputCoin)
		sold := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, amt))

		pb := poolBalances.Add(sold).Sub(bought)
		sb := senderBlances.Add(bought).Sub(sold)

		assertResult(t, keeper, ctx, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func TestTradeExactInputForOutput(t *testing.T) {
	ctx, keeper, sender, poolAddr, _, poolBalances, senderBlances := createReservePool(t)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(0))
	inputCoin := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100))
	input := types.Input{
		Address: sender,
		Coin:    inputCoin,
	}
	output := types.Output{
		Coin: outputCoin,
	}

	for i := 1; i < 1000; i++ {
		amt, err := keeper.tradeExactInputForOutput(ctx, input, output)
		ifNil(t, err)

		sold := sdk.NewCoins(inputCoin)
		bought := sdk.NewCoins(sdk.NewCoin("btc-min", amt))

		pb := poolBalances.Add(sold).Sub(bought)
		sb := senderBlances.Add(bought).Sub(sold)

		assertResult(t, keeper, ctx, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func assertResult(t *testing.T, keeper Keeper, ctx sdk.Context, reservePoolAddr, sender sdk.AccAddress, expectPoolBalance, expectSenderBalance sdk.Coins) {
	reservePoolBalances := keeper.ak.GetAccount(ctx, reservePoolAddr).GetCoins()
	require.Equal(t, expectPoolBalance.String(), reservePoolBalances.String())
	senderBlances := keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, expectSenderBalance.String(), senderBlances.String())
}

func ifNil(t *testing.T, err sdk.Error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	require.Nil(t, err, msg)
}
