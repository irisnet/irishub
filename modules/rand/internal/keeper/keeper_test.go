package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/rand/internal/keeper"
	"github.com/irisnet/irishub/modules/rand/internal/types"
	"github.com/irisnet/irishub/simapp"
)

// define testing variables
var (
	testTxBytes        = []byte("test-tx")
	testHeight         = int64(10000)
	testBlockInterval  = uint64(100)
	testConsumer       = sdk.AccAddress([]byte("test-consumer"))
	testReqID          = []byte("test-req-id")
	testRandNumerator  = int64(3)
	testRandDenomiator = int64(4)
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.keeper = &app.RandKeeper
}

func (suite *KeeperTestSuite) TestSetRand() {
	rand := types.NewRand(types.SHA256(testTxBytes), testHeight, big.NewRat(testRandNumerator, testRandDenomiator))
	suite.keeper.SetRand(suite.ctx, testReqID, rand)

	storedRand, err := suite.keeper.GetRand(suite.ctx, testReqID)
	suite.NoError(err)
	suite.Equal(rand, storedRand)
}

func (suite *KeeperTestSuite) TestRequestRand() {
	suite.ctx = suite.ctx.WithBlockHeight(testHeight).WithTxBytes(testTxBytes)

	_, err := suite.keeper.RequestRand(suite.ctx, testConsumer, testBlockInterval)
	suite.NoError(err)

	expectedRequest := types.NewRequest(testHeight, testConsumer, types.SHA256(testTxBytes))

	iterator := suite.keeper.IterateRandRequestQueueByHeight(suite.ctx, testHeight+int64(testBlockInterval))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		suite.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), request)

		suite.Equal(expectedRequest, request)
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
