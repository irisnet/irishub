package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/config"
	"github.com/irisnet/irishub/modules/coinswap/internal/keeper"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
	"github.com/irisnet/irishub/simapp"
)

const (
	denomIris      = config.Iris
	denomIrisAtto  = config.IrisAtto
	denomBTC       = "btc"
	denomBTCMin    = "btc-min"
	denomETH       = "eth"
	denomETHMin    = "eth-min"
	unidenomBTC    = types.FormatUniABSPrefix + "btc"
	unidenomBTCMin = types.FormatUniABSPrefix + "btc-min"
	unidenomETH    = types.FormatUniABSPrefix + "eth"
	unidenomETHMin = types.FormatUniABSPrefix + "eth-min"
)

// test that the params can be properly set and retrieved
type KeeperTestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.app = app

	_ = sdk.RegisterDenom(denomIris, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(denomIrisAtto, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(denomBTC, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(denomBTCMin, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(denomETH, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(denomETHMin, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(unidenomBTC, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(unidenomBTCMin, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(unidenomETH, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(unidenomETHMin, sdk.NewDecWithPrec(1, sdk.Precision))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestParams() {
	cases := []struct {
		params types.Params
	}{
		{types.DefaultParams()},
		{types.NewParams(sdk.NewDecWithPrec(5, 10))},
	}
	for _, tc := range cases {
		suite.app.CoinswapKeeper.SetParams(suite.ctx, tc.params)

		feeParam := suite.app.CoinswapKeeper.GetParams(suite.ctx)
		require.Equal(suite.T(), tc.params.Fee, feeParam.Fee)
	}
}

func (suite *KeeperTestSuite) TestKeeper_UpdateLiquidity() {
	amountInit, _ := sdk.NewIntFromString("10000000000000000000")

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

	uniId, _ := types.GetUniId(denomBTCMin, denomIrisAtto)
	poolAddr := keeper.GetReservePoolAddr(uniId)

	btcAmt, _ := sdk.NewIntFromString("1")
	depositCoin := sdk.NewCoin("btc-min", btcAmt)

	irisAmt, _ := sdk.NewIntFromString("10000000000000000000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, irisAmt, minReward, deadline.Unix(), addrSender)

	err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	require.Nil(suite.T(), err)

	reservePoolBalances := suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr).GetCoins()
	moduleAccountBalances := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName).GetCoins()
	require.Equal(suite.T(), "1btc-min,10000000000000000000iris-atto", reservePoolBalances.String())
	require.Equal(suite.T(), "10000000000000000000", moduleAccountBalances.AmountOf("uni:btc-min").String())

	senderBlances := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()
	require.Equal(suite.T(), "9999999999999999999btc-min,10000000000000000000uni:btc-min", senderBlances.String())

	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(
		sdk.NewInt(1),
		sdk.NewCoin("uni:btc-min", withdraw),
		sdk.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender,
	)

	err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	require.Nil(suite.T(), err)

	poolAccout := suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr)
	acc := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender)
	require.Equal(suite.T(), "", poolAccout.GetCoins().String())
	require.Equal(suite.T(), "10000000000000000000btc-min,10000000000000000000iris-atto", acc.GetCoins().String())
}
