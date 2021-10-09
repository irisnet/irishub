package token_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	tokenmodule "github.com/irisnet/irismod/modules/token"
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
	owner       = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	initAmt     = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin    = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite

	cdc    codec.JSONCodec
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

func (suite *HandlerSuite) issueToken(token types.Token) {
	err := suite.keeper.AddToken(suite.ctx, token)
	suite.NoError(err)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			token.MinUnit,
			sdk.NewIntWithDecimal(int64(token.InitialSupply), int(token.Scale)),
		),
	)

	err = suite.bk.MintCoins(suite.ctx, types.ModuleName, mintCoins)
	suite.NoError(err)

	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, mintCoins)
	suite.NoError(err)
}

func (suite *HandlerSuite) TestIssueToken() {
	h := tokenmodule.NewHandler(suite.keeper)

	nativeTokenAmt1 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, owner.String())

	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	nativeTokenAmt2 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	fee, err := suite.keeper.GetTokenIssueFee(suite.ctx, msg.Symbol)
	suite.NoError(err)

	suite.Equal(nativeTokenAmt1.Sub(fee.Amount), nativeTokenAmt2)

	mintTokenAmt := sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale))

	nativeTokenAmt3 := suite.bk.GetBalance(suite.ctx, owner, msg.MinUnit).Amount
	suite.Equal(nativeTokenAmt3, mintTokenAmt)
}

func (suite *HandlerSuite) TestMintToken() {
	token := types.NewToken("btc", "Bitcoin Network", "satoshi", 18, 1000, 2000, true, owner)
	suite.issueToken(token)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount
	suite.Equal(sdk.NewIntWithDecimal(int64(token.InitialSupply), int(token.Scale)), beginBtcAmt)

	beginNativeAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), denom).Amount

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintToken := types.NewMsgMintToken(token.Symbol, token.Owner, "", 1000)
	_, err := h(suite.ctx, msgMintToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount

	mintBtcAmt := sdk.NewIntWithDecimal(int64(msgMintToken.Amount), int(token.Scale))
	suite.Equal(beginBtcAmt.Add(mintBtcAmt), endBtcAmt)

	fee, err := suite.keeper.GetTokenMintFee(suite.ctx, token.Symbol)
	suite.NoError(err)

	endNativeAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), denom).Amount

	suite.Equal(beginNativeAmt.Sub(fee.Amount), endNativeAmt)
}

func (suite *HandlerSuite) TestBurnToken() {
	token := types.NewToken("btc", "Bitcoin Network", "satoshi", 18, 1000, 2000, true, owner)
	suite.issueToken(token)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount
	suite.Equal(sdk.NewIntWithDecimal(int64(token.InitialSupply), int(token.Scale)), beginBtcAmt)

	h := tokenmodule.NewHandler(suite.keeper)

	msgBurnToken := types.NewMsgBurnToken(token.Symbol, token.Owner, 200)
	_, err := h(suite.ctx, msgBurnToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount
	burnBtcAmt := sdk.NewIntWithDecimal(int64(msgBurnToken.Amount), int(token.Scale))

	suite.Equal(beginBtcAmt.Sub(burnBtcAmt), endBtcAmt)
}
