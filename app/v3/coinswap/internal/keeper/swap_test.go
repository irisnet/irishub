package keeper

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
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
		{"denom1 is native", native, "btc-min", "uni:btc", true},
		{"denom2 is native", "btc-min", native, "uni:btc", true},
		{"denom1 equals denom2", "btc-min", "btc-min", "uni:btc", false},
		{"neither denom is native", "eth-min", "btc-min", "uni:btc", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uniId, err := types.GetUniID(tc.denom1, tc.denom2)
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
	amount := sdk.NewInt(100000000)
	btcToken := sdk.NewCoin("btc-min", amount)
	irisToken := sdk.NewCoin(sdk.IrisAtto, amount)
	app := createTestApp(sdk.NewCoins(btcToken, irisToken).Sort(), 1)

	sender := app.accounts[0].GetAddress()
	uniID := createReservePool(app, btcToken.Denom)

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
	_, err := app.csk.HandleSwap(app.ctx, msg1)
	require.Nil(t, err)
	pool, existed := app.csk.GetPool(app.ctx, uniID)
	require.True(t, existed)
	require.Equal(t, "900btc-min,1112iris-atto,1000uni:btc-min", pool.Balance().String())
	senderBalances := app.ak.GetAccount(app.ctx, sender).GetCoins()
	require.Equal(t, "99999100btc-min,99998888iris-atto,1000uni:btc-min", senderBalances.String())

	// second swap
	_, err = app.csk.HandleSwap(app.ctx, msg1)
	require.Nil(t, err)
	pool, existed = app.csk.GetPool(app.ctx, uniID)
	require.True(t, existed)
	require.Equal(t, "800btc-min,1252iris-atto,1000uni:btc-min", pool.Balance().String())
	senderBalances = app.ak.GetAccount(app.ctx, sender).GetCoins()
	require.Equal(t, "99999200btc-min,99998748iris-atto,1000uni:btc-min", senderBalances.String())

	// third swap
	_, err = app.csk.HandleSwap(app.ctx, msg1)
	require.Nil(t, err)
	pool, existed = app.csk.GetPool(app.ctx, uniID)
	require.True(t, existed)
	require.Equal(t, "700btc-min,1432iris-atto,1000uni:btc-min", pool.Balance().String())
}

func TestKeeperDoubleSwap(t *testing.T) {
	amount := sdk.NewInt(100000000)
	btcToken := sdk.NewCoin("btc-min", amount)
	ethToken := sdk.NewCoin("eth-min", amount)
	irisToken := sdk.NewCoin(sdk.IrisAtto, amount)
	app := createTestApp(sdk.NewCoins(btcToken, irisToken, ethToken).Sort(), 1)

	sender := app.accounts[0].GetAddress()
	ctx := app.ctx

	btcUniID := createReservePool(app, btcToken.Denom)
	ethUniID := createReservePool(app, ethToken.Denom)

	senderBalances := app.ak.GetAccount(ctx, sender).GetCoins()
	fmt.Println(senderBalances.String())

	uniDenomBTC, _ := types.GetUniDenom(btcUniID)
	uniDenomETH, _ := types.GetUniDenom(ethUniID)

	msg := types.NewMsgSwapOrder(
		types.Input{Coin: sdk.NewCoin(btcToken.Denom, sdk.NewInt(1000)), Address: sender},
		types.Output{Coin: sdk.NewCoin(ethToken.Denom, sdk.NewInt(100))},
		time.Now().Add(1*time.Minute).Unix(),
		true,
	)

	// first swap buy order
	_, err := app.csk.HandleSwap(ctx, msg)
	require.True(t, err == nil)

	poolBTC, existed := app.csk.GetPool(app.ctx, btcUniID)
	require.True(t, existed)

	poolETH, existed := app.csk.GetPool(app.ctx, ethUniID)
	require.True(t, existed)

	poolBTCBalances := poolBTC.Balance()
	poolETHBalances := poolETH.Balance()
	senderBalances = app.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, fmt.Sprintf("1127%s,888%s,1000%s", btcToken.Denom, sdk.IrisAtto, uniDenomBTC), poolBTCBalances.String())
	require.Equal(t, fmt.Sprintf("900%s,1112%s,1000%s", ethToken.Denom, sdk.IrisAtto, uniDenomETH), poolETHBalances.String())
	require.Equal(t, fmt.Sprintf("99998873%s,99999100%s,99998000%s,1000%s,1000%s", btcToken.Denom, ethToken.Denom, sdk.IrisAtto, uniDenomBTC, uniDenomETH), senderBalances.String())

	// second swap buy order
	_, err = app.csk.HandleSwap(ctx, msg)
	require.NoError(t, err)

	poolBTC, existed = app.csk.GetPool(app.ctx, btcUniID)
	require.True(t, existed)

	poolETH, existed = app.csk.GetPool(app.ctx, ethUniID)
	require.True(t, existed)

	poolBTCBalances = poolBTC.Balance()
	poolETHBalances = poolETH.Balance()
	senderBalances = app.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, fmt.Sprintf("1339%s,748%s,1000%s", btcToken.Denom, sdk.IrisAtto, uniDenomBTC), poolBTCBalances.String())
	require.Equal(t, fmt.Sprintf("800%s,1252%s,1000%s", ethToken.Denom, sdk.IrisAtto, uniDenomETH), poolETHBalances.String())
	require.Equal(t, fmt.Sprintf("99998661%s,99999200%s,99998000%s,1000%s,1000%s", btcToken.Denom, ethToken.Denom, sdk.IrisAtto, uniDenomBTC, uniDenomETH), senderBalances.String())

	// swap sell order msg
	msg = types.NewMsgSwapOrder(
		types.Input{Coin: sdk.NewCoin(ethToken.Denom, sdk.NewInt(100)), Address: sender},
		types.Output{Coin: sdk.NewCoin(btcToken.Denom, sdk.NewInt(80))},
		time.Now().Add(1*time.Minute).Unix(),
		false,
	)

	// first swap sell order
	_, err = app.csk.HandleSwap(ctx, msg)
	require.True(t, err == nil)
	poolBTC, existed = app.csk.GetPool(app.ctx, btcUniID)
	require.True(t, existed)

	poolETH, existed = app.csk.GetPool(app.ctx, ethUniID)
	require.True(t, existed)

	poolBTCBalances = poolBTC.Balance()
	poolETHBalances = poolETH.Balance()
	senderBalances = app.ak.GetAccount(ctx, sender).GetCoins()

	require.Equal(t, fmt.Sprintf("1131%s,886%s,1000%s", btcToken.Denom, sdk.IrisAtto, uniDenomBTC), poolBTCBalances.String())
	require.Equal(t, fmt.Sprintf("900%s,1114%s,1000%s", ethToken.Denom, sdk.IrisAtto, uniDenomETH), poolETHBalances.String())
	require.Equal(t, fmt.Sprintf("99998869%s,99999100%s,99998000%s,1000%s,1000%s", btcToken.Denom, ethToken.Denom, sdk.IrisAtto, uniDenomBTC, uniDenomETH), senderBalances.String())

	// second swap sell order
	_, err = app.csk.HandleSwap(ctx, msg)
	require.True(t, err == nil)
	poolBTC, existed = app.csk.GetPool(app.ctx, btcUniID)
	require.True(t, existed)

	poolETH, existed = app.csk.GetPool(app.ctx, ethUniID)
	require.True(t, existed)

	poolBTCBalances = poolBTC.Balance()
	poolETHBalances = poolETH.Balance()
	senderBalances = app.ak.GetAccount(ctx, sender).GetCoins()
	fmt.Println(senderBalances.String())

	require.Equal(t, fmt.Sprintf("1006%s,997%s,1000%s", btcToken.Denom, sdk.IrisAtto, uniDenomBTC), poolBTCBalances.String())
	require.Equal(t, fmt.Sprintf("1000%s,1003%s,1000%s", ethToken.Denom, sdk.IrisAtto, uniDenomETH), poolETHBalances.String())
	require.Equal(t, fmt.Sprintf("99998994%s,99999000%s,99998000%s,1000%s,1000%s", btcToken.Denom, ethToken.Denom, sdk.IrisAtto, uniDenomBTC, uniDenomETH), senderBalances.String())
}

func TestTradeInputForExactOutput(t *testing.T) {
	amount := sdk.NewInt(100000000)
	btcToken := sdk.NewCoin("btc-min", amount)
	irisToken := sdk.NewCoin(sdk.IrisAtto, amount)
	app := createTestApp(sdk.NewCoins(btcToken, irisToken).Sort(), 1)
	sender := app.accounts[0]
	uniID := createReservePool(app, btcToken.Denom)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(100))
	inputCoin := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100000))
	input := types.Input{
		Address: sender.GetAddress(),
		Coin:    inputCoin,
	}
	output := types.Output{
		Coin: outputCoin,
	}

	pool, existed := app.csk.GetPool(app.ctx, uniID)
	require.True(t, existed)

	initSupplyOutput := pool.BalanceOf(outputCoin.Denom)
	maxCnt := int(initSupplyOutput.Div(outputCoin.Amount).Int64())

	balance := app.ak.GetAccount(app.ctx, sender.GetAddress()).GetCoins()
	for i := 1; i < 100; i++ {
		amt, err := app.csk.tradeInputForExactOutput(app.ctx, input, output)
		if i == maxCnt {
			require.NotNil(t, err)
			break
		}
		ifNil(t, err)

		bought := sdk.NewCoins(outputCoin)
		sold := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, amt))

		pool.Add(sold).Sub(bought)
		sb := balance.Add(bought).Sub(sold)

		assertResult(t, app.csk, app.ak, app.ctx, uniID, sender.GetAddress(), pool.Balance(), sb)

		balance = sb
	}
}

func TestTradeExactInputForOutput(t *testing.T) {
	amount := sdk.NewInt(100000000)
	btcToken := sdk.NewCoin("btc-min", amount)
	irisToken := sdk.NewCoin(sdk.IrisAtto, amount)
	app := createTestApp(sdk.NewCoins(btcToken, irisToken).Sort(), 1)
	sender := app.accounts[0]
	uniID := createReservePool(app, btcToken.Denom)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(0))
	inputCoin := sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100))
	input := types.Input{
		Address: sender.GetAddress(),
		Coin:    inputCoin,
	}
	output := types.Output{
		Coin: outputCoin,
	}

	pool, existed := app.csk.GetPool(app.ctx, uniID)
	require.True(t, existed)

	balance := app.ak.GetAccount(app.ctx, sender.GetAddress()).GetCoins()
	for i := 1; i < 1000; i++ {
		amt, err := app.csk.tradeExactInputForOutput(app.ctx, input, output)
		ifNil(t, err)

		sold := sdk.NewCoins(inputCoin)
		bought := sdk.NewCoins(sdk.NewCoin("btc-min", amt))

		pool.Add(sold).Sub(bought)
		sb := balance.Add(bought).Sub(sold)

		assertResult(t, app.csk, app.ak, app.ctx, uniID, sender.GetAddress(), pool.Balance(), sb)
		balance = sb
	}
}

func createReservePool(app TestApp, denom1 string) string {
	btcAmt, _ := sdk.NewIntFromString("1000")
	irisAmt, _ := sdk.NewIntFromString("1000")
	coin1 := sdk.NewCoin(denom1, btcAmt)
	coin2 := sdk.NewCoin(sdk.IrisAtto, irisAmt)
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	account := app.accounts[0]
	msg := types.NewMsgAddLiquidity(coin1, coin2.Amount, minReward, deadline.Unix(), account.GetAddress())
	_, _ = app.csk.HandleAddLiquidity(app.ctx, msg)

	uniID, _ := types.GetUniID(denom1, sdk.IrisAtto)
	return uniID
}

func assertResult(t *testing.T, keeper Keeper, ak auth.AccountKeeper, ctx sdk.Context, uniID string, sender sdk.AccAddress, expectPoolBalance, expectSenderBalance sdk.Coins) {
	pool, existed := keeper.GetPool(ctx, uniID)
	require.True(t, existed)
	require.Equal(t, expectPoolBalance.String(), pool.Balance().String())
	senderBalances := ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, expectSenderBalance.String(), senderBalances.String())
}

func ifNil(t *testing.T, err sdk.Error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	require.Nil(t, err, msg)
}
