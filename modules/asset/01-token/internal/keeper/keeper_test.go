package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/keeper"
	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
	"github.com/irisnet/irishub/simapp"
)

const (
	isCheck = false
)

var (
	denom    = types.IrisToken().MinUnit
	owner    = sdk.AccAddress([]byte("tokenTest"))
	initAmt  = sdk.NewIntWithDecimal(types.IrisToken().InitialSupply.Int64(), int(types.IrisToken().Scale))
	initCoin = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

type KeeperSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper keeper.Keeper
	sk     supply.Keeper
	bk     bank.Keeper
}

func (suite *KeeperSuite) SetupTest() {

	app := simapp.Setup(isCheck)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheck, abci.Header{})
	suite.keeper = app.AssetKeeper.TokenKeeper
	suite.bk = app.BankKeeper
	suite.sk = app.SupplyKeeper

	// set params
	suite.keeper.SetParamSet(suite.ctx, types.DefaultParams())

	// init tokens to addr
	err := suite.sk.MintCoins(suite.ctx, types.SubModuleName, initCoin)
	suite.NoError(err)
	err = suite.sk.SendCoinsFromModuleToAccount(suite.ctx, types.SubModuleName, owner, initCoin)
	suite.NoError(err)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperSuite))
}

func (suite *KeeperSuite) TestIssueToken() {
	ft := types.NewFungibleToken("btc", "Bitcoin network", 18, "satoshi", sdk.NewIntWithDecimal(21, 6), sdk.NewIntWithDecimal(21, 6), false, owner)
	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, ft))

	token, found := suite.keeper.GetToken(suite.ctx, ft.Symbol)
	suite.True(found)

	suite.Equal(ft.GetMinUnit(), token.GetMinUnit())
	suite.Equal(ft.Owner, ft.Owner)

	ftJson, _ := json.Marshal(ft)
	tokenJson, _ := json.Marshal(token)
	suite.Equal(ftJson, tokenJson)
}

func (suite *KeeperSuite) TestEditToken() {

	ft := types.NewFungibleToken("btc", "Bitcoin network", 18, "satoshi", sdk.NewIntWithDecimal(21, 6), sdk.NewIntWithDecimal(21, 6), false, owner)

	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, ft))

	token, found := suite.keeper.GetToken(suite.ctx, ft.Symbol)
	suite.True(found)

	suite.Equal(ft.GetMinUnit(), token.GetMinUnit())
	suite.Equal(ft.Owner, token.Owner)

	ftJson, _ := json.Marshal(ft)
	tokenJson, _ := json.Marshal(token)
	suite.Equal(ftJson, tokenJson)

	mintable := types.False
	msgEditToken := types.NewMsgEditToken("BTC Token", "btc", 22000000, mintable, owner)
	err = suite.keeper.EditToken(suite.ctx, msgEditToken)
	suite.NoError(err)

	token2, found := suite.keeper.GetToken(suite.ctx, msgEditToken.Symbol)
	suite.True(found)

	expToken := types.NewFungibleToken("btc", "BTC Token", 18, "satoshi", sdk.NewIntWithDecimal(21, 6), sdk.NewIntWithDecimal(22, 6), false, owner)

	expJson, _ := json.Marshal(expToken)
	actJson, _ := json.Marshal(token2)
	suite.Equal(expJson, actJson)

}

func (suite *KeeperSuite) TestMintToken() {

	ft := types.NewFungibleToken("btc", "Bitcoin network", 0, "satoshi", sdk.NewIntWithDecimal(1000, 0), sdk.NewIntWithDecimal(3000, 0), true, owner)
	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, ft))

	token, found := suite.keeper.GetToken(suite.ctx, "btc")
	suite.True(found)

	suite.Equal(ft.GetMinUnit(), token.GetMinUnit())
	suite.Equal(ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	suite.Equal(msgJson, assetJson)

	balance := suite.bk.GetCoins(suite.ctx, owner)
	amt := balance.AmountOf(ft.MinUnit)
	suite.Equal("1000", amt.String())

	msgMintToken := types.NewMsgMintToken(ft.Symbol, owner, nil, 1000)
	_, err = suite.keeper.MintToken(suite.ctx, msgMintToken)
	suite.NoError(err)

	balance = suite.bk.GetCoins(suite.ctx, owner)
	amt = balance.AmountOf(ft.MinUnit)
	suite.Equal("2000", amt.String())
}

func (suite *KeeperSuite) TestTransferToken() {

	ft := types.NewFungibleToken("btc", "Bitcoin network", 18, "satoshi", sdk.NewIntWithDecimal(21, 6), sdk.NewIntWithDecimal(21, 6), false, owner)

	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, ft))

	token, found := suite.keeper.GetToken(suite.ctx, ft.Symbol)
	suite.True(found)

	suite.Equal(ft.GetMinUnit(), token.GetMinUnit())
	suite.Equal(ft.Owner, token.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	suite.Equal(msgJson, assetJson)

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferToken{
		SrcOwner: ft.Owner,
		DstOwner: dstOwner,
		Symbol:   ft.Symbol,
	}
	err = suite.keeper.TransferToken(suite.ctx, msg)
	suite.NoError(err)

	token, found = suite.keeper.GetToken(suite.ctx, ft.Symbol)
	suite.True(found)
	suite.Equal(dstOwner, token.Owner)
}

func (suite *KeeperSuite) TestBurnToken() {
	burnCoin := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000)))
	err := suite.keeper.BurnToken(suite.ctx, types.MsgBurnToken{
		Sender: owner,
		Amount: burnCoin,
	})
	suite.NoError(err)

	balance := suite.bk.GetCoins(suite.ctx, owner)

	suite.Equal(balance, initCoin.Sub(burnCoin))
}
