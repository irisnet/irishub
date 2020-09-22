package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/irisnet/irismod/modules/token/keeper"
	"github.com/irisnet/irismod/modules/token/types"
	"github.com/irisnet/irismod/simapp"
)

const (
	isCheckTx = false
)

var (
	denom    = types.GetNativeToken().Symbol
	owner    = sdk.AccAddress([]byte("tokenTest"))
	initAmt  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	bk          bankkeeper.Keeper
	app         *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper
	suite.app = app

	// set params
	suite.keeper.SetParamSet(suite.ctx, types.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestIssueToken() {
	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, owner)

	err := suite.keeper.IssueToken(suite.ctx, *msg)
	require.NoError(suite.T(), err)

	suite.True(suite.keeper.HasToken(suite.ctx, msg.Symbol))

	token, err := suite.keeper.GetToken(suite.ctx, msg.Symbol)
	require.NoError(suite.T(), err)

	suite.Equal(msg.MinUnit, token.GetMinUnit())
	suite.Equal(msg.Owner, token.GetOwner())

	ftJson, _ := json.Marshal(msg)
	tokenJson, _ := json.Marshal(token)
	suite.Equal(ftJson, tokenJson)
}

func (suite *KeeperTestSuite) TestEditToken() {

	suite.TestIssueToken()

	mintable := types.True
	msgEditToken := types.NewMsgEditToken("Bitcoin Token", "btc", 22000000, mintable, owner)
	err := suite.keeper.EditToken(suite.ctx, *msgEditToken)
	require.NoError(suite.T(), err)

	token2, err := suite.keeper.GetToken(suite.ctx, msgEditToken.Symbol)
	require.NoError(suite.T(), err)

	expToken := types.NewToken("btc", "Bitcoin Token", "satoshi", 18, 21000000, 22000000, mintable.ToBool(), owner)

	expJson, _ := json.Marshal(expToken)
	actJson, _ := json.Marshal(token2)
	suite.Equal(expJson, actJson)

}

func (suite *KeeperTestSuite) TestMintToken() {

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 1000, 2000, true, owner)

	err := suite.keeper.IssueToken(suite.ctx, *msg)
	require.NoError(suite.T(), err)

	suite.True(suite.keeper.HasToken(suite.ctx, msg.Symbol))

	amt := suite.bk.GetBalance(suite.ctx, owner, msg.MinUnit)
	suite.Equal("1000000000000000000000satoshi", amt.String())

	msgMintToken := types.NewMsgMintToken(msg.Symbol, owner, nil, 1000)
	err = suite.keeper.MintToken(suite.ctx, *msgMintToken)
	require.NoError(suite.T(), err)

	amt = suite.bk.GetBalance(suite.ctx, owner, msg.MinUnit)
	suite.Equal("2000000000000000000000satoshi", amt.String())
}

func (suite *KeeperTestSuite) TestTransferToken() {

	suite.TestIssueToken()

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferTokenOwner{
		SrcOwner: owner,
		DstOwner: dstOwner,
		Symbol:   "btc",
	}
	err := suite.keeper.TransferTokenOwner(suite.ctx, msg)
	require.NoError(suite.T(), err)

	token, err := suite.keeper.GetToken(suite.ctx, "btc")
	require.NoError(suite.T(), err)
	suite.Equal(dstOwner, token.GetOwner())
}
