package keeper_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	"github.com/tendermint/tendermint/crypto/tmhash"
// 	tmbytes "github.com/tendermint/tendermint/libs/bytes"
// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/irisnet/irismod/modules/htlc/keeper"
// 	"github.com/irisnet/irismod/modules/htlc/types"
// 	"github.com/irisnet/irismod/simapp"
// )

// var (
// 	id               tmbytes.HexBytes
// 	sender           sdk.AccAddress
// 	recipient        sdk.AccAddress
// 	deputyAddress    sdk.AccAddress
// 	idStr            string
// 	senderStr        string
// 	recipientStr     string
// 	deputyAddressStr string

// 	initCoinAmt          = sdk.NewInt(100)
// 	receiverOnOtherChain = "receiverOnOtherChain"
// 	senderOnOtherChain   = "senderOnOtherChain"
// 	amount               = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
// 	secret               = tmbytes.HexBytes(tmhash.Sum([]byte("secret")))
// 	secretStr            = secret.String()
// 	timestamp            = uint64(1580000000)
// 	hashLock             = tmbytes.HexBytes(tmhash.Sum(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
// 	hashLockStr          = hashLock.String()
// 	timeLock             = uint64(50)
// 	transfer             = true
// 	notTransfer          = false
// 	incoming             = types.Incoming
// 	outgoing             = types.Outgoing
// 	invalid              = types.Invalid
// )

// type KeeperTestSuite struct {
// 	suite.Suite

// 	cdc    codec.JSONMarshaler
// 	ctx    sdk.Context
// 	keeper *keeper.Keeper
// 	app    *simapp.SimApp
// }

// func TestKeeperTestSuite(t *testing.T) {
// 	suite.Run(t, new(KeeperTestSuite))
// }

// func (suite *KeeperTestSuite) SetupTest() {
// 	app := simapp.Setup(false)
// 	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})

// 	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
// 	suite.keeper = &app.HTLCKeeper
// 	suite.app = app

// 	suite.setTestAddrs()
// }

// func (suite *KeeperTestSuite) setTestAddrs() {
// 	testAddrs := simapp.AddTestAddrs(suite.app, suite.ctx, 10, initCoinAmt)
// 	sender = testAddrs[0]
// 	recipient = testAddrs[1]
// 	deputyAddress = testAddrs[2]
// 	id = tmbytes.HexBytes(tmhash.Sum(append(append(append(hashLock, sender...), recipient...), []byte(amount.String())...)))
// 	idStr = id.String()
// 	senderStr = sender.String()
// 	recipientStr = recipient.String()
// 	deputyAddressStr = deputyAddress.String()
// }

// func (suite *KeeperTestSuite) TestCreateHTLC() {
// 	_, err := suite.keeper.CreateHTLC(
// 		suite.ctx,
// 		sender,
// 		recipient,
// 		receiverOnOtherChain,
// 		senderOnOtherChain,
// 		amount,
// 		hashLock,
// 		timestamp,
// 		timeLock,
// 		notTransfer,
// 	)
// 	suite.NoError(err)

// 	htlc, found := suite.keeper.GetHTLC(suite.ctx, id)
// 	suite.True(found)

// 	suite.Equal(idStr, htlc.Id)
// 	suite.Equal(senderStr, htlc.Sender)
// 	suite.Equal(recipientStr, htlc.To)
// 	suite.Equal(receiverOnOtherChain, htlc.ReceiverOnOtherChain)
// 	suite.Equal(senderOnOtherChain, htlc.SenderOnOtherChain)
// 	suite.Equal(amount, htlc.Amount)
// 	suite.Equal("", htlc.Secret)
// 	suite.Equal(timestamp, htlc.Timestamp)
// 	suite.Equal(uint64(suite.ctx.BlockHeight())+timeLock, htlc.ExpirationHeight)
// 	suite.Equal(types.Open, htlc.State)
// 	suite.Equal(uint64(0), htlc.ClosedBlock)
// 	suite.Equal(notTransfer, htlc.Transfer)
// 	suite.Equal(invalid, htlc.Direction)

// 	htlcs := make([]types.HTLC, 0)
// 	suite.keeper.IterateHTLCExpiredQueueByHeight(
// 		suite.ctx,
// 		htlc.ExpirationHeight,
// 		func(hlock tmbytes.HexBytes, h types.HTLC) bool {
// 			htlcs = append(htlcs, h)
// 			return false
// 		},
// 	)

// 	suite.Len(htlcs, 1)
// 	suite.Equal(htlc, htlcs[0])

// 	senderCoin := suite.app.BankKeeper.GetBalance(suite.ctx, sender, sdk.DefaultBondDenom)
// 	suite.Equal(initCoinAmt.Sub(amount.AmountOf(sdk.DefaultBondDenom)), senderCoin.Amount)

// 	htlcModAccCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.keeper.GetHTLCAccount(suite.ctx).GetAddress(), sdk.DefaultBondDenom)
// 	suite.Equal(amount.AmountOf(sdk.DefaultBondDenom), htlcModAccCoin.Amount)
// }

// func (suite *KeeperTestSuite) TestCreateIncomingHTLT() {
// 	_, err := suite.keeper.CreateHTLC(
// 		suite.ctx,
// 		sender,
// 		recipient,
// 		receiverOnOtherChain,
// 		senderOnOtherChain,
// 		amount,
// 		hashLock,
// 		timestamp,
// 		timeLock,
// 		transfer,
// 	)
// 	suite.NoError(err)

// 	htlc, found := suite.keeper.GetHTLC(suite.ctx, id)
// 	suite.True(found)

// 	suite.Equal(idStr, htlc.Id)
// 	suite.Equal(senderStr, htlc.Sender)
// 	suite.Equal(recipientStr, htlc.To)
// 	suite.Equal(receiverOnOtherChain, htlc.ReceiverOnOtherChain)
// 	suite.Equal(senderOnOtherChain, htlc.SenderOnOtherChain)
// 	suite.Equal(amount, htlc.Amount)
// 	suite.Equal("", htlc.Secret)
// 	suite.Equal(timestamp, htlc.Timestamp)
// 	suite.Equal(uint64(suite.ctx.BlockHeight())+timeLock, htlc.ExpirationHeight)
// 	suite.Equal(types.Open, htlc.State)
// 	suite.Equal(uint64(0), htlc.ClosedBlock)
// 	suite.Equal(notTransfer, htlc.Transfer)
// 	suite.Equal(invalid, htlc.Direction)

// 	htlcs := make([]types.HTLC, 0)
// 	suite.keeper.IterateHTLCExpiredQueueByHeight(
// 		suite.ctx,
// 		htlc.ExpirationHeight,
// 		func(hlock tmbytes.HexBytes, h types.HTLC) bool {
// 			htlcs = append(htlcs, h)
// 			return false
// 		},
// 	)

// 	suite.Len(htlcs, 1)
// 	suite.Equal(htlc, htlcs[0])
// }

// func (suite *KeeperTestSuite) TestCreateOutgoingHTLT() {
// 	err := suite.keeper.CreateHTLC(suite.ctx, sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
// 	suite.NoError(err)

// 	// TODO
// }

// func (suite *KeeperTestSuite) TestClaimHTLC() {
// 	err := suite.keeper.CreateHTLC()
// 	suite.NoError(err)

// 	htlc, found := suite.keeper.GetHTLC(suite.ctx, hashLock)
// 	suite.True(found)

// 	suite.Equal(secretStr, htlc.Secret)
// 	suite.Equal(timestamp, htlc.Timestamp)
// 	suite.Equal(types.Completed, htlc.State)

// 	recipientCoin := suite.app.BankKeeper.GetBalance(suite.ctx, recipient, sdk.DefaultBondDenom)
// 	suite.Equal(initCoinAmt.Add(amount.AmountOf(sdk.DefaultBondDenom)), recipientCoin.Amount)

// 	htlcModAccCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.keeper.GetHTLCAccount(suite.ctx).GetAddress(), sdk.DefaultBondDenom)
// 	suite.Equal(sdk.ZeroInt(), htlcModAccCoin.Amount)
// }

// func (suite *KeeperTestSuite) TestClaimIncomingHTLT() {
// 	err := suite.keeper.CreateIncomingHTLT()
// 	suite.NoError(err)

// 	// TODO
// }

// func (suite *KeeperTestSuite) TestClaimOutgoingHTLT() {
// 	err := suite.keeper.CreateOutgoingHTLT()
// 	suite.NoError(err)

// 	// TODO
// }

// func (suite *KeeperTestSuite) TestRefundHTLC() {
// 	err := suite.keeper.CreateHTLC()
// 	suite.NoError(err)

// 	err := suite.keeper.RefundHTLC(suite.ctx, hashLock)
// 	suite.NoError(err)

// 	htlc, found := suite.keeper.GetHTLC(suite.ctx, hashLock)
// 	suite.True(found)

// 	suite.Equal(types.Refunded, htlc.State)

// 	senderCoin := suite.app.BankKeeper.GetBalance(suite.ctx, sender, sdk.DefaultBondDenom)
// 	suite.Equal(initCoinAmt, senderCoin.Amount)

// 	htlcModAccCoin := suite.app.BankKeeper.GetBalance(suite.ctx, suite.keeper.GetHTLCAccount(suite.ctx).GetAddress(), sdk.DefaultBondDenom)
// 	suite.Equal(sdk.ZeroInt(), htlcModAccCoin.Amount)
// }

// func (suite *KeeperTestSuite) TestRefundIncomingHTLT() {
// 	err := suite.keeper.CreateIncomingHTLT()
// 	suite.NoError(err)

// 	// TODO
// }

// func (suite *KeeperTestSuite) TestRefundOutgoingHTLT() {
// 	err := suite.keeper.CreateOutgoingHTLT()
// 	suite.NoError(err)

// 	// TODO
// }
