package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/irisnet/irishub/modules/asset/internal/keeper"
	"github.com/irisnet/irishub/modules/asset/internal/types"
	"github.com/irisnet/irishub/simapp"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper keeper.Keeper
	sk     supply.Keeper
	bk     bank.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.keeper = app.AssetKeeper
	suite.bk = app.BankKeeper
	suite.sk = app.SupplyKeeper

	// set params
	suite.keeper.SetParamSet(suite.ctx, types.DefaultParams())
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

func (suite *KeeperTestSuite) TestMintTokenKeeper() {
	addr := sdk.AccAddress([]byte("addr"))

	// mint tokens to addr
	mintAmount, _ := sdk.NewIntFromString("1000000000000000000000000000")
	mintCoin := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, mintAmount)}
	err := suite.sk.MintCoins(suite.ctx, types.ModuleName, mintCoin)
	suite.NoError(err)
	err = suite.sk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr, mintCoin)
	suite.NoError(err)

	ft := types.NewFungibleToken(types.NATIVE, "btc", "btc", 0, "", "satoshi", sdk.NewIntWithDecimal(1000, 0), sdk.NewIntWithDecimal(10000, 0), true, addr)
	err = suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, "btc"))

	token, found := suite.keeper.GetToken(suite.ctx, "btc")
	suite.True(found)

	suite.Equal(ft.GetDenom(), token.GetDenom())
	suite.Equal(ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	suite.Equal(msgJson, assetJson)

	msgMintToken := types.NewMsgMintToken("btc", addr, nil, 1000)
	err = suite.keeper.MintToken(suite.ctx, msgMintToken)
	suite.NoError(err)

	balance := suite.bk.GetCoins(suite.ctx, addr)
	amt := balance.AmountOf("btc-min")
	suite.Equal("2000", amt.String())
}

func (suite *KeeperTestSuite) TestTransferOwnerKeeper() {
	srcOwner := sdk.AccAddress([]byte("TokenSrcOwner"))

	ft := types.NewFungibleToken(types.NATIVE, "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(21000000, 0), true, srcOwner)

	err := suite.keeper.IssueToken(suite.ctx, ft)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, "i.btc"))

	token, found := suite.keeper.GetToken(suite.ctx, "i.btc")
	suite.True(found)

	suite.Equal(ft.GetDenom(), token.GetDenom())
	suite.Equal(ft.Owner, token.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	suite.Equal(msgJson, assetJson)

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		TokenID:  "btc",
	}
	err = suite.keeper.TransferTokenOwner(suite.ctx, msg)
	suite.NoError(err)

	token, found = suite.keeper.GetToken(suite.ctx, "i.btc")
	suite.True(found)
	suite.Equal(dstOwner, token.Owner)
}
