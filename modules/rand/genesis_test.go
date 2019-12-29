package rand_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/rand"
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

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper rand.Keeper
}

func (suite *GenesisTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.keeper = app.RandKeeper
}

func TestQuerierSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestExportGenesis() {
	suite.ctx = suite.ctx.WithBlockHeight(testHeight).WithTxBytes(testTxBytes)

	// request rands
	_, err := suite.keeper.RequestRand(suite.ctx, testConsumer1, testBlockInterval1)
	suite.NoError(err)
	_, err = suite.keeper.RequestRand(suite.ctx, testConsumer2, testBlockInterval2)
	suite.NoError(err)

	// precede to the new block
	suite.ctx = suite.ctx.WithBlockHeight(testNewHeight)

	// get the pending requests from queue
	storedRequests := make(map[int64][]rand.Request)
	suite.keeper.IterateRandRequestQueue(suite.ctx, func(h int64, r rand.Request) bool {
		storedRequests[h] = append(storedRequests[h], r)
		return false
	})
	suite.Equal(2, len(storedRequests))

	// export genesis
	genesis := rand.ExportGenesis(suite.ctx, suite.keeper)
	exportedRequests := genesis.PendingRandRequests
	suite.Equal(2, len(exportedRequests))

	// assert that exported requests are consistent with requests in queue
	for height, requests := range exportedRequests {
		h, _ := strconv.ParseInt(height, 10, 64)
		storedHeight := h + testNewHeight - 1
		suite.Equal(storedRequests[storedHeight], requests)
	}
}
