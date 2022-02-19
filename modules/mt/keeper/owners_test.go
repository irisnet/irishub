package keeper_test

import (
	"github.com/irisnet/irismod/modules/mt/keeper"
)

func (suite *KeeperSuite) TestGetOwners() {

	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenURIHash, tokenData, address2)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID, tokenID3, tokenNm3, tokenURI, tokenURIHash, tokenData, address3)
	suite.NoError(err)

	owners := suite.keeper.GetOwners(suite.ctx)
	suite.Equal(3, len(owners))

	err = suite.keeper.MintMT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID2, tokenID2, tokenNm2, tokenURI, tokenURIHash, tokenData, address2)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID2, tokenID3, tokenNm3, tokenURI, tokenURIHash, tokenData, address3)
	suite.NoError(err)

	owners = suite.keeper.GetOwners(suite.ctx)
	suite.Equal(3, len(owners))

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}
