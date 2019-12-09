package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	addrSender := sdk.AccAddress([]byte("addrSender"))
	addrTo := sdk.AccAddress([]byte("addrTo"))

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrTo)
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrSender, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))
	require.True(suite.T(), suite.app.BankKeeper.GetCoins(suite.ctx, addrSender).IsEqual(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000))))

	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := types.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(suite.ctx.BlockHeight())
	state := types.OPEN
	initSecret := make([]byte, 0)

	_, err := suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	require.NotNil(suite.T(), err)

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
	require.Nil(suite.T(), suite.app.AccountKeeper.GetAccount(suite.ctx, htlcAddr))

	err = suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc, hashLock)
	require.Nil(suite.T(), err)

	htlcAcc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	require.NotNil(suite.T(), htlcAcc)

	amountCreatedHTLC := htlcAcc.GetCoins()
	require.True(suite.T(), amount.IsEqual(amountCreatedHTLC))

	finalSenderAccAmt := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()
	require.True(suite.T(), originSenderAccAmt.Sub(amount).IsEqual(finalSenderAccAmt))

	htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	require.Nil(suite.T(), err)

	require.Equal(suite.T(), addrSender, htlc.Sender)
	require.Equal(suite.T(), addrTo, htlc.To)
	require.Equal(suite.T(), receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(suite.T(), amount, htlc.Amount)
	require.Equal(suite.T(), []byte(nil), htlc.Secret)
	require.Equal(suite.T(), timestamp, htlc.Timestamp)
	require.Equal(suite.T(), expireHeight, htlc.ExpireHeight)
	require.Equal(suite.T(), state, htlc.State)

	store := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	require.True(suite.T(), store.Has(keeper.KeyHTLCExpireQueue(htlc.ExpireHeight, hashLock)))
}

func (suite *KeeperTestSuite) TestClaimHTLC() {
	addrSender := sdk.AccAddress([]byte("addrSender"))
	addrTo := sdk.AccAddress([]byte("addrTo"))

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrTo)
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrSender, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrTo, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 50000)))
	require.True(suite.T(), suite.app.BankKeeper.GetCoins(suite.ctx, addrSender).IsEqual(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000))))

	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret1 := []byte("___abcdefghijklmnopqrstuvwxyz___")
	secret2 := []byte("___00000000000000000000000000___")
	timestamp := uint64(1580000000)
	timestampNil := uint64(0)
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(suite.ctx.BlockHeight())
	state := types.OPEN
	initSecret := make([]byte, 0)

	testData := []struct {
		expectPass           bool
		senderAddr           []byte
		toAddr               []byte
		receiverOnOtherChain string
		amount               sdk.Coins
		secret               []byte
		timestamp            uint64
		hashLock             []byte
		timeLock             uint64
		expireHeight         uint64
		state                types.HTLCState
		initSecret           []byte
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
			require.Nil(suite.T(), err, "TestData: %d", i)

			htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			require.Nil(suite.T(), err, "TestData: %d", i)
			require.Equal(suite.T(), types.OPEN, htlc.State, "TestData: %d", i)

			htlcAcc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

			originHTLCAmount := htlcAcc.GetCoins()
			originReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			_, err = suite.app.HTLCKeeper.ClaimHTLC(suite.ctx, td.hashLock, td.secret)
			require.Nil(suite.T(), err, "TestData: %d", i)

			htlc, _ = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			require.Equal(suite.T(), types.COMPLETED, htlc.State, "TestData: %d", i)

			store := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
			require.True(suite.T(), !store.Has(keeper.KeyHTLCExpireQueue(htlc.ExpireHeight, td.hashLock)))

			htlcAcc = suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

			claimedHTLCAmount := htlcAcc.GetCoins()
			claimedReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			require.True(suite.T(), originHTLCAmount.Sub(amount).IsEqual(claimedHTLCAmount), "TestData: %d", i)
			require.True(suite.T(), originReceiverAmount.Add(amount).IsEqual(claimedReceiverAmount), "TestData: %d", i)
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
			require.Nil(suite.T(), err, "TestData: %d", i)

			htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			require.Nil(suite.T(), err, "TestData: %d", i)
			require.Equal(suite.T(), types.OPEN, htlc.State, "TestData: %d", i)

			htlcAddr := suite.app.SupplyKeeper.GetModuleAddress(types.ModuleName)

			originHTLCAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, htlcAddr).GetCoins()
			originReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			_, err = suite.app.HTLCKeeper.ClaimHTLC(suite.ctx, td.hashLock, td.secret)
			require.NotNil(suite.T(), err, "TestData: %d", i)

			htlc, _ = suite.app.HTLCKeeper.GetHTLC(suite.ctx, td.hashLock)
			require.Equal(suite.T(), types.OPEN, htlc.State, "TestData: %d", i)

			claimedHTLCAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, htlcAddr).GetCoins()
			claimedReceiverAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrTo).GetCoins()

			require.True(suite.T(), originHTLCAmount.IsEqual(claimedHTLCAmount), "TestData: %d", i)
			require.True(suite.T(), originReceiverAmount.IsEqual(claimedReceiverAmount), "TestData: %d", i)
		}
	}
}

func (suite *KeeperTestSuite) TestRefundHTLC() {
	addrSender := sdk.AccAddress([]byte("addrSender"))
	addrTo := sdk.AccAddress([]byte("addrTo"))

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrSender)
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addrTo)
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, addrSender, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))
	require.True(suite.T(), suite.app.BankKeeper.GetCoins(suite.ctx, addrSender).IsEqual(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000))))

	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamp := uint64(1580000000)
	hashLock := types.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock := uint64(50)
	expireHeight := timeLock + uint64(suite.ctx.BlockHeight())
	state := types.EXPIRED
	initSecret := make([]byte, 0)

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
	require.Nil(suite.T(), err)

	htlc, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), types.EXPIRED, htlc.State)

	htlcAcc := suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

	originHTLCAmount := htlcAcc.GetCoins()
	originSenderAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()

	_, err = suite.app.HTLCKeeper.RefundHTLC(suite.ctx, hashLock)
	require.Nil(suite.T(), err)

	htlc, _ = suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
	require.Equal(suite.T(), types.REFUNDED, htlc.State)

	htlcAcc = suite.app.SupplyKeeper.GetModuleAccount(suite.ctx, types.ModuleName)

	claimedHTLCAmount := htlcAcc.GetCoins()
	claimedSenderAmount := suite.app.AccountKeeper.GetAccount(suite.ctx, addrSender).GetCoins()

	require.True(suite.T(), originHTLCAmount.Sub(amount).IsEqual(claimedHTLCAmount))
	require.True(suite.T(), originSenderAmount.Add(amount).IsEqual(claimedSenderAmount))
}
