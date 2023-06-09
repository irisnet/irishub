package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/simapp"
)

const (
	denomBTC = "btc"
	denomETH = "eth"
)

var (
	denomStandard = sdk.DefaultBondDenom
	addrSender1   sdk.AccAddress
	addrSender2   sdk.AccAddress
)

// test that the params can be properly set and retrieved
type TestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *simapp.SimApp
	queryClient types.QueryClient
}

func (suite *TestSuite) SetupTest() {
	app := setupWithGenesisAccounts(suite.T())
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
	}
	for _, tc := range cases {
		suite.app.CoinswapKeeper.SetParams(suite.ctx, tc.params)

		feeParam := suite.app.CoinswapKeeper.GetParams(suite.ctx)
		suite.Equal(tc.params.Fee, feeParam.Fee)
	}
}

func setupWithGenesisAccounts(t *testing.T) *simapp.SimApp {
	amountInitStandard, _ := sdkmath.NewIntFromString("30000000000000000000")
	amountInitBTC, _ := sdkmath.NewIntFromString("3000000000")

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
	app := simapp.SetupWithGenesisAccounts(t, genAccs, acc1Balances, acc2Balances)
	return app
}

func (suite *TestSuite) TestLiquidity() {
	btcAmt, _ := sdkmath.NewIntFromString("100")
	standardAmt, _ := sdkmath.NewIntFromString("10000000000000000000")
	depositCoin := sdk.NewCoin(denomBTC, btcAmt)
	minReward := sdkmath.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)

	msg := types.NewMsgAddLiquidity(
		depositCoin,
		standardAmt,
		minReward,
		deadline.Unix(),
		addrSender1.String(),
	)
	_, err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	poolId := types.GetPoolId(denomBTC)
	pool, has := suite.app.CoinswapKeeper.GetPool(suite.ctx, poolId)
	suite.Require().True(has)

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	suite.Require().NoError(err)

	lptDenom := pool.LptDenom

	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal(
		"10000000000000000000",
		suite.app.BankKeeper.GetSupply(suite.ctx, lptDenom).Amount.String(),
	)

	expCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 100),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 2999999900),
		sdk.NewCoin(
			denomStandard,
			sdkmath.NewIntWithDecimal(2, 19).Sub(sdk.NewIntFromUint64(5000)),
		),
		sdk.NewCoin(lptDenom, sdkmath.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	// test add liquidity (pool exists)
	expLptDenom, _ := suite.app.CoinswapKeeper.GetLptDenomFromDenoms(
		suite.ctx,
		denomBTC,
		denomStandard,
	)
	suite.Require().Equal(expLptDenom, lptDenom)
	btcAmt, _ = sdkmath.NewIntFromString("201")
	standardAmt, _ = sdkmath.NewIntFromString("20000000000000000000")
	depositCoin = sdk.NewCoin(denomBTC, btcAmt)
	minReward = sdkmath.NewInt(1)
	deadline = time.Now().Add(1 * time.Minute)

	msg = types.NewMsgAddLiquidity(
		depositCoin,
		standardAmt,
		minReward,
		deadline.Unix(),
		addrSender2.String(),
	)
	_, err = suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.NoError(err)

	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender2Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender2)
	suite.Equal(
		"30000000000000000000",
		suite.app.BankKeeper.GetSupply(suite.ctx, lptDenom).Amount.String(),
	)

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 301),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(3, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 2999999799),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(1, 19)),
		sdk.NewCoin(lptDenom, sdkmath.NewIntWithDecimal(2, 19)),
	)
	suite.Equal(expCoins.Sort().String(), sender2Balances.Sort().String())

	// Test remove liquidity (remove part)
	withdraw, _ := sdkmath.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(
		sdkmath.NewInt(1),
		sdk.NewCoin(lptDenom, withdraw),
		sdkmath.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender1.String(),
	)

	_, err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	suite.NoError(err)

	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal(
		"20000000000000000000",
		suite.app.BankKeeper.GetSupply(suite.ctx, lptDenom).Amount.String(),
	)

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 3000000000),
		sdk.NewCoin(
			denomStandard,
			sdkmath.NewIntWithDecimal(3, 19).Sub(sdk.NewIntFromUint64(5000)),
		),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 201),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(2, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.String())

	// Test remove liquidity (remove all)
	withdraw, _ = sdkmath.NewIntFromString("20000000000000000000")
	msgRemove = types.NewMsgRemoveLiquidity(
		sdkmath.NewInt(1),
		sdk.NewCoin(lptDenom, withdraw),
		sdkmath.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender2.String(),
	)

	_, err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	suite.NoError(err)

	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal("0", suite.app.BankKeeper.GetSupply(suite.ctx, lptDenom).Amount.String())

	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 3000000000),
		sdk.NewCoin(
			denomStandard,
			sdkmath.NewIntWithDecimal(3, 19).Sub(sdkmath.NewIntFromUint64(5000)),
		),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())
	suite.Equal("", reservePoolBalances.String())
}

// TestLiquidity2 tests functionality of add liquidity unilaterally.
func (suite *TestSuite) TestLiquidity2() {
	// 1. initial liquidity (create pool)
	btcAmt, _ := sdkmath.NewIntFromString("100")                  // 10^2
	stdAmt, _ := sdkmath.NewIntFromString("10000000000000000000") // 10^19
	initMsg := types.NewMsgAddLiquidity(
		sdk.NewCoin(denomBTC, btcAmt),
		stdAmt,
		sdkmath.NewInt(1),
		time.Now().Add(1*time.Minute).Unix(),
		addrSender1.String(),
	)

	_, err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, initMsg)
	suite.NoError(err)

	pool, exist := suite.app.CoinswapKeeper.GetPool(suite.ctx, types.GetPoolId(denomBTC))
	suite.Require().True(exist)

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	suite.Require().NoError(err)

	// 1.1 lptAmt
	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal(
		"10000000000000000000",
		suite.app.BankKeeper.GetSupply(suite.ctx, pool.LptDenom).Amount.String(),
	)

	// 1.2 poolBalances
	expCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 100),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	// 1.3 accountBalances
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(
			denomBTC,
			2999999900,
		), // 3*10^9 - 100
		sdk.NewCoin(
			denomStandard,
			sdkmath.NewIntWithDecimal(2, 19).Sub(sdk.NewIntFromUint64(5000)),
		), // 2*10^19 - 5000
		sdk.NewCoin(
			pool.LptDenom,
			sdkmath.NewIntWithDecimal(1, 19),
		), // 10^19
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	// 2. add liquidity unilaterally

	btcAmt, _ = sdkmath.NewIntFromString("100")
	addMsg := types.NewMsgAddUnilateralLiquidity(
		denomBTC,
		sdk.NewCoin(denomBTC, btcAmt),
		sdkmath.NewInt(1),
		time.Now().Add(1*time.Minute).Unix(),
		addrSender2.String(),
	)

	_, err = suite.app.CoinswapKeeper.AddUnilateralLiquidity(suite.ctx, addMsg)
	suite.NoError(err)

	// 2.1 lptAmt
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender2Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender2)
	suite.Equal(
		"14135062787267695755",
		suite.app.BankKeeper.GetSupply(suite.ctx, pool.LptDenom).Amount.String(),
	) // todo theoretical lpt ammount

	// 2.2 poolBalances
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 200),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	// 2.3 accountBalances
	lptAmt, _ := sdkmath.NewIntFromString("4135062787267695755")
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 2999999900),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(3, 19)),
		sdk.NewCoin(pool.LptDenom, lptAmt),
	)
	suite.Equal(expCoins.Sort().String(), sender2Balances.Sort().String())
}

// TestLiquidity3 tests functionality of remove liquidity unilaterally.
func (suite *TestSuite) TestLiquidity3() {
	// 1. initial liquidity (create pool)
	btcAmt, _ := sdkmath.NewIntFromString("100")                  // 10^2
	stdAmt, _ := sdkmath.NewIntFromString("10000000000000000000") // 10^19
	initMsg := types.NewMsgAddLiquidity(
		sdk.NewCoin(denomBTC, btcAmt),
		stdAmt,
		sdkmath.NewInt(1),
		time.Now().Add(1*time.Minute).Unix(),
		addrSender1.String(),
	)

	_, err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, initMsg)
	suite.NoError(err)

	pool, exist := suite.app.CoinswapKeeper.GetPool(suite.ctx, types.GetPoolId(denomBTC))
	suite.Require().True(exist)

	poolAddr, err := sdk.AccAddressFromBech32(pool.EscrowAddress)
	suite.Require().NoError(err)

	// 1.1 lptAmt
	reservePoolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal(
		"10000000000000000000",
		suite.app.BankKeeper.GetSupply(suite.ctx, pool.LptDenom).Amount.String(),
	)

	// 1.2 poolBalances
	expCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 100),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	// 1.3 accountBalances
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(
			denomBTC,
			2999999900,
		), // 3*10^9 - 100
		sdk.NewCoin(
			denomStandard,
			sdkmath.NewIntWithDecimal(2, 19).Sub(sdk.NewIntFromUint64(5000)),
		), // 2*10^19 - 5000
		sdk.NewCoin(
			pool.LptDenom,
			sdkmath.NewIntWithDecimal(1, 19),
		), // 10^19
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())

	// 2. remove liquidity unilaterally

	btcAmt, _ = sdkmath.NewIntFromString("1")
	removeMsg := types.NewMsgRemoveUnilateralLiquidity(
		denomBTC,
		sdk.NewCoin(denomBTC, btcAmt),       // at least 1
		sdkmath.NewInt(5000000000000000000), // 5 * 10^18
		time.Now().Add(1*time.Minute).Unix(),
		addrSender1.String(),
	)

	_, err = suite.app.CoinswapKeeper.RemoveUnilateralLiquidity(suite.ctx, removeMsg)
	suite.NoError(err)

	// 2.1 lptAmt
	reservePoolBalances = suite.app.BankKeeper.GetAllBalances(suite.ctx, poolAddr)
	sender1Balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addrSender1)
	suite.Equal(
		"5000000000000000000",
		suite.app.BankKeeper.GetSupply(suite.ctx, pool.LptDenom).Amount.String(),
	)

	// 2.2 poolBalances
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 26),
		sdk.NewCoin(denomStandard, sdkmath.NewIntWithDecimal(1, 19)),
	)
	suite.Equal(expCoins.Sort().String(), reservePoolBalances.Sort().String())

	// 2.3 accountBalances
	lptAmt, _ := sdkmath.NewIntFromString("5000000000000000000")
	irisAmt, _ := sdkmath.NewIntFromString("19999999999999995000")
	expCoins = sdk.NewCoins(
		sdk.NewInt64Coin(denomBTC, 2999999974),
		sdk.NewCoin(denomStandard, irisAmt),
		sdk.NewCoin(pool.LptDenom, lptAmt),
	)
	suite.Equal(expCoins.Sort().String(), sender1Balances.Sort().String())
}
