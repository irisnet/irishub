package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/simapp"
)

const (
	denomStandard = sdk.DefaultBondDenom
	denomBTC      = "btc"
	denomETH      = "eth"
)

var (
	addrSender1 sdk.AccAddress
	addrSender2 sdk.AccAddress
)

// test that the params can be properly set and retrieved
type TestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *simapp.SimApp
	queryClient types.QueryClient
}

func (suite *TestSuite) SetupTest() {
	app := setupWithGenesisAccounts()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.CoinswapKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient

	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/\-]{2,127}`
	})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestParams() {
	cases := []struct {
		params types.Params
	}{
		{types.DefaultParams()},
		{types.NewParams(sdk.NewDecWithPrec(5, 10))},
	}
	for _, tc := range cases {
		suite.app.CoinswapKeeper.SetParams(suite.ctx, tc.params)

		feeParam := suite.app.CoinswapKeeper.GetParams(suite.ctx)
		suite.Equal(tc.params.Fee, feeParam.Fee)
	}
}

func setupWithGenesisAccounts() *simapp.SimApp {
	amountInitStandard, _ := sdk.NewIntFromString("30000000000000000000")
	amountInitBTC, _ := sdk.NewIntFromString("3000000000")

	addrSender1 = sdk.AccAddress(tmhash.SumTruncated([]byte("addrSender1")))
	addrSender2 = sdk.AccAddress(tmhash.SumTruncated([]byte("addrSender2")))
	acc1Balances := banktypes.Balance{
		Address: addrSender1.String(),
		Coins: sdk.NewCoins(
			sdk.NewCoin(denomStandard, amountInitStandard),
			sdk.NewCoin(denomBTC, amountInitBTC),
		),
	}

	acc2Balances := banktypes.Balance{
		Address: addrSender2.String(),
		Coins: sdk.NewCoins(
			sdk.NewCoin(denomStandard, amountInitStandard),
			sdk.NewCoin(denomBTC, amountInitBTC),
		),
	}

	acc1 := &authtypes.BaseAccount{
		Address: addrSender1.String(),
	}
	acc2 := &authtypes.BaseAccount{
		Address: addrSender2.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc1, acc2}
	app := simapp.SetupWithGenesisAccounts(genAccs, acc1Balances, acc2Balances)
	return app
}

func (suite *TestSuite) TestLiquidity() {
	btcAmt, _ := sdk.NewIntFromString("100")
	standardAmt, _ := sdk.NewIntFromString("10000000000000000000")
	depositCoin := sdk.NewCoin(denomBTC, btcAmt)
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)

	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender1.String())
	_, err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	poolId := types.GetPoolId(denomBTC)
	pool, has := suite.app.CoinswapKeeper.GetPool(suite.ctx, poolId)
	suite.Require().True(has)

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	suite.Require().NoError(err)

	lptDenom := pool.LptDenom

	supply := suite.app.BankKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal("10000000000000000000", supply.AmountOf(lptDenom).String())

	expCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 100),
		sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 2999999900),
		sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(2, 19)),
		sdk.NewCoin(lptDenom, sdk.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	// test add liquidity (pool exists)
	expLptDenom, _ := suite.app.CoinswapKeeper.GetLptDenomFromDenoms(suite.ctx, denomBTC, denomStandard)
	suite.Require().Equal(expLptDenom, lptDenom)
	btcAmt, _ = sdk.NewIntFromString("201")
	standardAmt, _ = sdk.NewIntFromString("20000000000000000000")
	depositCoin = sdk.NewCoin(denomBTC, btcAmt)
	minReward = sdk.NewInt(1)
	deadline = time.Now().Add(1 * time.Minute)

	msg = types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender2.String())
	_, err = suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	supply = suite.app.BankKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender2Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender2)
	suite.Equal("30000000000000000000", supply.AmountOf(lptDenom).String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 301),
		sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(3, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 2999999799),
		sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(1, 19)),
		sdk.NewCoin(lptDenom, sdk.NewIntWithDecimal(2, 19)),
	)
	suite.Equal(expCoins.Sort().String(), sender2Balances.Sort().String())

	// Test remove liquidity (remove part)
	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(
		sdk.NewInt(1),
		sdk.NewCoin(lptDenom, withdraw),
		sdk.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender1.String(),
	)

	_, err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	suite.NoError(err)

	supply = suite.app.BankKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal("20000000000000000000", supply.AmountOf(lptDenom).String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 3000000000),
		sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(3, 19)),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 201),
		sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(2, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.String())

	// Test remove liquidity (remove all)
	withdraw, _ = sdk.NewIntFromString("20000000000000000000")
	msgRemove = types.NewMsgRemoveLiquidity(
		sdk.NewInt(1),
		sdk.NewCoin(lptDenom, withdraw),
		sdk.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender2.String(),
	)

	_, err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	suite.NoError(err)

	supply = suite.app.BankKeeper.GetSupply(suite.ctx).GetTotal()
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal("0", supply.AmountOf(lptDenom).String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 3000000000),
		sdk.NewCoin(denomStandard, sdk.NewIntWithDecimal(3, 19)),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())
	suite.Equal("", reservePoolBalances.String())
}
