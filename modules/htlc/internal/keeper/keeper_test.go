package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/htlc/internal/keeper"
	"github.com/irisnet/irishub/modules/htlc/internal/types"
	"github.com/irisnet/irishub/simapp"
)

type KeeperTestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.app = app
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestCreateHTLC() {
	addrSender := sdk.AccAddress("addrSender")
	addrTo := sdk.AccAddress("addrTo")

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrTo)
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrSender, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))
	suite.True(suite.app.BankKeeper.GetCoins(suite.ctx, addrSender).IsEqual(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000))))

	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret := types.HTLCSecret("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := types.HTLCHashLock(types.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(suite.ctx.BlockHeight())
	state := types.OPEN
	initSecret := make(types.HTLCSecret, 0)

	_, err := suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	suite.NotNil(err)

	htlc := types.NewHTLC(
		addrSender,
		addrTo,
		receiverOnOtherChain,
		amount,
		initSecret,
		timestamp,
		expireHeight,
		state,
	)

	originSenderAccAmt := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()

	htlcAddr := suite.app.SupplyKeeper.GetModuleAddress(types.ModuleName)
	suite.Nil(suite.app.AccountKeeper.GetAccount(suite.ctx, htlcAddr))

	err = suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc, hashLock)
	suite.Nil(err)

	htlcAcc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	suite.NotNil(htlcAcc)

	amountCreatedHTLC := htlcAcc.GetCoins()
	suite.True(amount.IsEqual(amountCreatedHTLC))

	finalSenderAccAmt := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()
	suite.True(originSenderAccAmt.Sub(amount).IsEqual(finalSenderAccAmt))

	htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	suite.Nil(err)

	suite.Equal(addrSender, htlc.Sender)
	suite.Equal(addrTo, htlc.To)
	suite.Equal(receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	suite.Equal(amount, htlc.Amount)
	suite.Equal(types.HTLCSecret(nil), htlc.Secret)
	suite.Equal(timestamp, htlc.Timestamp)
	suite.Equal(expireHeight, htlc.ExpireHeight)
	suite.Equal(state, htlc.State)

	store := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	suite.True(store.Has(keeper.KeyHTLCExpireQueue(htlc.ExpireHeight, hashLock)))
}

func (suite *KeeperTestSuite) TestClaimHTLC() {
	addrSender := sdk.AccAddress("addrSender")
	addrTo := sdk.AccAddress("addrTo")

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrTo)
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrSender, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrTo, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 50000)))
	suite.True(suite.app.BankKeeper.GetCoins(suite.ctx, addrSender).IsEqual(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000))))

	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret1 := types.HTLCSecret("___abcdefghijklmnopqrstuvwxyz___")
	secret2 := types.HTLCSecret("___00000000000000000000000000___")
	timestamp := uint64(1580000000)
	timestampNil := uint64(0)
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(suite.ctx.BlockHeight())
	state := types.OPEN
	initSecret := make(types.HTLCSecret, 0)

	testData := []struct {
		expectPass           bool
		senderAddr           []byte
		toAddr               []byte
		receiverOnOtherChain string
		amount               sdk.Coins
		secret               types.HTLCSecret
		timestamp            uint64
		hashLock             types.HTLCHashLock
		timeLock             uint64
		expireHeight         uint64
		state                types.HTLCState
		initSecret           types.HTLCSecret
	}{
		// timestamp > 0
		{true, addrSender, addrTo, receiverOnOtherChain, amount, secret1, timestamp, types.GetHashLock(secret1, timestamp), timeLock, expireHeight, state, initSecret},
		// timestamp = 0
		{true, addrSender, addrTo, receiverOnOtherChain, amount, secret1, timestampNil, types.GetHashLock(secret1, timestampNil), timeLock, expireHeight, state, initSecret},
		// invalid secret
		{false, addrSender, addrTo, receiverOnOtherChain, amount, secret1, timestampNil, types.GetHashLock(secret2, timestampNil), timeLock, expireHeight, state, initSecret},
	}

	for i, td := range testData {
		if td.expectPass {
			htlc := types.NewHTLC(
				td.senderAddr,
				td.toAddr,
				td.receiverOnOtherChain,
				td.amount,
				td.initSecret,
				td.timestamp,
				td.expireHeight,
				td.state,
			)

			err := suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc, td.hashLock)
			suite.Nil(err, "TestData: %d", i)

			htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			suite.Nil(err, "TestData: %d", i)
			suite.Equal(types.OPEN, htlc.State, "TestData: %d", i)

			htlcAcc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

			originHTLCAmount := htlcAcc.GetCoins()
			originReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			_, err = suite.app.HTLCKeeper.ClaimHTLC(suite.ctx, td.hashLock, td.secret)
			suite.Nil(err, "TestData: %d", i)

			htlc, _ = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			suite.Equal(types.COMPLETED, htlc.State, "TestData: %d", i)

			store := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
			suite.True(!store.Has(keeper.KeyHTLCExpireQueue(htlc.ExpireHeight, td.hashLock)))

			htlcAcc = suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

			claimedHTLCAmount := htlcAcc.GetCoins()
			claimedReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			suite.True(originHTLCAmount.Sub(amount).IsEqual(claimedHTLCAmount), "TestData: %d", i)
			suite.True(originReceiverAmount.Add(amount).IsEqual(claimedReceiverAmount), "TestData: %d", i)
		} else {
			htlc := types.NewHTLC(
				td.senderAddr,
				td.toAddr,
				td.receiverOnOtherChain,
				td.amount,
				td.initSecret,
				td.timestamp,
				td.expireHeight,
				td.state,
			)

			err := suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc, td.hashLock)
			suite.Nil(err, "TestData: %d", i)

			htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			suite.Nil(err, "TestData: %d", i)
			suite.Equal(types.OPEN, htlc.State, "TestData: %d", i)

			htlcAddr := suite.app.SupplyKeeper.GetModuleAddress(types.ModuleName)

			originHTLCAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, htlcAddr).GetCoins()
			originReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			_, err = suite.app.HTLCKeeper.ClaimHTLC(suite.ctx, td.hashLock, td.secret)
			suite.NotNil(err, "TestData: %d", i)

			htlc, _ = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			suite.Equal(types.OPEN, htlc.State, "TestData: %d", i)

			claimedHTLCAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, htlcAddr).GetCoins()
			claimedReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			suite.True(originHTLCAmount.IsEqual(claimedHTLCAmount), "TestData: %d", i)
			suite.True(originReceiverAmount.IsEqual(claimedReceiverAmount), "TestData: %d", i)
		}
	}
}

func (suite *KeeperTestSuite) TestRefundHTLC() {
	addrSender := sdk.AccAddress("addrSender")
	addrTo := sdk.AccAddress("addrTo")

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrTo)
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrSender, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))
	suite.True(suite.app.BankKeeper.GetCoins(suite.ctx, addrSender).IsEqual(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000))))

	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret := types.HTLCSecret("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := types.HTLCHashLock(types.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(suite.ctx.BlockHeight())
	state := types.EXPIRED
	initSecret := make(types.HTLCSecret, 0)

	htlc := types.NewHTLC(
		addrSender,
		addrTo,
		receiverOnOtherChain,
		amount,
		initSecret,
		timestamp,
		expireHeight,
		state,
	)

	err := suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc, hashLock)
	suite.Nil(err)

	htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	suite.Nil(err)
	suite.Equal(types.EXPIRED, htlc.State)

	htlcAcc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

	originHTLCAmount := htlcAcc.GetCoins()
	originSenderAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()

	_, err = suite.app.HTLCKeeper.RefundHTLC(suite.ctx, hashLock)
	suite.Nil(err)

	htlc, _ = suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	suite.Equal(types.REFUNDED, htlc.State)

	htlcAcc = suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

	claimedHTLCAmount := htlcAcc.GetCoins()
	claimedSenderAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()

	suite.True(originHTLCAmount.Sub(amount).IsEqual(claimedHTLCAmount))
	suite.True(originSenderAmount.Add(amount).IsEqual(claimedSenderAmount))
}
