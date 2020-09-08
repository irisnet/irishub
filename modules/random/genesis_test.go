package random_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/random"
	"github.com/irisnet/irishub/modules/random/keeper"
	"github.com/irisnet/irishub/modules/random/types"
	"github.com/irisnet/irishub/simapp"
)

// define testing variables
var (
	testTxBytes        = []byte("test-tx")
	testHeight         = int64(10000)
	testNewHeight      = testHeight + 50
	testBlockInterval1 = uint64(100)
	testBlockInterval2 = uint64(200)
	testConsumer1      = sdk.AccAddress("test-consumer1")
	testConsumer2      = sdk.AccAddress("test-consumer2")
)

type GenesisTestSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *GenesisTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = app.RandomKeeper
}

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestExportGenesis() {
	suite.ctx = suite.ctx.WithBlockHeight(testHeight).WithTxBytes(testTxBytes)

	// request rands
	_, err := suite.keeper.RequestRandom(suite.ctx, testConsumer1, testBlockInterval1, false, sdk.NewCoins())
	suite.NoError(err)
	_, err = suite.keeper.RequestRandom(suite.ctx, testConsumer2, testBlockInterval2, false, sdk.NewCoins())
	suite.NoError(err)

	// precede to the new block
	suite.ctx = suite.ctx.WithBlockHeight(testNewHeight)

	// get the pending requests from queue
	storedRequests := make(map[int64][]types.Request)
	suite.keeper.IterateRandomRequestQueue(suite.ctx, func(h int64, r types.Request) bool {
		storedRequests[h] = append(storedRequests[h], r)
		return false
	})
	suite.Equal(2, len(storedRequests))

	// export genesis
	genesis := random.ExportGenesis(suite.ctx, suite.keeper)
	exportedRequests := genesis.PendingRandomRequests
	suite.Equal(2, len(exportedRequests))

	// assert that exported requests are consistent with requests in queue
	for height, requests := range exportedRequests {
		h, _ := strconv.ParseInt(height, 10, 64)
		storedHeight := h + testNewHeight - 1
		suite.Equal(storedRequests[storedHeight], requests.Requests)
	}
}
