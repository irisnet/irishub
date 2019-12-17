package htlc_test

import (
	"github.com/irisnet/irishub/modules/htlc"
)

func (suite *KeeperTestSuite) TestBeginBlocker() {
	// create HTLCs
	err := suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc1, hashLocks[0])
	suite.Nil(err)

	newBlockHeight := int64(50)
	suite.ctx = suite.ctx.WithBlockHeight(newBlockHeight)
	htlc.BeginBlocker(suite.ctx, suite.app.HTLCKeeper)

	expectHTLC, err := suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLocks[0])
	suite.Nil(err)
	suite.Equal(expectHTLC.State, htlc.EXPIRED)
}
