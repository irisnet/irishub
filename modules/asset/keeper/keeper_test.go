package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/modules/asset/keeper"
	"github.com/irisnet/irishub/modules/asset/types"
	"github.com/irisnet/irishub/simapp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper keeper.Keeper
	bk     bank.Keeper
	app    *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.keeper = app.AssetKeeper
	suite.bk = app.BankKeeper
	suite.app = app
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestKeeper_IssueToken() {
	addr := sdk.AccAddress([]byte("addr"))

	ft := types.NewFungibleToken(types.NATIVE, "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, addr)
	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, "btc"))

	token, found := suite.keeper.GetToken(suite.ctx, "btc")
	suite.True(found)

	suite.Equal(ft.GetDenom(), token.GetDenom())
	suite.Equal(ft.Owner, ft.Owner)

	ftJson, _ := json.Marshal(ft)
	tokenJson, _ := json.Marshal(token)
	suite.Equal(ftJson, tokenJson)
}

func (suite *KeeperTestSuite) TestKeeper_EditToken() {
	addr := sdk.AccAddress([]byte("addr"))

	ft := types.NewFungibleToken(types.NATIVE, "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(21000000, 0), true, addr)

	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, "i.btc"))

	token, found := suite.keeper.GetToken(suite.ctx, "i.btc")
	suite.True(found)

	suite.Equal(ft.GetDenom(), token.GetDenom())
	suite.Equal(ft.Owner, token.Owner)

	ftJson, _ := json.Marshal(ft)
	tokenJson, _ := json.Marshal(token)
	suite.Equal(ftJson, tokenJson)

	mintable := types.False
	msgEditToken := types.NewMsgEditToken("BTC Token", "btc", "btc", "btc", 0, mintable, addr)
	err = suite.keeper.EditToken(suite.ctx, msgEditToken)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestMintTokenKeeper(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))

	amtCoin, _ := sdk.NewIntFromString("1000000000000000000000000000")
	coin := sdk.Coins{sdk.NewCoin("iris-atto", amtCoin)}
	suite.bk.AddCoins(suite.ctx, addr, coin)

	ft := types.NewFungibleToken(types.NATIVE, "btc", "btc", 0, "", "satoshi", sdk.NewIntWithDecimal(1000, 0), sdk.NewIntWithDecimal(10000, 0), true, addr)
	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	assert.True(t, suite.keeper.HasToken(suite.ctx, "btc"))

	token, found := suite.keeper.GetToken(suite.ctx, "btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	suite.Equal(t, msgJson, assetJson)

	msgMintToken := types.NewMsgMintToken("btc", addr, nil, 1000)
	err = suite.keeper.MintToken(suite.ctx, msgMintToken)
	assert.NoError(t, err)

	balance := suite.bk.GetCoins(suite.ctx, addr)
	amt := balance.AmountOf("btc-min")
	assert.Equal(t, "2000", amt.String())
}

func (suite *KeeperTestSuite) TestTransferOwnerKeeper(t *testing.T) {
	srcOwner := sdk.AccAddress([]byte("TokenSrcOwner"))

	ft := types.NewFungibleToken(types.NATIVE, "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(21000000, 0), true, srcOwner)

	err := suite.keeper.IssueToken(suite.ctx, ft)
	assert.NoError(t, err)

	assert.True(t, suite.keeper.HasToken(suite.ctx, "i.btc"))

	token, found := suite.keeper.GetToken(suite.ctx, "i.btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, token.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		TokenId:  "btc",
	}
	err = suite.keeper.TransferTokenOwner(suite.ctx, msg)
	assert.NoError(t, err)

	token, found = suite.keeper.GetToken(suite.ctx, "i.btc")
	assert.True(t, found)
	assert.Equal(t, dstOwner, token.Owner)
}
