package htlc_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/htlc"
	"github.com/irisnet/irishub/simapp"
)

var (
	senderAddrs   []sdk.AccAddress
	receiverAddrs []sdk.AccAddress

	receiverOnOtherChain string
	amount               sdk.Coins
	secret               []byte
	timestamps           []uint64
	hashLocks            []htlc.HTLCHashLock
	timeLocks            []uint64
	expireHeights        []uint64
	state                htlc.HTLCState
	initSecret           htlc.HTLCSecret

	// construct HTLCs
	htlc1 htlc.HTLC
	htlc2 htlc.HTLC
)

type TestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.app = app

	initVars(suite)
}

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func initVars(suite *TestSuite) {
	senderAddrs = []sdk.AccAddress{sdk.AccAddress("sender1"), sdk.AccAddress("sender2")}
	receiverAddrs = []sdk.AccAddress{sdk.AccAddress("receiver1"), sdk.AccAddress("receiver2")}

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, senderAddrs[0])
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, senderAddrs[1])
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, receiverAddrs[0])
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, receiverAddrs[1])

	_ = suite.app.BankKeeper.SetCoins(suite.ctx, senderAddrs[0], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))
	_ = suite.app.BankKeeper.SetCoins(suite.ctx, senderAddrs[1], sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000)))

	receiverOnOtherChain = "receiverOnOtherChain"
	amount = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000)))
	secret = []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamps = []uint64{uint64(1580000000), 0}
	hashLocks = []htlc.HTLCHashLock{htlc.GetHashLock(secret, timestamps[0]), htlc.GetHashLock(secret, timestamps[1])}
	timeLocks = []uint64{50, 100}
	expireHeights = []uint64{timeLocks[0] + uint64(suite.ctx.BlockHeight()), timeLocks[1] + uint64(suite.ctx.BlockHeight())}
	state = htlc.OPEN
	initSecret = htlc.HTLCSecret{}

	// construct HTLCs
	htlc1 = htlc.NewHTLC(senderAddrs[0], receiverAddrs[0], receiverOnOtherChain, amount, initSecret, timestamps[0], expireHeights[0], state)
	htlc2 = htlc.NewHTLC(senderAddrs[1], receiverAddrs[1], receiverOnOtherChain, amount, initSecret, timestamps[1], expireHeights[1], state)
}

func (suite *TestSuite) TestExportGenesis() {
	// create HTLCs
	err := suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc1, hashLocks[0])
	suite.NoError(err)
	err = suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc2, hashLocks[1])
	suite.NoError(err)

	newBlockHeight := int64(50)
	suite.ctx = suite.ctx.WithBlockHeight(newBlockHeight)
	htlc.BeginBlocker(suite.ctx, suite.app.HTLCKeeper)

	// export genesis
	exportedGenesis := htlc.ExportGenesis(suite.ctx, suite.app.HTLCKeeper)
	exportedHTLCs := exportedGenesis.PendingHTLCs
	suite.Equal(1, len(exportedHTLCs))

	for hashLockHex, tmpHTLC := range exportedHTLCs {
		// assert the state must be OPEN
		suite.True(tmpHTLC.State == htlc.OPEN)

		hashLock, err := hex.DecodeString(hashLockHex)
		suite.NoError(err)

		// assert the HTLC with the given hash lock exists
		htlcInStore, err := suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
		suite.NoError(err)

		// assert the expiration height is new
		newExpireHeight := htlcInStore.ExpireHeight - uint64(newBlockHeight) + 1
		suite.Equal(newExpireHeight, tmpHTLC.ExpireHeight)

		// assert the exported HTLC is consistent with the HTLC in store except for the expiration height
		htlcInStore.ExpireHeight = newExpireHeight
		suite.Equal(htlcInStore, tmpHTLC)
	}

	e := htlc.ValidateGenesis(exportedGenesis)
	suite.NoError(e)

	// assert the expired HTLCs(htlc1) have been refunded
	_, err = suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLocks[0])
	suite.Error(err)
}
