package keeper_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/irisnet/irishub/modules/random/keeper"
	"github.com/irisnet/irishub/modules/random/types"
	"github.com/irisnet/irishub/simapp"
)

// define testing variables
var (
	testTxBytes          = []byte("test_tx")
	testHeight           = int64(10000)
	testBlockInterval    = uint64(100)
	testConsumer, _      = sdk.AccAddressFromHex(crypto.AddressHash([]byte("test_consumer")).String())
	testReqID            = []byte("test_req_id")
	testRandomNumerator  = int64(3)
	testRandomDenomiator = int64(4)
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper keeper.Keeper
	app    *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.app = app
	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = app.RandomKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetRandom() {
	rand := types.NewRandom(types.SHA256(testTxBytes), testHeight, big.NewRat(testRandomNumerator, testRandomDenomiator).FloatString(types.RandPrec))
	suite.keeper.SetRandom(suite.ctx, testReqID, rand)

	storedRandom, err := suite.keeper.GetRandom(suite.ctx, testReqID)
	suite.NoError(err)
	randJson, _ := json.Marshal(rand)
	storedRandomJson, _ := json.Marshal(storedRandom)
	suite.Equal(string(randJson), string(storedRandomJson))
}

func (suite *KeeperTestSuite) TestRequestRandom() {
	suite.ctx = suite.ctx.WithBlockHeight(testHeight).WithTxBytes(testTxBytes)

	request, err := suite.keeper.RequestRandom(suite.ctx, testConsumer, testBlockInterval, false, nil)
	suite.NoError(err)

	expectedRequest := types.NewRequest(testHeight, testConsumer, types.SHA256(testTxBytes), false, nil, nil)
	suite.Equal(request, expectedRequest)

	iterator := suite.keeper.IterateRandomRequestQueueByHeight(suite.ctx, testHeight+int64(testBlockInterval))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		suite.app.AppCodec().MustUnmarshalBinaryBare(iterator.Value(), &request)
		suite.Equal(expectedRequest, request)
	}
}
