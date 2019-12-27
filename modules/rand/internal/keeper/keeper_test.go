package keeper_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/rand"
	"github.com/irisnet/irishub/simapp"
)

// define testing variables
var (
	testTxBytes        = []byte("test_tx")
	testHeight         = int64(10000)
	testBlockInterval  = uint64(100)
	testConsumer, _    = sdk.AccAddressFromHex(crypto.AddressHash([]byte("test_consumer")).String())
	testReqID          = []byte("test_req_id")
	testRandNumerator  = int64(3)
	testRandDenomiator = int64(4)
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper rand.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.keeper = app.RandKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetRand() {
	rand := rand.NewRand(rand.SHA256(testTxBytes), testHeight, big.NewRat(testRandNumerator, testRandDenomiator).FloatString(rand.RandPrec))
	suite.keeper.SetRand(suite.ctx, testReqID, rand)

	storedRand, err := suite.keeper.GetRand(suite.ctx, testReqID)
	suite.NoError(err)
	randJson, _ := json.Marshal(rand)
	storedRandJson, _ := json.Marshal(storedRand)
	suite.Equal(string(randJson), string(storedRandJson))
}

func (suite *KeeperTestSuite) TestRequestRand() {
	suite.ctx = suite.ctx.WithBlockHeight(testHeight).WithTxBytes(testTxBytes)

	request, err := suite.keeper.RequestRand(suite.ctx, testConsumer, testBlockInterval)
	suite.NoError(err)

	expectedRequest := rand.NewRequest(testHeight, testConsumer, rand.SHA256(testTxBytes))
	suite.Equal(request, expectedRequest)

	iterator := suite.keeper.IterateRandRequestQueueByHeight(suite.ctx, testHeight+int64(testBlockInterval))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request rand.Request
		suite.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)
		suite.Equal(expectedRequest, request)
	}
}
