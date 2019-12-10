package htlc_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	modulehtlc "github.com/irisnet/irishub/modules/htlc"
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

func (suite *KeeperTestSuite) TestExportHTLCGenesis() {
	// define variables
	senderAddrs := []sdk.AccAddress{sdk.AccAddress([]byte("sender1")), sdk.AccAddress([]byte("sender2"))}
	receiverAddrs := []sdk.AccAddress{sdk.AccAddress([]byte("receiver1")), sdk.AccAddress([]byte("receiver2"))}

	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, senderAddrs[0])
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, senderAddrs[1])
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, receiverAddrs[0])
	_ = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, receiverAddrs[1])

	// _ = suite.app.BankKeeper.SetCoins(suite.ctx, senderAddrs[0], sdk.NewCoins(sdk.NewInt64Coin(config.Iris, 100000)))

	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)))
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamps := []uint64{uint64(1580000000), 0}
	hashLocks := []types.HTLCHashLock{types.GetHashLock(secret, timestamps[0]), types.GetHashLock(secret, timestamps[1])}
	timeLocks := []uint64{50, 100}
	expireHeights := []uint64{timeLocks[0] + uint64(suite.ctx.BlockHeight()), timeLocks[1] + uint64(suite.ctx.BlockHeight())}
	state := modulehtlc.OPEN
	initSecret := make(types.HTLCSecret, 0)

	// construct HTLCs
	htlc1 := modulehtlc.NewHTLC(senderAddrs[0], receiverAddrs[0], receiverOnOtherChain, amount, initSecret, timestamps[0], expireHeights[0], state)
	htlc2 := modulehtlc.NewHTLC(senderAddrs[1], receiverAddrs[1], receiverOnOtherChain, amount, initSecret, timestamps[1], expireHeights[1], state)

	// create HTLCs
	err := suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc1, hashLocks[0])
	require.Nil(suite.T(), err)
	err = suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc2, hashLocks[1])
	require.Nil(suite.T(), err)

	newBlockHeight := int64(50)
	suite.ctx = suite.ctx.WithBlockHeight(newBlockHeight)
	modulehtlc.BeginBlocker(suite.ctx, suite.app.HTLCKeeper)

	// export genesis
	exportedGenesis := modulehtlc.ExportGenesis(suite.ctx, suite.app.HTLCKeeper)
	exportedHTLCs := exportedGenesis.PendingHTLCs
	require.Equal(suite.T(), 1, len(exportedHTLCs))

	for hashLockHex, htlc := range exportedHTLCs {
		// assert the state must be OPEN
		require.True(suite.T(), htlc.State == modulehtlc.OPEN)

		hashLock, err := hex.DecodeString(hashLockHex)

		// assert the HTLC with the given hash lock exists
		htlcInStore, err := suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLock)
		require.Nil(suite.T(), err)

		// assert the expiration height is new
		newExpireHeight := htlcInStore.ExpireHeight - uint64(newBlockHeight) + 1
		require.Equal(suite.T(), newExpireHeight, htlc.ExpireHeight)

		// assert the exported HTLC is consistant with the HTLC in store except for the expiration height
		htlcInStore.ExpireHeight = newExpireHeight
		require.Equal(suite.T(), htlcInStore, htlc)
	}

	// assert the expired HTLCs(htlc1) have been refunded
	htlc, _ := suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLocks[0])
	require.Equal(suite.T(), modulehtlc.REFUNDED, htlc.State)
}
