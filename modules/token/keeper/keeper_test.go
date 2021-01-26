package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
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
	owner    = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
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

func (suite *KeeperTestSuite) setToken(token types.Token) {
	err := suite.keeper.AddToken(suite.ctx, token)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) issueToken(token types.Token) {
	suite.setToken(token)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			token.MinUnit,
			sdk.NewIntWithDecimal(int64(token.InitialSupply), int(token.Scale)),
		),
	)

	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, mintCoins)
	suite.NoError(err)

	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, mintCoins)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestIssueToken() {
	token := types.NewToken("btc", "Bitcoin Network", "satoshi", 18, 21000000, 21000000, false, owner)

	err := suite.keeper.IssueToken(
		suite.ctx, token.Symbol, token.Name,
		token.MinUnit, token.Scale, token.InitialSupply,
		token.MaxSupply, token.Mintable, token.GetOwner(),
	)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, token.Symbol))

	issuedToken, err := suite.keeper.GetToken(suite.ctx, token.Symbol)
	suite.NoError(err)

	suite.Equal(token.MinUnit, issuedToken.GetMinUnit())
	suite.Equal(token.Owner, issuedToken.GetOwner().String())

	suite.EqualValues(&token, issuedToken.(*types.Token))
}

func (suite *KeeperTestSuite) TestEditToken() {
	token := types.NewToken("btc", "Bitcoin Network", "satoshi", 18, 21000000, 21000000, false, owner)
	suite.setToken(token)

	symbol := "btc"
	name := "Bitcoin Token"
	mintable := types.True
	maxSupply := uint64(22000000)

	err := suite.keeper.EditToken(suite.ctx, symbol, name, maxSupply, mintable, owner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetToken(suite.ctx, symbol)
	suite.NoError(err)

	expToken := types.NewToken("btc", "Bitcoin Token", "satoshi", 18, 21000000, 22000000, mintable.ToBool(), owner)

	suite.EqualValues(newToken.(*types.Token), &expToken)
}

func (suite *KeeperTestSuite) TestMintToken() {
	token := types.NewToken("btc", "Bitcoin Network", "satoshi", 18, 1000, 2000, true, owner)
	suite.issueToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit)
	suite.Equal("1000000000000000000000satoshi", amt.String())

	mintAmount := uint64(1000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintToken(suite.ctx, token.Symbol, mintAmount, recipient, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit)
	suite.Equal("2000000000000000000000satoshi", amt.String())

	// mint token without owner

	err = suite.keeper.MintToken(suite.ctx, token.Symbol, mintAmount, owner, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")

	token = types.NewToken("atom", "Cosmos Hub", "uatom", 6, 1000, 2000, true, sdk.AccAddress{})
	suite.issueToken(token)

	err = suite.keeper.MintToken(suite.ctx, token.Symbol, mintAmount, owner, sdk.AccAddress{})
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestBurnToken() {
	token := types.NewToken("btc", "Bitcoin Network", "satoshi", 18, 1000, 2000, true, owner)
	suite.issueToken(token)

	amt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit)
	suite.Equal("1000000000000000000000satoshi", amt.String())

	burnedAmount := uint64(200)

	err := suite.keeper.BurnToken(suite.ctx, token.Symbol, burnedAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit)
	suite.Equal("800000000000000000000satoshi", amt.String())
}

func (suite *KeeperTestSuite) TestTransferToken() {
	token := types.NewToken("btc", "Bitcoin Network", "satoshi", 18, 21000000, 21000000, false, owner)
	suite.setToken(token)

	dstOwner := sdk.AccAddress(tmhash.SumTruncated([]byte("TokenDstOwner")))

	err := suite.keeper.TransferTokenOwner(suite.ctx, token.Symbol, token.GetOwner(), dstOwner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetToken(suite.ctx, token.Symbol)
	suite.NoError(err)

	suite.Equal(dstOwner, newToken.GetOwner())
}
