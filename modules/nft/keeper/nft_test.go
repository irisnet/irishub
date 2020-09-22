package keeper_test

import (
	"github.com/irisnet/irismod/modules/nft/keeper"
)

func (suite *KeeperSuite) TestGetNFT() {
	// MintNFT shouldn't fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// GetNFT should get the NFT
	receivedNFT, err := suite.keeper.GetNFT(suite.ctx, denomID, tokenID)
	suite.NoError(err)
	suite.Equal(receivedNFT.GetID(), tokenID)
	suite.True(receivedNFT.GetOwner().Equals(address))
	suite.Equal(receivedNFT.GetURI(), tokenURI)

	// MintNFT shouldn't fail when collection exists
	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenData, address)
	suite.NoError(err)

	// GetNFT should get the NFT when collection exists
	receivedNFT2, err := suite.keeper.GetNFT(suite.ctx, denomID, tokenID2)
	suite.NoError(err)
	suite.Equal(receivedNFT2.GetID(), tokenID2)
	suite.True(receivedNFT2.GetOwner().Equals(address))
	suite.Equal(receivedNFT2.GetURI(), tokenURI)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *KeeperSuite) TestGetNFTs() {
	err := suite.keeper.MintNFT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID2, tokenNm2, tokenURI, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID2, tokenID3, tokenNm3, tokenURI, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID3, tokenNm3, tokenURI, tokenData, address)
	suite.NoError(err)

	nfts := suite.keeper.GetNFTs(suite.ctx, denomID2)
	suite.Len(nfts, 3)
}

func (suite *KeeperSuite) TestAuthorize() {
	err := suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	_, err = suite.keeper.Authorize(suite.ctx, denomID, tokenID, address2)
	suite.Error(err)

	_, err = suite.keeper.Authorize(suite.ctx, denomID, tokenID, address)
	suite.NoError(err)
}

func (suite *KeeperSuite) TestHasNFT() {
	// IsNFT should return false
	isNFT := suite.keeper.HasNFT(suite.ctx, denomID, tokenID)
	suite.False(isNFT)

	// MintNFT shouldn't fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// IsNFT should return true
	isNFT = suite.keeper.HasNFT(suite.ctx, denomID, tokenID)
	suite.True(isNFT)
}
