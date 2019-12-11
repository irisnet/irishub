package keeper_test

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/coinswap/internal/keeper"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

func (suite *KeeperTestSuite) TestGetUniId() {
	cases := []struct {
		name         string
		denom1       string
		denom2       string
		expectResult string
		expectPass   bool
	}{
		{"denom1 is denomStandard", denomStandard, denomBTC, unidenomBTC, true},
		{"denom2 is denomStandard", denomETH, denomStandard, unidenomETH, true},
		{"denom1 equals denom2", denomBTC, denomBTC, unidenomBTC, false},
		{"neither denom is denomStandard", denomETH, denomBTC, unidenomBTC, false},
	}

	for _, tc := range cases {
		uniDenom, err := types.GetUniDenomFromDenoms(tc.denom1, tc.denom2)
		if tc.expectPass {
			suite.Equal(tc.expectResult, uniDenom)
		} else {
			suite.NotNil(err)
		}
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
		suite.Equal(tcase.expect, actual)
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
		suite.Equal(tcase.expect, actual)
	}
}

func (suite *KeeperTestSuite) TestKeeperSwap() {
	sender, reservePoolAddr, err, reservePoolBalances, senderBlances := createReservePool(suite)
	moduleAccount := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

	outputCoin := sdk.NewCoin("btc", sdk.NewInt(100))
	inputCoin := sdk.NewCoin(denomStandard, sdk.NewInt(1000))

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
	suite.Nil(err)
	moduleAccountBalances := moduleAccount.GetCoins()
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	suite.Equal(fmt.Sprintf("900%s,1112%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal("1000", moduleAccountBalances.AmountOf(unidenomBTC).String())
	senderBlances = suite.app.AccountKeeper.GetAccount(suite.ctx, sender).GetCoins()
	suite.Equal(fmt.Sprintf("99999100%s,99998888%s,1000%s", denomBTC, denomStandard, unidenomBTC), senderBlances.String())

	// second swap
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg1)
	suite.Nil(err)
	moduleAccountBalances = moduleAccount.GetCoins()
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	suite.Equal(fmt.Sprintf("800%s,1252%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal("1000", moduleAccountBalances.AmountOf(unidenomBTC).String())
	senderBlances = suite.app.AccountKeeper.GetAccount(suite.ctx, sender).GetCoins()
	suite.Equal(fmt.Sprintf("99999200%s,99998748%s,1000%s", denomBTC, denomStandard, unidenomBTC), senderBlances.String())

	// third swap
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg1)
	suite.Nil(err)
	moduleAccountBalances = moduleAccount.GetCoins()
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	suite.Equal(fmt.Sprintf("700%s,1432%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal("1000", moduleAccountBalances.AmountOf(unidenomBTC).String())
}

func createReservePool(suite *KeeperTestSuite) (sdk.AccAddress, sdk.AccAddress, sdk.Error, sdk.Coins, sdk.Coins) {
	amountInit, _ := sdk.NewIntFromString("100000000")
	addrSender := sdk.AccAddress([]byte("addrSender"))
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.BankKeeper.SetCoins(
		suite.ctx,
		addrSender,
		sdk.NewCoins(
			sdk.NewCoin(denomStandard, amountInit),
			sdk.NewCoin(denomBTC, amountInit),
		),
	)

	denom1 := denomBTC
	denom2 := denomStandard
	uniDenom, _ := types.GetUniDenomFromDenoms(denom1, denom2)
	reservePoolAddr := keeper.GetReservePoolAddr(uniDenom)

	btcAmt, _ := sdk.NewIntFromString("1000")
	depositCoin := sdk.NewCoin(denomBTC, btcAmt)

	standardAmt, _ := sdk.NewIntFromString("1000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender)
	err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	//assert
	suite.Nil(err)
	moduleAccount := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	moduleAccountBalances := moduleAccount.GetCoins()
	reservePoolBalances := suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	suite.Equal(fmt.Sprintf("1000%s,1000%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal("1000", moduleAccountBalances.AmountOf(unidenomBTC).String())
	senderBlances := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()
	suite.Equal(fmt.Sprintf("99999000%s,99999000%s,1000%s", denomBTC, denomStandard, unidenomBTC), senderBlances.String())
	return addrSender, reservePoolAddr, err, reservePoolBalances, senderBlances
}

func (suite *KeeperTestSuite) TestTradeInputForExactOutput() {
	sender, poolAddr, _, poolBalances, senderBlances := createReservePool(suite)

	outputCoin := sdk.NewCoin(denomBTC, sdk.NewInt(100))
	inputCoin := sdk.NewCoin(denomStandard, sdk.NewInt(100000))
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
			suite.NotNil(err)
			break
		}
		ifNil(suite, err)

		bought := sdk.NewCoins(outputCoin)
		sold := sdk.NewCoins(sdk.NewCoin(denomStandard, amt))

		pb := poolBalances.Add(sold).Sub(bought)
		sb := senderBlances.Add(bought).Sub(sold)

		assertResult(suite, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func (suite *KeeperTestSuite) TestTradeExactInputForOutput() {
	sender, poolAddr, _, poolBalances, senderBlances := createReservePool(suite)

	outputCoin := sdk.NewCoin(denomBTC, sdk.NewInt(0))
	inputCoin := sdk.NewCoin(denomStandard, sdk.NewInt(100))
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
		bought := sdk.NewCoins(sdk.NewCoin(denomBTC, amt))

		pb := poolBalances.Add(sold).Sub(bought)
		sb := senderBlances.Add(bought).Sub(sold)

		assertResult(suite, poolAddr, sender, pb, sb)

		poolBalances = pb
		senderBlances = sb
	}
}

func assertResult(suite *KeeperTestSuite, reservePoolAddr, sender sdk.AccAddress, expectPoolBalance, expectSenderBalance sdk.Coins) {
	reservePoolBalances := suite.app.AccountKeeper.GetAccount(suite.ctx, reservePoolAddr).GetCoins()
	suite.Equal(expectPoolBalance.String(), reservePoolBalances.String())
	senderBlances := suite.app.AccountKeeper.GetAccount(suite.ctx, sender).GetCoins()
	suite.Equal(expectSenderBalance.String(), senderBlances.String())
}

func ifNil(suite *KeeperTestSuite, err sdk.Error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	suite.Nil(err, msg)
}
