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
		fmt.Println(fmt.Sprintf("expect:%s,actual:%s", tcase.expect.String(), actual.String()))
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
		fmt.Println(fmt.Sprintf("expect:%s,actual:%s", tcase.expect.String(), actual.String()))
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

	// first swap buy order
	err := suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBlances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
	suite.Equal(fmt.Sprintf("900%s,1112%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal(fmt.Sprintf("99999100%s,99998888%s,1000%s", denomBTC, denomStandard, unidenomBTC), senderBlances.String())

	// second swap buy order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBlances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
	suite.Equal(fmt.Sprintf("800%s,1252%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal(fmt.Sprintf("99999200%s,99998748%s,1000%s", denomBTC, denomStandard, unidenomBTC), senderBlances.String())

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
	senderBlances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
	suite.Equal(fmt.Sprintf("446%s,2252%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal(fmt.Sprintf("99999554%s,99997748%s,1000%s", denomBTC, denomStandard, unidenomBTC), senderBlances.String())

	// second swap sell order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBlances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender)
	suite.Equal(fmt.Sprintf("310%s,3252%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal(fmt.Sprintf("99999690%s,99996748%s,1000%s", denomBTC, denomStandard, unidenomBTC), senderBlances.String())
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

	// first swap buy order
	err := suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBTCBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrBTC)
	reservePoolETHBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrETH)
	sender1Blances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender1)
	suite.Equal(fmt.Sprintf("1127%s,888%s", denomBTC, denomStandard), reservePoolBTCBalances.String())
	suite.Equal(fmt.Sprintf("900%s,1112%s", denomETH, denomStandard), reservePoolETHBalances.String())
	suite.Equal(fmt.Sprintf("99998873%s,100%s,99999000%s,1000%s", denomBTC, denomETH, denomStandard, unidenomBTC), sender1Blances.String())

	// second swap buy order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBTCBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrBTC)
	reservePoolETHBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrETH)
	sender1Blances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender1)
	suite.Equal(fmt.Sprintf("1339%s,748%s", denomBTC, denomStandard), reservePoolBTCBalances.String())
	suite.Equal(fmt.Sprintf("800%s,1252%s", denomETH, denomStandard), reservePoolETHBalances.String())
	suite.Equal(fmt.Sprintf("99998661%s,200%s,99999000%s,1000%s", denomBTC, denomETH, denomStandard, unidenomBTC), sender1Blances.String())

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
	sender2Blances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sender2)
	suite.Equal(fmt.Sprintf("696%s,1442%s", denomBTC, denomStandard), reservePoolBTCBalances.String())
	suite.Equal(fmt.Sprintf("1800%s,558%s", denomETH, denomStandard), reservePoolETHBalances.String())
	suite.Equal(fmt.Sprintf("643%s,99998000%s,99999000%s,1000%s", denomBTC, denomETH, denomStandard, unidenomETH), sender2Blances.String())

	// second swap sell order
	err = suite.app.CoinswapKeeper.Swap(suite.ctx, msg)
	suite.NoError(err)
	reservePoolBTCBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrBTC)
	reservePoolETHBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddrETH)
	sender2Blances = suite.app.BankKeeper.GetAllBalances(suite.ctx, sender2)
	suite.Equal(fmt.Sprintf("613%s,1640%s", denomBTC, denomStandard), reservePoolBTCBalances.String())
	suite.Equal(fmt.Sprintf("2800%s,360%s", denomETH, denomStandard), reservePoolETHBalances.String())
	suite.Equal(fmt.Sprintf("726%s,99997000%s,99999000%s,1000%s", denomBTC, denomETH, denomStandard, unidenomETH), sender2Blances.String())
}

func createReservePool(suite *TestSuite, denom string) (sdk.AccAddress, sdk.AccAddress) {
	amountInit, _ := sdk.NewIntFromString("100000000")
	addrSender := sdk.AccAddress(getRandomString(20))
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.BankKeeper.AddCoins(
		suite.ctx,
		addrSender,
		sdk.NewCoins(
			sdk.NewCoin(denomStandard, amountInit),
			sdk.NewCoin(denom, amountInit),
		),
	)

	uniDenom, err := types.GetUniDenomFromDenom(denom)
	suite.NoError(err)
	reservePoolAddr := types.GetReservePoolAddr(uniDenom)

	depositAmt, _ := sdk.NewIntFromString("1000")
	depositCoin := sdk.NewCoin(denom, depositAmt)

	standardAmt, _ := sdk.NewIntFromString("1000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender.String())
	_, err = suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	moduleAccountBalances := suite.app.BankKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, reservePoolAddr)
	senderBlances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender)
	suite.Equal("1000", moduleAccountBalances.AmountOf(uniDenom).String())
	suite.Equal(fmt.Sprintf("1000%s,1000%s", denom, denomStandard), reservePoolBalances.String())
	suite.Equal(fmt.Sprintf("99999000%s,99999000%s,1000%s", denom, denomStandard, uniDenom), senderBlances.String())
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
