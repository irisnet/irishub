package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/coinswap/internal/types"
	"github.com/irisnet/irishub/simapp"
)

const (
	denomStandard = types.StandardDenom
	denomBTC      = "btc"
	denomETH      = "eth"
	unidenomBTC   = types.FormatUniABSPrefix + denomBTC
	unidenomETH   = types.FormatUniABSPrefix + denomETH
)

var (
	addrSender1 sdk.AccAddress
	addrSender2 sdk.AccAddress
)

// test that the params can be properly set and retrieved
type TestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.app = app
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestParams() {
	cases := []struct {
		params types.Params
	}{
		{types.DefaultParams()},
		{types.NewParams(sdk.NewDecWithPrec(5, 10), denomStandard)},
	}
	for _, tc := range cases {
		suite.app.CoinswapKeeper.SetParams(suite.ctx, tc.params)

		feeParam := suite.app.CoinswapKeeper.GetParams(suite.ctx)
		suite.Equal(tc.params.Fee, feeParam.Fee)
	}
}

func initVars(suite *TestSuite) {
	amountInitStandard, _ := sdk.NewIntFromString("30000000000000000000")
	amountInitBTC, _ := sdk.NewIntFromString("3000000000")

	addrSender1 = sdk.AccAddress("addrSender1")
	addrSender2 = sdk.AccAddress("addrSender2")
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender1)
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender2)
	_ = suite.app.BankKeeper.SetCoins(
		suite.ctx,
		addrSender1,
		sdk.NewCoins(
			sdk.NewCoin(denomStandard, amountInitStandard),
			sdk.NewCoin(denomBTC, amountInitBTC),
		),
	)
	_ = suite.app.BankKeeper.SetCoins(
		suite.ctx,
		addrSender2,
		sdk.NewCoins(
			sdk.NewCoin(denomStandard, amountInitStandard),
			sdk.NewCoin(denomBTC, amountInitBTC),
		),
	)
}

func (suite *TestSuite) TestLiquidity() {
	initVars(suite)

	// test add liquidity (poor not exist)
	uniDenom, _ := types.GetUniDenomFromDenoms(denomBTC, denomStandard)
	suite.Equal(uniDenom, unidenomBTC)
	poolAddr := types.GetReservePoolAddr(uniDenom)
	btcAmt, _ := sdk.NewIntFromString("100")
	standardAmt, _ := sdk.NewIntFromString("10000000000000000000")
	depositCoin := sdk.NewCoin(denomBTC, btcAmt)
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)

	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender1)
	err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	moduleAccountBalances := suite.app.SupplyKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances := suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr).GetCoins()
	sender1Blances := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender1).GetCoins()
	suite.Equal("10000000000000000000", moduleAccountBalances.AmountOf(unidenomBTC).String())
	suite.Equal(fmt.Sprintf("100%s,10000000000000000000%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal(fmt.Sprintf("2999999900%s,20000000000000000000%s,10000000000000000000%s", denomBTC, denomStandard, unidenomBTC), sender1Blances.String())

	// test add liquidity (poor exist)
	uniDenom, _ = types.GetUniDenomFromDenoms(denomBTC, denomStandard)
	suite.Equal(uniDenom, unidenomBTC)
	poolAddr = types.GetReservePoolAddr(uniDenom)
	btcAmt, _ = sdk.NewIntFromString("201")
	standardAmt, _ = sdk.NewIntFromString("20000000000000000000")
	depositCoin = sdk.NewCoin(denomBTC, btcAmt)
	minReward = sdk.NewInt(1)
	deadline = time.Now().Add(1 * time.Minute)

	msg = types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender2)
	err = suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	moduleAccountBalances = suite.app.SupplyKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr).GetCoins()
	sender2Blances := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender2).GetCoins()
	suite.Equal("30000000000000000000", moduleAccountBalances.AmountOf(unidenomBTC).String())
	suite.Equal(fmt.Sprintf("301%s,30000000000000000000%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal(fmt.Sprintf("2999999799%s,10000000000000000000%s,20000000000000000000%s", denomBTC, denomStandard, unidenomBTC), sender2Blances.String())

	// Test remove liquidity (remove part)
	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(
		sdk.NewInt(1),
		sdk.NewCoin(unidenomBTC, withdraw),
		sdk.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender1,
	)

	err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	suite.NoError(err)

	moduleAccountBalances = suite.app.SupplyKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr).GetCoins()
	sender1Blances = suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender1).GetCoins()
	suite.Equal("20000000000000000000", moduleAccountBalances.AmountOf(unidenomBTC).String())
	suite.Equal(fmt.Sprintf("3000000000%s,30000000000000000000%s", denomBTC, denomStandard), sender1Blances.String())
	suite.Equal(fmt.Sprintf("201%s,20000000000000000000%s", denomBTC, denomStandard), reservePoolBalances.String())

	// Test remove liquidity (remove all)
	withdraw, _ = sdk.NewIntFromString("20000000000000000000")
	msgRemove = types.NewMsgRemoveLiquidity(
		sdk.NewInt(1),
		sdk.NewCoin(unidenomBTC, withdraw),
		sdk.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender2,
	)

	err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	suite.NoError(err)

	moduleAccountBalances = suite.app.SupplyKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances = suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr).GetCoins()
	sender1Blances = suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender1).GetCoins()
	suite.Equal("0", moduleAccountBalances.AmountOf(unidenomBTC).String())
	suite.Equal(fmt.Sprintf("3000000000%s,30000000000000000000%s", denomBTC, denomStandard), sender1Blances.String())
	suite.Equal("", reservePoolBalances.String())
}
