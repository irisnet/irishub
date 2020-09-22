package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/simapp"
)

var (
	initCoinAmt = sdk.NewInt(100)

	sender    sdk.AccAddress
	recipient sdk.AccAddress

	receiverOnOtherChain = "receiverOnOtherChain"
	amount               = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret               = tmbytes.HexBytes(tmhash.Sum([]byte("secret")))
	timestamp            = uint64(1580000000)
	hashLock             = tmbytes.HexBytes(tmhash.Sum(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	timeLock             = uint64(50)
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = &app.HTLCKeeper
	suite.app = app

	suite.setTestAddrs()
}

func (suite *KeeperTestSuite) setTestAddrs() {
	testAddrs := simapp.AddTestAddrs(suite.app, suite.ctx, 2, initCoinAmt)

	sender = testAddrs[0]
	recipient = testAddrs[1]
}

func (suite *KeeperTestSuite) setHTLC(state types.HTLCState) {
	expirationHeight := uint64(suite.ctx.BlockHeight()) + timeLock
	htlc := types.NewHTLC(sender, recipient, receiverOnOtherChain, amount, secret, timestamp, expirationHeight, state)

	suite.keeper.SetHTLC(suite.ctx, htlc, hashLock)

	if state == types.Open || state == types.Expired {
		suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, sender, types.ModuleName, amount)
	}
}

func (suite *KeeperTestSuite) TestCreateHTLC() {
	err := suite.keeper.CreateHTLC(suite.ctx, sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	suite.NoError(err)

	htlc, found := suite.keeper.GetHTLC(suite.ctx, hashLock)
	suite.True(found)

	suite.Equal(sender, htlc.Sender)
	suite.Equal(recipient, htlc.To)
	suite.Equal(receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	suite.Equal(amount, htlc.Amount)
	suite.Equal(tmbytes.HexBytes(nil), htlc.Secret)
	suite.Equal(timestamp, htlc.Timestamp)
	suite.Equal(uint64(suite.ctx.BlockHeight())+timeLock, htlc.ExpirationHeight)
	suite.Equal(types.Open, htlc.State)

	htlcs := make([]types.HTLC, 0)
	suite.keeper.IterateHTLCExpiredQueueByHeight(
		suite.ctx,
		htlc.ExpirationHeight,
		func(hlock tmbytes.HexBytes, h types.HTLC) bool {
			htlcs = append(htlcs, h)
			return false
		},
	)

	suite.Len(htlcs, 1)
	suite.Equal(htlc, htlcs[0])

	senderCoin := suite.app.BankKeeper.GetBalance(suite.ctx, sender, sdk.DefaultBondDenom)
	suite.Equal(initCoinAmt.Sub(amount.AmountOf(sdk.DefaultBondDenom)), senderCoin.Amount)

	htlcModAccCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.keeper.GetHTLCAccount(suite.ctx).GetAddress(), sdk.DefaultBondDenom)
	suite.Equal(amount.AmountOf(sdk.DefaultBondDenom), htlcModAccCoin.Amount)
}

func (suite *KeeperTestSuite) TestClaimHTLC() {
	suite.setHTLC(types.Open)

	err := suite.keeper.ClaimHTLC(suite.ctx, hashLock, secret)
	suite.NoError(err)

	htlc, found := suite.keeper.GetHTLC(suite.ctx, hashLock)
	suite.True(found)

	suite.Equal(secret, htlc.Secret)
	suite.Equal(timestamp, htlc.Timestamp)
	suite.Equal(types.Completed, htlc.State)

	recipientCoin := suite.app.BankKeeper.GetBalance(suite.ctx, recipient, sdk.DefaultBondDenom)
	suite.Equal(initCoinAmt.Add(amount.AmountOf(sdk.DefaultBondDenom)), recipientCoin.Amount)

	htlcModAccCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.keeper.GetHTLCAccount(suite.ctx).GetAddress(), sdk.DefaultBondDenom)
	suite.Equal(sdk.ZeroInt(), htlcModAccCoin.Amount)
}

func (suite *KeeperTestSuite) TestRefundHTLC() {
	suite.setHTLC(types.Expired)

	err := suite.keeper.RefundHTLC(suite.ctx, hashLock)
	suite.NoError(err)

	htlc, found := suite.keeper.GetHTLC(suite.ctx, hashLock)
	suite.True(found)

	suite.Equal(types.Refunded, htlc.State)

	senderCoin := suite.app.BankKeeper.GetBalance(suite.ctx, sender, sdk.DefaultBondDenom)
	suite.Equal(initCoinAmt, senderCoin.Amount)

	htlcModAccCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.keeper.GetHTLCAccount(suite.ctx).GetAddress(), sdk.DefaultBondDenom)
	suite.Equal(sdk.ZeroInt(), htlcModAccCoin.Amount)
}
