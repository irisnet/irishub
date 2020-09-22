package token_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/irisnet/irismod/modules/token"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	"github.com/irisnet/irismod/modules/token/types"
	"github.com/irisnet/irismod/simapp"
)

const (
	isCheckTx = false
)

var (
	nativeToken = types.GetNativeToken()
	denom       = nativeToken.Symbol
	owner       = sdk.AccAddress([]byte("tokenTest"))
	initAmt     = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin    = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper tokenkeeper.Keeper
	bk     bankkeeper.Keeper
}

func (suite *HandlerSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper

	// set params
	suite.keeper.SetParamSet(suite.ctx, types.DefaultParams())
	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func (suite *HandlerSuite) TestIssueToken() {
	h := token.NewHandler(suite.keeper)

	nativeTokenAmt1 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, owner)

	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	nativeTokenAmt2 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	fee := suite.keeper.GetTokenIssueFee(suite.ctx, msg.Symbol)

	suite.Equal(nativeTokenAmt1.Sub(fee.Amount), nativeTokenAmt2)

	mintTokenAmt := sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale))

	nativeTokenAmt3 := suite.bk.GetBalance(suite.ctx, owner, msg.MinUnit).Amount
	suite.Equal(nativeTokenAmt3, mintTokenAmt)
}

func (suite *HandlerSuite) TestMintToken() {
	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 1000, 2000, true, owner)

	err := suite.keeper.IssueToken(suite.ctx, *msg)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, msg.Symbol))

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, owner, msg.MinUnit).Amount
	suite.Equal(sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale)), beginBtcAmt)

	beginNativeAmt := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	h := token.NewHandler(suite.keeper)

	msgMintToken := types.NewMsgMintToken(msg.Symbol, owner, nil, 1000)
	_, err = h(suite.ctx, msgMintToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, owner, msg.MinUnit).Amount

	mintBtcAmt := sdk.NewIntWithDecimal(int64(msgMintToken.Amount), int(msg.Scale))
	suite.Equal(beginBtcAmt.Add(mintBtcAmt), endBtcAmt)

	fee := suite.keeper.GetTokenMintFee(suite.ctx, msg.Symbol)
	endNativeAmt := suite.bk.GetBalance(suite.ctx, owner, denom).Amount
	suite.Equal(beginNativeAmt.Sub(fee.Amount), endNativeAmt)
}
