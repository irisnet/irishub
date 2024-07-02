package keeper_test

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/crypto"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"mods.irisnet.org/modules/random/keeper"
	"mods.irisnet.org/modules/random/types"
	"mods.irisnet.org/simapp"
)

// define testing variables
var (
	testTxBytes       = []byte("test_tx")
	testHeight        = int64(10000)
	testBlockInterval = uint64(100)
	testConsumer, _   = sdk.AccAddressFromHexUnsafe(
		crypto.AddressHash([]byte("test_consumer")).String(),
	)
	testReqID            = []byte("test_req_id")
	testRandomNumerator  = int64(3)
	testRandomDenomiator = int64(4)
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.LegacyAmino
	ctx    sdk.Context
	keeper keeper.Keeper
	app    *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	depInjectOptions := simapp.DepinjectOptions{
		Config:    AppConfig,
		Providers: []interface{}{},
		Consumers: []interface{}{&suite.keeper},
	}

	app := simapp.Setup(suite.T(), false, depInjectOptions)

	suite.app = app
	suite.cdc = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetRandom() {
	random := types.NewRandom(
		hex.EncodeToString(types.SHA256(testTxBytes)),
		testHeight,
		big.NewRat(testRandomNumerator, testRandomDenomiator).FloatString(types.RandPrec),
	)
	suite.keeper.SetRandom(suite.ctx, testReqID, random)

	storedRandom, err := suite.keeper.GetRandom(suite.ctx, testReqID)
	suite.NoError(err)
	randJson, _ := json.Marshal(random)
	storedRandomJson, _ := json.Marshal(storedRandom)
	suite.Equal(string(randJson), string(storedRandomJson))
}

func (suite *KeeperTestSuite) TestRequestRandom() {
	suite.ctx = suite.ctx.WithBlockHeight(testHeight).WithTxBytes(testTxBytes)

	request, err := suite.keeper.RequestRandom(
		suite.ctx,
		testConsumer,
		testBlockInterval,
		false,
		nil,
	)
	suite.NoError(err)

	expectedRequest := types.NewRequest(
		testHeight,
		testConsumer.String(),
		hex.EncodeToString(types.SHA256(testTxBytes)),
		false,
		nil,
		"",
	)
	suite.Equal(request, expectedRequest)

	iterator := suite.keeper.IterateRandomRequestQueueByHeight(
		suite.ctx,
		testHeight+int64(testBlockInterval),
	)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		suite.app.AppCodec().MustUnmarshal(iterator.Value(), &request)
		suite.Equal(expectedRequest, request)
	}
}
