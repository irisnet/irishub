package token_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	tokenmodule "github.com/irisnet/irismod/modules/token"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	"github.com/irisnet/irismod/modules/token/types"
	v1 "github.com/irisnet/irismod/modules/token/types/v1"
	"github.com/irisnet/irismod/simapp"
)

const (
	isCheckTx = false
)

var (
	nativeToken = v1.GetNativeToken()
	denom       = nativeToken.Symbol
	owner       = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	initAmt     = sdkmath.NewIntWithDecimal(100000000, int(6))
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
	app := simapp.Setup(suite.T(), isCheckTx)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper

	// set params
	err := suite.keeper.SetParams(suite.ctx, v1.DefaultParams())
	suite.NoError(err)

	// init tokens to addr
	err = suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func (suite *HandlerSuite) issueToken(token v1.Token) {
	err := suite.keeper.AddToken(suite.ctx, token, true)
	suite.NoError(err)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			token.MinUnit,
			sdkmath.NewIntWithDecimal(int64(token.InitialSupply), int(token.Scale)),
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

	msg := v1.NewMsgIssueToken(
		"btc",
		"satoshi",
		"Bitcoin Network",
		18,
		21000000,
		21000000,
		false,
		owner.String(),
	)

	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	nativeTokenAmt2 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount

	fee, err := suite.keeper.GetTokenIssueFee(suite.ctx, msg.Symbol)
	suite.NoError(err)

	suite.Equal(nativeTokenAmt1.Sub(fee.Amount), nativeTokenAmt2)

	mintTokenAmt := sdkmath.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale))

	nativeTokenAmt3 := suite.bk.GetBalance(suite.ctx, owner, msg.MinUnit).Amount
	suite.Equal(nativeTokenAmt3, mintTokenAmt)
}

func (suite *HandlerSuite) TestMintToken() {
	token := v1.NewToken("btc", "Bitcoin Network", "satoshi", 18, 1000, 2000, true, owner)
	suite.issueToken(token)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount
	suite.Equal(
		sdkmath.NewIntWithDecimal(int64(token.InitialSupply), int(token.Scale)),
		beginBtcAmt,
	)

	beginNativeAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), denom).Amount

	h := tokenmodule.NewHandler(suite.keeper)

	msgMintToken := &v1.MsgMintToken{
		Coin: sdk.Coin{
			Denom:  token.MinUnit,
			Amount: sdkmath.NewInt(1000),
		},
		Owner: token.Owner,
	}
	_, err := h(suite.ctx, msgMintToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount

	suite.Equal(beginBtcAmt.Add(msgMintToken.Coin.Amount), endBtcAmt)

	fee, err := suite.keeper.GetTokenMintFee(suite.ctx, token.Symbol)
	suite.NoError(err)

	endNativeAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), denom).Amount

	suite.Equal(beginNativeAmt.Sub(fee.Amount), endNativeAmt)
}

func (suite *HandlerSuite) TestBurnToken() {
	token := v1.NewToken("btc", "Bitcoin Network", "satoshi", 18, 1000, 2000, true, owner)
	suite.issueToken(token)

	beginBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount
	suite.Equal(
		sdkmath.NewIntWithDecimal(int64(token.InitialSupply), int(token.Scale)),
		beginBtcAmt,
	)

	h := tokenmodule.NewHandler(suite.keeper)

	msgBurnToken := &v1.MsgBurnToken{
		Coin: sdk.Coin{
			Denom:  token.MinUnit,
			Amount: sdkmath.NewInt(1000),
		},
		Sender: token.Owner,
	}
	_, err := h(suite.ctx, msgBurnToken)
	suite.NoError(err)

	endBtcAmt := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.MinUnit).Amount

	suite.Equal(beginBtcAmt.Sub(msgBurnToken.Coin.Amount), endBtcAmt)
}
