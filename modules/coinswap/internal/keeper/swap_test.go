package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/coinswap/internal/keeper"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

var (
	native = "iris-atto"
)

func (suite *KeeperTestSuite) TestGetUniId() {

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
		suite.T().Run(tc.name, func(t *testing.T) {
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
	fee   sdk.Dec
}
type SwapCase struct {
	data   Data
	expect sdk.Int
}

func (suite *KeeperTestSuite) TestGetInputPrice() {
	var datas = []SwapCase{
		{
			data:   Data{delta: sdk.NewInt(100), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
			expect: sdk.NewInt(90),
		},
		{
			data:   Data{delta: sdk.NewInt(200), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
			expect: sdk.NewInt(166),
		},
		{
			data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
			expect: sdk.NewInt(230),
		},
		{
			data:   Data{delta: sdk.NewInt(1000), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
			expect: sdk.NewInt(499),
		},
		{
			data:   Data{delta: sdk.NewInt(1000), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.ZeroDec()},
			expect: sdk.NewInt(500),
		},
	}
	for _, tcase := range datas {
		data := tcase.data
		actual := keeper.GetInputPrice(data.delta, data.x, data.y, data.fee)
		fmt.Println(fmt.Sprintf("expect:%s,actual:%s", tcase.expect.String(), actual.String()))
		require.Equal(suite.T(), tcase.expect, actual)
	}
}

func (suite *KeeperTestSuite) TestGetOutputPrice() {
	var datas = []SwapCase{
		{
			data:   Data{delta: sdk.NewInt(100), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
			expect: sdk.NewInt(112),
		},
		{
			data:   Data{delta: sdk.NewInt(200), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
			expect: sdk.NewInt(251),
		},
		{
			data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
			expect: sdk.NewInt(430),
		},
		{
			data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.ZeroDec()},
			expect: sdk.NewInt(429),
		},
	}
	for _, tcase := range datas {
		data := tcase.data
		actual := keeper.GetOutputPrice(data.delta, data.x, data.y, data.fee)
		fmt.Println(fmt.Sprintf("expect:%s,actual:%s", tcase.expect.String(), actual.String()))
		require.Equal(suite.T(), tcase.expect, actual)
	}
}

func (suite *KeeperTestSuite) TestKeeperSwap() {
	sender, reservePoolAddr, err, reservePoolBalances, senderBlances := createReservePool(suite)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(100))
	inputCoin := sdk.NewCoin(native, sdk.NewInt(1000))

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
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg1)
	require.Nil(suite.T(), err)
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	require.Equal(suite.T(), "900btc-min,1112iris-atto,1000uni:btc-min", reservePoolBalances.String())
	senderBlances = suite.app.AccountKeeper.GetAccount(suite.ctx, sender).GetCoins()
	require.Equal(suite.T(), "99999100btc-min,99998888iris-atto,1000uni:btc-min", senderBlances.String())

	// second swap
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg1)
	require.Nil(suite.T(), err)
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	require.Equal(suite.T(), "800btc-min,1252iris-atto,1000uni:btc-min", reservePoolBalances.String())
	senderBlances = suite.app.AccountKeeper.GetAccount(suite.ctx, sender).GetCoins()
	require.Equal(suite.T(), "99999200btc-min,99998748iris-atto,1000uni:btc-min", senderBlances.String())

	// third swap
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg1)
	require.Nil(suite.T(), err)
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	require.Equal(suite.T(), "700btc-min,1432iris-atto,1000uni:btc-min", reservePoolBalances.String())
}

func createReservePool(suite *KeeperTestSuite) (sdk.AccAddress, sdk.AccAddress, sdk.Error, sdk.Coins, sdk.Coins) {
	amountInit, _ := sdk.NewIntFromString("100000000")
	addrSender := sdk.AccAddress([]byte("addrSender"))
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.BankKeeper.SetCoins(
		suite.ctx,
		addrSender,
		sdk.NewCoins(
			sdk.NewCoin(denomIrisAtto, amountInit),
			sdk.NewCoin(denomBTCMin, amountInit),
		),
	)

	denom1 := "btc-min"
	denom2 := native
	uniId, _ := types.GetUniId(denom1, denom2)
	reservePoolAddr := keeper.GetReservePoolAddr(uniId)

	btcAmt, _ := sdk.NewIntFromString("1000")
	depositCoin := sdk.NewCoin("btc-min", btcAmt)

	irisAmt, _ := sdk.NewIntFromString("1000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, irisAmt, minReward, deadline.Unix(), addrSender)
	err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	//assert
	require.Nil(suite.T(), err)
	reservePoolBalances := suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	require.Equal(suite.T(), "1000btc-min,1000iris-atto,1000uni:btc-min", reservePoolBalances.String())
	senderBlances := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()
	require.Equal(suite.T(), "99999000btc-min,99999000iris-atto,1000uni:btc-min", senderBlances.String())
	return addrSender, reservePoolAddr, err, reservePoolBalances, senderBlances
}

func (suite *KeeperTestSuite) TestTradeInputForExactOutput() {
	sender, poolAddr, _, poolBalances, senderBlances := createReservePool(suite)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(100))
	inputCoin := sdk.NewCoin(native, sdk.NewInt(100000))
	input := types.Input{
		Address: sender,
		Coin:    inputCoin,
	}
	output := types.Output{
		Coin: outputCoin,
	}

	initSupplyOutput := poolBalances.AmountOf(outputCoin.Denom)
	maxCnt := int(initSupplyOutput.Quo(outputCoin.Amount).Int64())

	for i := 1; i < 100; i++ {
		amt, err := suite.app.CoinswapKeeper.TradeInputForExactOutput(suite.ctx, input, output)
		if i == maxCnt {
			require.NotNil(suite.T(), err)
			break
		}
		ifNil(suite, err)

		bought := sdk.NewCoins(outputCoin)
		sold := sdk.NewCoins(sdk.NewCoin(native, amt))

		pb := poolBalances.Add(sold).Sub(bought)
		sb := senderBlances.Add(bought).Sub(sold)

		assertResult(suite, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func (suite *KeeperTestSuite) TestTradeExactInputForOutput() {
	sender, poolAddr, _, poolBalances, senderBlances := createReservePool(suite)

	outputCoin := sdk.NewCoin("btc-min", sdk.NewInt(0))
	inputCoin := sdk.NewCoin(native, sdk.NewInt(100))
	input := types.Input{
		Address: sender,
		Coin:    inputCoin,
	}
	output := types.Output{
		Coin: outputCoin,
	}

	for i := 1; i < 1000; i++ {
		amt, err := suite.app.CoinswapKeeper.TradeExactInputForOutput(suite.ctx, input, output)
		ifNil(suite, err)

		sold := sdk.NewCoins(inputCoin)
		bought := sdk.NewCoins(sdk.NewCoin("btc-min", amt))

		pb := poolBalances.Add(sold).Sub(bought)
		sb := senderBlances.Add(bought).Sub(sold)

		assertResult(suite, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func assertResult(suite *KeeperTestSuite, reservePoolAddr, sender sdk.AccAddress, expectPoolBalance, expectSenderBalance sdk.Coins) {
	reservePoolBalances := suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	require.Equal(suite.T(), expectPoolBalance.String(), reservePoolBalances.String())
	senderBlances := suite.app.AccountKeeper.GetAccount(suite.ctx, sender).GetCoins()
	require.Equal(suite.T(), expectSenderBalance.String(), senderBlances.String())
}

func ifNil(suite *KeeperTestSuite, err sdk.Error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	require.Nil(suite.T(), err, msg)
}
