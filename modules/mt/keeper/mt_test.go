package keeper_test

import (
	"github.com/irisnet/irismod/modules/mt/keeper"
)

func (suite *KeeperSuite) TestGetMT() {
	// MintMT shouldn't fail when collection does not exist
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	// GetMT should get the MT
	receivedMT, err := suite.keeper.GetMT(suite.ctx, denomID, tokenID)
	suite.NoError(err)
	suite.Equal(receivedMT.GetID(), tokenID)
	suite.True(receivedMT.GetOwner().Equals(address))
	suite.Equal(receivedMT.GetURI(), tokenURI)

	// MintMT shouldn't fail when collection exists
	err = suite.keeper.MintMT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	// GetMT should get the MT when collection exists
	receivedMT2, err := suite.keeper.GetMT(suite.ctx, denomID, tokenID2)
	suite.NoError(err)
	suite.Equal(receivedMT2.GetID(), tokenID2)
	suite.True(receivedMT2.GetOwner().Equals(address))
	suite.Equal(receivedMT2.GetURI(), tokenURI)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *KeeperSuite) TestGetMTs() {
	err := suite.keeper.MintMT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID2, tokenID2, tokenNm2, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID2, tokenID3, tokenNm3, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID, tokenID3, tokenNm3, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	mts := suite.keeper.GetMTs(suite.ctx, denomID2)
	suite.Len(mts, 3)
}

func (suite *KeeperSuite) TestAuthorize() {
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	_, err = suite.keeper.Authorize(suite.ctx, denomID, tokenID, address2)
	suite.Error(err)

	_, err = suite.keeper.Authorize(suite.ctx, denomID, tokenID, address)
	suite.NoError(err)
}

func (suite *KeeperSuite) TestHasMT() {
	// IsMT should return false
	isMT := suite.keeper.HasMT(suite.ctx, denomID, tokenID)
	suite.False(isMT)

	// MintMT shouldn't fail when collection does not exist
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	// IsMT should return true
	isMT = suite.keeper.HasMT(suite.ctx, denomID, tokenID)
	suite.True(isMT)
}
