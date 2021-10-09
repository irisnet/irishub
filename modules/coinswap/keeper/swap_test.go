package keeper_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/coinswap/keeper"
	"github.com/irisnet/irismod/modules/coinswap/types"
)

func TestSwapSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
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

func (suite *TestSuite) TestGetInputPrice() {
	var datas = []SwapCase{{
		data:   Data{delta: sdk.NewInt(100), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
		expect: sdk.NewInt(90),
	}, {
		data:   Data{delta: sdk.NewInt(200), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
		expect: sdk.NewInt(166),
	}, {
		data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
		expect: sdk.NewInt(230),
	}, {
		data:   Data{delta: sdk.NewInt(1000), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
		expect: sdk.NewInt(499),
	}, {
		data:   Data{delta: sdk.NewInt(1000), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.ZeroDec()},
		expect: sdk.NewInt(500),
	}}
	for _, tcase := range datas {
		data := tcase.data
		actual := keeper.GetInputPrice(data.delta, data.x, data.y, data.fee)
		suite.Equal(tcase.expect, actual)
	}
}

func (suite *TestSuite) TestGetOutputPrice() {
	var datas = []SwapCase{{
		data:   Data{delta: sdk.NewInt(100), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
		expect: sdk.NewInt(112),
	}, {
		data:   Data{delta: sdk.NewInt(200), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
		expect: sdk.NewInt(251),
	}, {
		data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.NewDecWithPrec(3, 3)},
		expect: sdk.NewInt(430),
	}, {
		data:   Data{delta: sdk.NewInt(300), x: sdk.NewInt(1000), y: sdk.NewInt(1000), fee: sdk.ZeroDec()},
		expect: sdk.NewInt(429),
	}}
	for _, tcase := range datas {
		data := tcase.data
		actual := keeper.GetOutputPrice(data.delta, data.x, data.y, data.fee)
		suite.Equal(tcase.expect, actual)
	}
}

func (suite *TestSuite) TestSwap() {
	sender, reservePoolAddr := createReservePool(suite, denomBTC)

	// swap buy order msg
	msg := types.NewMsgSwapOrder(
		types.Input{Coin: sdk.NewCoin(denomStandard, sdk.NewInt(1000)), Address: sender.String()},
		types.Output{Coin: sdk.NewCoin(denomBTC, sdk.NewInt(100)), Address: sender.String()},
		time.Now().Add(1*time.Minute).Unix(),
		true,
	)

	poolId := types.GetPoolId(denomBTC)
	pool, has := suite.app.CoinswapKeeper.GetPool(suite.ctx, poolId)
	suite.Require().True(has)

	lptDenom := pool.LptDenom

	// first swap buy order
	err := suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)

	expCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 900),
		sdk.NewInt64Coin(denomStandard, 1112),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 99999100),
		sdk.NewInt64Coin(denomStandard, 99998888),
		sdk.NewInt64Coin(lptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), senderBalances.Sort().String())

	// second swap buy order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 800),
		sdk.NewInt64Coin(denomStandard, 1252),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 99999200),
		sdk.NewInt64Coin(denomStandard, 99998748),
		sdk.NewInt64Coin(lptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), senderBalances.Sort().String())

	// swap sell order msg
	msg = types.NewMsgSwapOrder(
		types.Input{Coin: sdk.NewCoin(denomStandard, sdk.NewInt(1000)), Address: sender.String()},
		types.Output{Coin: sdk.NewCoin(denomBTC, sdk.NewInt(100)), Address: sender.String()},
		time.Now().Add(1*time.Minute).Unix(),
		false,
	)

	// first swap sell order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 446),
		sdk.NewInt64Coin(denomStandard, 2252),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 99999554),
		sdk.NewInt64Coin(denomStandard, 99997748),
		sdk.NewInt64Coin(lptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), senderBalances.Sort().String())

	// second swap sell order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
	suite.Equal(fmt.Sprintf("310%s,3252%s", denomBTC, denomStandard), reservePoolBalances.String())
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 310),
		sdk.NewInt64Coin(denomStandard, 3252),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 99999690),
		sdk.NewInt64Coin(denomStandard, 99996748),
		sdk.NewInt64Coin(lptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), senderBalances.Sort().String())
}

func (suite *TestSuite) TestDoubleSwap() {
	sender1, reservePoolAddrBTC := createReservePool(suite, denomBTC)
	sender2, reservePoolAddrETH := createReservePool(suite, denomETH)

	// swap buy order msg
	msg := types.NewMsgSwapOrder(
		types.Input{Coin: sdk.NewCoin(denomBTC, sdk.NewInt(1000)), Address: sender1.String()},
		types.Output{Coin: sdk.NewCoin(denomETH, sdk.NewInt(100)), Address: sender1.String()},
		time.Now().Add(1*time.Minute).Unix(),
		true,
	)

	poolId := types.GetPoolId(denomBTC)
	pool, has := suite.app.CoinswapKeeper.GetPool(suite.ctx, poolId)
	suite.Require().True(has)

	poolIdETH := types.GetPoolId(denomETH)
	poolETH, has := suite.app.CoinswapKeeper.GetPool(suite.ctx, poolIdETH)
	suite.Require().True(has)

	lptDenom := pool.LptDenom

	// first swap buy order
	err := suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBTCBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrBTC)
	reservePoolETHBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrETH)
	sender1Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender1)
	expCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 1127),
		sdk.NewInt64Coin(denomStandard, 888),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBTCBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomETH, 900),
		sdk.NewInt64Coin(denomStandard, 1112),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolETHBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 99998873),
		sdk.NewInt64Coin(denomETH, 100),
		sdk.NewInt64Coin(denomStandard, 99999000),
		sdk.NewInt64Coin(lptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	// second swap buy order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBTCBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrBTC)
	reservePoolETHBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrETH)
	sender1Balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender1)
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 1339),
		sdk.NewInt64Coin(denomStandard, 748),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBTCBalances.Sort().String())

	suite.Equal(fmt.Sprintf("800%s,1252%s", denomETH, denomStandard), reservePoolETHBalances.String())
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomETH, 800),
		sdk.NewInt64Coin(denomStandard, 1252),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolETHBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 99998661),
		sdk.NewInt64Coin(denomETH, 200),
		sdk.NewInt64Coin(denomStandard, 99999000),
		sdk.NewInt64Coin(lptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	// swap sell order msg
	msg = types.NewMsgSwapOrder(
		types.Input{Coin: sdk.NewCoin(denomETH, sdk.NewInt(1000)), Address: sender2.String()},
		types.Output{Coin: sdk.NewCoin(denomBTC, sdk.NewInt(80)), Address: sender2.String()},
		time.Now().Add(1*time.Minute).Unix(),
		false,
	)

	// first swap sell order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBTCBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrBTC)
	reservePoolETHBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrETH)
	sender2Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender2)
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 696),
		sdk.NewInt64Coin(denomStandard, 1442),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBTCBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomETH, 1800),
		sdk.NewInt64Coin(denomStandard, 558),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolETHBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 643),
		sdk.NewInt64Coin(denomETH, 99998000),
		sdk.NewInt64Coin(denomStandard, 99999000),
		sdk.NewInt64Coin(poolETH.LptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), sender2Balances.Sort().String())

	// second swap sell order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBTCBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrBTC)
	reservePoolETHBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrETH)
	sender2Balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender2)
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 613),
		sdk.NewInt64Coin(denomStandard, 1640),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBTCBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomETH, 2800),
		sdk.NewInt64Coin(denomStandard, 360),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolETHBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 726),
		sdk.NewInt64Coin(denomETH, 99997000),
		sdk.NewInt64Coin(denomStandard, 99999000),
		sdk.NewInt64Coin(poolETH.LptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), sender2Balances.Sort().String())
}

func createReservePool(suite *TestSuite, denom string) (sdk.AccAddress, sdk.AccAddress) {
	amountInit, _ := sdk.NewIntFromString("100000000")
	addrSender := sdk.AccAddress(getRandomString(20))
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)

	coins := sdk.NewCoins(
		sdk.NewCoin(denomStandard, amountInit),
		sdk.NewCoin(denom, amountInit),
	)

	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addrSender, coins)
	suite.NoError(err)

	depositAmt, _ := sdk.NewIntFromString("1000")
	depositCoin := sdk.NewCoin(denom, depositAmt)

	standardAmt, _ := sdk.NewIntFromString("1000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender.String())
	_, err = suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	poolId := types.GetPoolId(denom)
	pool, has := suite.app.CoinswapKeeper.GetPool(suite.ctx, poolId)
	suite.Require().True(has)
	reservePoolAddr := types.GetReservePoolAddr(pool.LptDenom)

	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBlances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender)
	suite.Equal("1000", suite.app.BankKeeper.GetSupply(suite.ctx, pool.LptDenom).Amount.String())

	expCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denom, 1000),
		sdk.NewInt64Coin(denomStandard, 1000),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denom, 99999000),
		sdk.NewInt64Coin(denomStandard, 99999000),
		sdk.NewInt64Coin(pool.LptDenom, 1000),
	)
	suite.Equal(expCoins.Sort().String(), senderBlances.Sort().String())
	return addrSender, reservePoolAddr
}

func (suite *TestSuite) TestTradeInputForExactOutput() {
	sender, poolAddr := createReservePool(suite, denomBTC)

	outputCoin := sdk.NewCoin(denomBTC, sdk.NewInt(100))
	inputCoin := sdk.NewCoin(denomStandard, sdk.NewInt(100000))
	input := types.Input{
		Address: sender.String(),
		Coin:    inputCoin,
	}
	output := types.Output{
		Address: sender.String(),
		Coin:    outputCoin,
	}

	poolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	senderBlances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)

	initSupplyOutput := poolBalances.AmountOf(outputCoin.Denom)
	maxCnt := int(initSupplyOutput.Quo(outputCoin.Amount).Int64())

	for i := 1; i < 100; i++ {
		amt, err := suite.app.CoinswapKeeper.TradeInputForExactOutput(suite.ctx, input, output)
		if i == maxCnt {
			suite.Error(err)
			break
		}
		suite.NoError(err)

		bought := sdk.NewCoins(outputCoin)
		sold := sdk.NewCoins(sdk.NewCoin(denomStandard, amt))

		pb := poolBalances.Add(sold...).Sub(bought)
		sb := senderBlances.Add(bought...).Sub(sold)

		assertResult(suite, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func (suite *TestSuite) TestTradeExactInputForOutput() {
	sender, poolAddr := createReservePool(suite, denomBTC)

	outputCoin := sdk.NewCoin(denomBTC, sdk.NewInt(0))
	inputCoin := sdk.NewCoin(denomStandard, sdk.NewInt(100))
	input := types.Input{
		Address: sender.String(),
		Coin:    inputCoin,
	}
	output := types.Output{
		Address: sender.String(),
		Coin:    outputCoin,
	}

	poolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	senderBlances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)

	for i := 1; i < 1000; i++ {
		amt, err := suite.app.CoinswapKeeper.TradeExactInputForOutput(suite.ctx, input, output)
		suite.NoError(err)

		sold := sdk.NewCoins(inputCoin)
		bought := sdk.NewCoins(sdk.NewCoin(denomBTC, amt))

		pb := poolBalances.Add(sold...).Sub(bought)
		sb := senderBlances.Add(bought...).Sub(sold)

		assertResult(suite, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func assertResult(suite *TestSuite, reservePoolAddr, sender sdk.AccAddress, expectPoolBalance, expectSenderBalance sdk.Coins) {
	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBlances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
	suite.Equal(expectPoolBalance.String(), reservePoolBalances.String())
	suite.Equal(expectSenderBalance.String(), senderBlances.String())
}

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
