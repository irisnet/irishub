package htlc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub/modules/htlc"
)

func TestABCISuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestBeginBlocker() {
	// create HTLCs
	err := suite.app.HTLCKeeper.CreateHTLC(suite.ctx, htlc1, hashLocks[0])
	suite.NoError(err)

	newBlockHeight := int64(50)
	suite.ctx = suite.ctx.WithBlockHeight(newBlockHeight)
	htlc.BeginBlocker(suite.ctx, suite.app.HTLCKeeper)

	expectHTLC, err := suite.app.HTLCKeeper.GetHTLC(suite.ctx, hashLocks[0])
	suite.NoError(err)
	suite.Equal(expectHTLC.State, htlc.EXPIRED)
}
