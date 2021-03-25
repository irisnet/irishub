package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/htlc"
	"github.com/irisnet/irismod/modules/htlc/keeper"
	"github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/simapp"
)

const (
	BNB_DENOM             = "htltbnb"
	OTHER_DENOM           = "htltinc"
	STARING_BNB_BALANCE   = int64(3000000000000)
	STARING_OTHER_BALANCE = int64(3000000000000)
)

type HTLCTestSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp

	deputy     sdk.AccAddress
	addrs      []sdk.AccAddress
	timestamps []uint64
	hashLocks  []tmbytes.HexBytes
	secrets    []tmbytes.HexBytes
}

func TestHTLCTestSuite(t *testing.T) {
	suite.Run(t, new(HTLCTestSuite))
}

func (suite *HTLCTestSuite) SetupTest() {
	app := simapp.SetupWithGenesisHTLC(NewHTLTGenesis(TestDeputy))
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: tmtime.Now()})

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.keeper = &app.HTLCKeeper
	suite.app = app

	_, addrs := GeneratePrivKeyAddressPairs(20)
	suite.addrs = addrs
	suite.deputy = TestDeputy

	coins := cs(c(BNB_DENOM, STARING_BNB_BALANCE), c(OTHER_DENOM, STARING_OTHER_BALANCE))
	for _, acc := range addrs {
		_ = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
		_ = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, acc, coins)
	}

	suite.setTestParams()
	suite.GenerateHTLCDetails()
}

func (suite *HTLCTestSuite) setTestParams() {
	params := suite.keeper.GetParams(suite.ctx)
	params.AssetParams[1].Active = true
	suite.keeper.SetParams(suite.ctx, params)
}

func (suite *HTLCTestSuite) GenerateHTLCDetails() {
	var timestamps []uint64
	var randomSecrets []tmbytes.HexBytes
	var randomHashLocks []tmbytes.HexBytes
	for i := 0; i < 15; i++ {
		timestamp := ts(i)
		randomSecret, _ := GenerateRandomSecret()
		randomHashLock := types.GetHashLock(randomSecret, timestamp)

		timestamps = append(timestamps, timestamp)
		randomSecrets = append(randomSecrets, randomSecret)
		randomHashLocks = append(randomHashLocks, randomHashLock)
	}
	suite.timestamps = timestamps
	suite.secrets = randomSecrets
	suite.hashLocks = randomHashLocks
}

func (suite *HTLCTestSuite) TestCreateHTLC() {
	currentTmTime := tmtime.Now()
	type args struct {
		sender             sdk.AccAddress
		receiver           sdk.AccAddress
		receiverOtherChain string
		senderOtherChain   string
		amount             sdk.Coins
		hashLock           tmbytes.HexBytes
		timestamp          uint64
		timeLock           uint64
		transfer           bool
		direction          types.SwapDirection
	}
	testCases := []struct {
		name          string
		blockTime     time.Time
		args          args
		expectPass    bool
		shouldBeFound bool
	}{
		{
			"incoming htlt",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[0],
				timestamp:          suite.timestamps[0],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[1],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 50000)),
				transfer:           true,
				direction:          types.Incoming,
			},
			true,
			true,
		},
		{
			"incoming htlt rate limited",
			currentTmTime.Add(time.Minute * 10),
			args{
				hashLock:           suite.hashLocks[12],
				timestamp:          suite.timestamps[12],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[1],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c("htltinc", 50000000000)),
				transfer:           true,
				direction:          types.Incoming,
			},
			true,
			true,
		},
		{
			"incoming htlt over rate limit",
			currentTmTime.Add(time.Minute * 10),
			args{
				hashLock:           suite.hashLocks[13],
				timestamp:          suite.timestamps[13],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[1],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c("htltinc", 50000000001)),
				transfer:           true,
				direction:          types.Incoming,
			},
			false,
			false,
		},
		{
			"outgoing htlt",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[0],
				timestamp:          suite.timestamps[0],
				timeLock:           MinTimeLock,
				sender:             suite.addrs[1],
				receiver:           suite.deputy,
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 50000)),
				transfer:           true,
				direction:          types.Outgoing,
			},
			true,
			true,
		},
		{

			"outgoing htlt amount not greater than fixed fee",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[1],
				timestamp:          suite.timestamps[1],
				timeLock:           MinTimeLock,
				sender:             suite.addrs[1],
				receiver:           suite.addrs[2],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 1000)),
				transfer:           true,
				direction:          types.Outgoing,
			},
			false,
			false,
		},
		{
			"unsupported asset",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[2],
				timestamp:          suite.timestamps[2],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[2],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c("htltxyz", 50000)),
				transfer:           true,
				direction:          types.Incoming,
			},
			false,
			false,
		},
		{
			"outside timestamp range",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[3],
				timestamp:          suite.timestamps[3] - 2000,
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[3],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 50000)),
				transfer:           true,
				direction:          types.Incoming,
			},
			false,
			false,
		},
		{
			"future timestamp",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[4],
				timestamp:          suite.timestamps[4] + 5000,
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[4],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 50000)),
				transfer:           true,
				direction:          types.Incoming,
			},
			false,
			false,
		},
		{
			"small height span on outgoing htlt",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[5],
				timestamp:          suite.timestamps[5],
				timeLock:           uint64(100),
				sender:             suite.addrs[5],
				receiver:           suite.deputy,
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 50000)),
				transfer:           true,
				direction:          types.Outgoing,
			},
			false,
			false,
		},
		{
			"big height span on outgoing htlt",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[6],
				timestamp:          suite.timestamps[6],
				timeLock:           uint64(300),
				sender:             suite.addrs[6],
				receiver:           suite.deputy,
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 50000)),
				transfer:           true,
				direction:          types.Outgoing,
			},
			false,
			false,
		},
		{
			"zero amount",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[7],
				timestamp:          suite.timestamps[7],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[7],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 0)),
				transfer:           true,
				direction:          types.Incoming,
			},
			false,
			false,
		},
		{
			"duplicate htlt",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[0],
				timestamp:          suite.timestamps[0],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[1],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 50000)),
				transfer:           true,
				direction:          types.Incoming,
			},
			false,
			true,
		},
		{
			"exactly at maximum amount",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[10],
				timestamp:          suite.timestamps[10],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[4],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 1000000000000)), // 10,000 BNB
				transfer:           true,
				direction:          types.Incoming,
			},
			true,
			true,
		},
		{
			"above maximum amount",
			currentTmTime,
			args{
				hashLock:           suite.hashLocks[11],
				timestamp:          suite.timestamps[11],
				timeLock:           MinTimeLock,
				sender:             suite.deputy,
				receiver:           suite.addrs[5],
				senderOtherChain:   SenderOnOtherChain,
				receiverOtherChain: ReceiverOnOtherChain,
				amount:             cs(c(BNB_DENOM, 1000000000001)), // 10,001 BNB
				transfer:           true,
				direction:          types.Incoming,
			},
			false,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(
			tc.name,
			func() {
				// Increment current asset supply to support outgoing htlcs
				suite.ctx = suite.ctx.WithBlockTime(tc.blockTime)
				if tc.args.direction == types.Outgoing {
					err := suite.keeper.IncrementCurrentAssetSupply(suite.ctx, tc.args.amount[0])
					suite.Nil(err, tc.name)
				}

				// Load asset denom (required for zero coins test case)
				var htlcAssetDenom string
				if len(tc.args.amount) == 1 {
					htlcAssetDenom = tc.args.amount[0].Denom
				} else {
					htlcAssetDenom = BNB_DENOM
				}

				senderBalancePre := suite.app.BankKeeper.GetBalance(suite.ctx, tc.args.sender, htlcAssetDenom)
				assetSupplyPre, _ := suite.keeper.GetAssetSupply(suite.ctx, htlcAssetDenom)

				// Create htlt
				id, err := suite.keeper.CreateHTLC(
					suite.ctx,
					tc.args.sender,
					tc.args.receiver,
					tc.args.receiverOtherChain,
					tc.args.senderOtherChain,
					tc.args.amount,
					tc.args.hashLock,
					tc.args.timestamp,
					tc.args.timeLock,
					tc.args.transfer,
				)

				// Load sender's account after htlt creation
				senderBalancePost := suite.app.BankKeeper.GetBalance(suite.ctx, tc.args.sender, htlcAssetDenom)
				assetSupplyPost, _ := suite.keeper.GetAssetSupply(suite.ctx, htlcAssetDenom)

				// Load expected htlt ID
				expectedHTLCID := types.GetID(tc.args.sender, tc.args.receiver, tc.args.amount, tc.args.hashLock)
				suite.Equal(expectedHTLCID, id, tc.name)

				if tc.expectPass {
					suite.NoError(err, tc.name)
					// Check incoming/outgoing asset supply increased
					switch tc.args.direction {
					case types.Incoming:
						suite.Equal(assetSupplyPre.IncomingSupply.Add(tc.args.amount[0]).String(), assetSupplyPost.IncomingSupply.String(), tc.name)
					case types.Outgoing:
						// Check coins moved
						suite.Equal(senderBalancePre.Sub(tc.args.amount[0]).String(), senderBalancePost.String(), tc.name)
						suite.Equal(assetSupplyPre.OutgoingSupply.Add(tc.args.amount[0]).String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					default:
						suite.Fail("should not have invalid direction", tc.name)
					}

					// Check htlt in store
					actualhtlc, found := suite.keeper.GetHTLC(suite.ctx, expectedHTLCID)
					suite.True(found, tc.name)
					suite.NotNil(actualhtlc, tc.name)

					// Confirm htlt contents
					expectedhtlc := types.HTLC{
						Id:                   expectedHTLCID.String(),
						Sender:               tc.args.sender.String(),
						To:                   tc.args.receiver.String(),
						ReceiverOnOtherChain: tc.args.receiverOtherChain,
						SenderOnOtherChain:   tc.args.senderOtherChain,
						Amount:               tc.args.amount,
						HashLock:             tc.args.hashLock.String(),
						Secret:               "",
						Timestamp:            tc.args.timestamp,
						ExpirationHeight:     uint64(suite.ctx.BlockHeight()) + tc.args.timeLock,
						State:                types.Open,
						ClosedBlock:          0,
						Transfer:             tc.args.transfer,
						Direction:            tc.args.direction,
					}
					suite.Equal(expectedhtlc, actualhtlc, tc.name)
				} else {
					suite.Error(err, tc.name)
					// Check coins not moved
					suite.Equal(senderBalancePre, senderBalancePost, tc.name)

					// Check incoming/outgoing asset supply not increased
					switch tc.args.direction {
					case types.Incoming:
						suite.Equal(assetSupplyPre.IncomingSupply, assetSupplyPost.IncomingSupply, tc.name)
					case types.Outgoing:
						suite.Equal(assetSupplyPre.OutgoingSupply, assetSupplyPost.OutgoingSupply, tc.name)
					default:
						suite.Fail("should not have invalid direction", tc.name)
					}

					// Check if htlt found in store
					_, found := suite.keeper.GetHTLC(suite.ctx, expectedHTLCID)
					if !tc.shouldBeFound {
						suite.False(found, tc.name)
					} else {
						suite.True(found, tc.name)
					}
				}
			},
		)
	}
}

func (suite *HTLCTestSuite) TestClaimHtlc() {
	suite.SetupTest()
	currentTmTime := tmtime.Now()
	invalidSecret, _ := GenerateRandomSecret()
	type args struct {
		amount    sdk.Coins
		htlcID    tmbytes.HexBytes
		secret    tmbytes.HexBytes
		direction types.SwapDirection
	}
	testCases := []struct {
		name       string
		claimCtx   sdk.Context
		args       args
		expectPass bool
	}{{
		"normal incoming htlc",
		suite.ctx,
		args{
			amount:    cs(c(BNB_DENOM, 50000)),
			htlcID:    []byte{},
			secret:    []byte{},
			direction: types.Incoming,
		},
		true,
	}, {
		"normal incoming htlc rate-limited",
		suite.ctx.WithBlockTime(currentTmTime.Add(time.Minute * 10)),
		args{
			amount:    cs(c(OTHER_DENOM, 50000)),
			htlcID:    []byte{},
			secret:    []byte{},
			direction: types.Incoming,
		},
		true,
	}, {
		"normal outgoing htlc",
		suite.ctx,
		args{
			amount:    cs(c(BNB_DENOM, 50000)),
			htlcID:    []byte{},
			secret:    []byte{},
			direction: types.Outgoing,
		},
		true,
	}, {
		"invalid random number",
		suite.ctx,
		args{
			amount:    cs(c(BNB_DENOM, 50000)),
			htlcID:    []byte{},
			secret:    invalidSecret,
			direction: types.Incoming,
		},
		false,
	}, {
		"wrong htlc ID",
		suite.ctx,
		args{
			amount: cs(c(BNB_DENOM, 50000)),
			htlcID: types.GetID(
				suite.addrs[6],
				suite.addrs[6],
				cs(c(BNB_DENOM, 50000)),
				suite.hashLocks[3],
			),
			secret:    []byte{},
			direction: types.Outgoing,
		},
		false,
	}, {
		"past expiration",
		suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 220),
		args{
			amount:    cs(c(BNB_DENOM, 50000)),
			htlcID:    []byte{},
			secret:    []byte{},
			direction: types.Incoming,
		},
		false,
	}}

	for i, tc := range testCases {
		suite.GenerateHTLCDetails()
		suite.Run(
			tc.name,
			func() {
				expectedRecipient := suite.addrs[5]
				sender := suite.deputy

				// Set sender to other and increment current asset supply for outgoing htlc
				if tc.args.direction == types.Outgoing {
					sender = suite.addrs[6]
					expectedRecipient = suite.deputy
					err := suite.keeper.IncrementCurrentAssetSupply(suite.ctx, tc.args.amount[0])
					suite.Nil(err, tc.name)
				}

				// Create htlc
				id, err := suite.keeper.CreateHTLC(
					suite.ctx,
					sender,
					expectedRecipient,
					ReceiverOnOtherChain,
					SenderOnOtherChain,
					tc.args.amount,
					suite.hashLocks[i],
					suite.timestamps[i],
					MinTimeLock,
					true,
				)
				suite.NoError(err, tc.name)

				realHTLCID := types.GetID(sender, expectedRecipient, tc.args.amount, suite.hashLocks[i])
				suite.Equal(realHTLCID, id, tc.name)

				// If args contains an invalid htlc ID claim attempt will use it instead of the real htlc ID
				var htlcID tmbytes.HexBytes
				if len(tc.args.htlcID) == 0 {
					htlcID = realHTLCID
				} else {
					htlcID = tc.args.htlcID
				}

				// If args contains an invalid random number claim attempt will use it instead of the real random number
				var claimSecret tmbytes.HexBytes
				if len(tc.args.secret) == 0 {
					claimSecret = suite.secrets[i]
				} else {
					claimSecret = tc.args.secret
				}

				// Run the beginblocker before attempting claim
				htlc.BeginBlocker(tc.claimCtx, *suite.keeper)

				expectedRecipientBalancePre := suite.app.BankKeeper.GetBalance(suite.ctx, expectedRecipient, tc.args.amount[0].Denom)
				assetSupplyPre, _ := suite.keeper.GetAssetSupply(tc.claimCtx, tc.args.amount[0].Denom)

				// Attempt to claim htlc
				_, _, _, err = suite.keeper.ClaimHTLC(tc.claimCtx, htlcID, claimSecret)
				expectedRecipientBalancePost := suite.app.BankKeeper.GetBalance(suite.ctx, expectedRecipient, tc.args.amount[0].Denom)
				assetSupplyPost, _ := suite.keeper.GetAssetSupply(tc.claimCtx, tc.args.amount[0].Denom)

				if tc.expectPass {
					suite.NoError(err, tc.name)

					// Check asset supply changes
					switch tc.args.direction {
					case types.Incoming:
						// Check coins moved
						suite.Equal(expectedRecipientBalancePre.Add(tc.args.amount[0]).String(), expectedRecipientBalancePost.String(), tc.name)
						// Check incoming supply decreased
						suite.Equal(assetSupplyPre.IncomingSupply.Amount.Sub(tc.args.amount[0].Amount).String(), assetSupplyPost.IncomingSupply.Amount.String(), tc.name)
						// Check current supply increased
						suite.Equal(assetSupplyPre.CurrentSupply.Add(tc.args.amount[0]).String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						// Check outgoing supply not changed
						suite.Equal(assetSupplyPre.OutgoingSupply.String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					case types.Outgoing:
						// Check incoming supply not changed
						suite.Equal(assetSupplyPre.IncomingSupply.String(), assetSupplyPost.IncomingSupply.String(), tc.name)
						// Check current supply decreased
						suite.Equal(assetSupplyPre.CurrentSupply.Sub(tc.args.amount[0]).String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						// Check outgoing supply decreased
						suite.Equal(assetSupplyPre.OutgoingSupply.Sub(tc.args.amount[0]).String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					default:
						suite.Fail("should not have invalid direction")
					}
				} else {
					suite.Error(err, tc.name)
					// Check coins not moved
					suite.Equal(expectedRecipientBalancePre.String(), expectedRecipientBalancePost.String(), tc.name)

					// Check asset supply has not changed
					switch tc.args.direction {
					case types.Incoming:
						suite.Equal(assetSupplyPre.IncomingSupply.String(), assetSupplyPost.IncomingSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.CurrentSupply.String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.OutgoingSupply.String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					case types.Outgoing:
						suite.Equal(assetSupplyPre.IncomingSupply.String(), assetSupplyPost.IncomingSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.CurrentSupply.String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.OutgoingSupply.String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					default:
						suite.Fail("should not have invalid direction", tc.name)
					}
				}
			},
		)
	}
}

func (suite *HTLCTestSuite) TestRefundHTLC() {
	suite.SetupTest()

	type args struct {
		htlcID    []byte
		direction types.SwapDirection
	}
	testCases := []struct {
		name       string
		refundCtx  sdk.Context
		args       args
		expectPass bool
	}{{
		"normal incoming htlt",
		suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 220),
		args{
			htlcID:    []byte{},
			direction: types.Incoming,
		},
		true,
	}, {
		"normal outgoing htlt",
		suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 220),
		args{
			htlcID:    []byte{},
			direction: types.Outgoing,
		},
		true,
	}, {
		"before expiration",
		suite.ctx,
		args{
			htlcID:    []byte{},
			direction: types.Incoming,
		},
		false,
	}, {
		"wrong htlcID",
		suite.ctx,
		args{
			htlcID: types.GetID(
				suite.addrs[1],
				suite.addrs[1],
				cs(c(BNB_DENOM, 50000)),
				suite.hashLocks[6],
			),
			direction: types.Incoming,
		},
		false,
	}}

	for i, tc := range testCases {
		suite.GenerateHTLCDetails()
		suite.Run(
			tc.name,
			func() {
				// Create atomic swap
				expectedRefundAmount := cs(c(BNB_DENOM, 50000))
				sender := suite.deputy
				expectedRecipient := suite.addrs[9]

				// Set sender to other and increment current asset supply for outgoing swap
				if tc.args.direction == types.Outgoing {
					sender = suite.addrs[6]
					expectedRecipient = suite.deputy
					err := suite.keeper.IncrementCurrentAssetSupply(suite.ctx, expectedRefundAmount[0])
					suite.Nil(err)
				}

				// Create htlc
				id, err := suite.keeper.CreateHTLC(
					suite.ctx,
					sender,
					expectedRecipient,
					ReceiverOnOtherChain,
					SenderOnOtherChain,
					expectedRefundAmount,
					suite.hashLocks[i],
					suite.timestamps[i],
					MinTimeLock,
					true,
				)
				suite.NoError(err, tc.name)

				realHTLCID := types.GetID(sender, expectedRecipient, expectedRefundAmount, suite.hashLocks[i])
				suite.Equal(realHTLCID, id, tc.name)

				originalSenderBalancePre := suite.app.BankKeeper.GetBalance(tc.refundCtx, sender, expectedRefundAmount[0].Denom)
				assetSupplyPre, _ := suite.keeper.GetAssetSupply(tc.refundCtx, expectedRefundAmount[0].Denom)

				// Run the beginblocker before attempting refund
				htlc.BeginBlocker(tc.refundCtx, *suite.keeper)

				originalSenderBalancePost := suite.app.BankKeeper.GetBalance(tc.refundCtx, sender, expectedRefundAmount[0].Denom)
				assetSupplyPost, _ := suite.keeper.GetAssetSupply(tc.refundCtx, expectedRefundAmount[0].Denom)

				if tc.expectPass {
					// Check asset supply changes
					switch tc.args.direction {
					case types.Incoming:
						// Check incoming supply decreased
						suite.Equal(assetSupplyPre.IncomingSupply.Sub(expectedRefundAmount[0]).String(), assetSupplyPost.IncomingSupply.String(), tc.name)
						// Check current, outgoing supply not changed
						suite.Equal(assetSupplyPre.CurrentSupply.String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.OutgoingSupply.String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					case types.Outgoing:
						// Check coins moved
						suite.Equal(originalSenderBalancePre.Add(expectedRefundAmount[0]).String(), originalSenderBalancePost.String(), tc.name)
						// Check incoming, current supply not changed
						suite.Equal(assetSupplyPre.IncomingSupply.String(), assetSupplyPost.IncomingSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.CurrentSupply.String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						// Check outgoing supply decreased
						suite.Equal(assetSupplyPre.OutgoingSupply.Sub(expectedRefundAmount[0]).String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					default:
						suite.Fail("should not have invalid direction")
					}
				} else {
					suite.Equal(originalSenderBalancePre.String(), originalSenderBalancePost.String(), tc.name)

					// Check asset supply has not changed
					switch tc.args.direction {
					case types.Incoming:
						suite.Equal(assetSupplyPre.IncomingSupply.String(), assetSupplyPost.IncomingSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.CurrentSupply.String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.OutgoingSupply.String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					case types.Outgoing:
						suite.Equal(assetSupplyPre.IncomingSupply.String(), assetSupplyPost.IncomingSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.CurrentSupply.String(), assetSupplyPost.CurrentSupply.String(), tc.name)
						suite.Equal(assetSupplyPre.OutgoingSupply.String(), assetSupplyPost.OutgoingSupply.String(), tc.name)
					default:
						suite.Fail("should not have invalid direction", tc.name)
					}
				}
			},
		)
	}
}
