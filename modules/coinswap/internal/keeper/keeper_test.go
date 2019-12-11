package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/modules/coinswap/internal/keeper"
	"github.com/irisnet/irishub/modules/coinswap/internal/types"
	"github.com/irisnet/irishub/simapp"
)

const (
	denomStandard = types.StandardDenom
	denomBTC      = "btc"
	denomETH      = "eth"
	unidenomBTC   = types.FormatUniABSPrefix + "btc"
	unidenomETH   = types.FormatUniABSPrefix + "eth"
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
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestParams() {
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

func (suite *KeeperTestSuite) TestKeeper_UpdateLiquidity() {
	amountInit, _ := sdk.NewIntFromString("10000000000000000000")

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

	uniDenom, _ := types.GetUniDenomFromDenoms(denomBTC, denomStandard)
	suite.Equal(uniDenom, unidenomBTC)
	poolAddr := keeper.GetReservePoolAddr(uniDenom)

	btcAmt, _ := sdk.NewIntFromString("1")
	depositCoin := sdk.NewCoin(denomBTC, btcAmt)

	standardAmt, _ := sdk.NewIntFromString("10000000000000000000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, standardAmt, minReward, deadline.Unix(), addrSender)

	err := suite.app.CoinswapKeeper.AddLiquidity(suite.ctx, msg)
	suite.Nil(err)

	reservePoolBalances := suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr).GetCoins()
	moduleAccountBalances := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName).GetCoins()
	suite.Equal(fmt.Sprintf("1%s,10000000000000000000%s", denomBTC, denomStandard), reservePoolBalances.String())
	suite.Equal("10000000000000000000", moduleAccountBalances.AmountOf(unidenomBTC).String())

	senderBlances := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()
	suite.Equal(fmt.Sprintf("9999999999999999999%s,10000000000000000000%s", denomBTC, unidenomBTC), senderBlances.String())

	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(
		sdk.NewInt(1),
		sdk.NewCoin(unidenomBTC, withdraw),
		sdk.NewInt(1),
		suite.ctx.BlockHeader().Time.Unix(),
		addrSender,
	)

	err = suite.app.CoinswapKeeper.RemoveLiquidity(suite.ctx, msgRemove)
	suite.Nil(err)

	poolAccout := suite.app.AccountKeeper.GetAccount(suite.ctx, poolAddr)
	acc := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender)
	suite.Equal("", poolAccout.GetCoins().String())
	suite.Equal(fmt.Sprintf("10000000000000000000%s,10000000000000000000%s", denomBTC, denomStandard), acc.GetCoins().String())
}
