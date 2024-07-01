package keeper_test

import (
	"mods.irisnet.org/modules/nft/keeper"
	"mods.irisnet.org/modules/nft/types"
)

func (suite *KeeperSuite) TestSetCollection() {
	nft := types.NewBaseNFT(tokenID, tokenNm, address, tokenURI, tokenURIHash, tokenData)
	// create a new NFT and add it to the collection created with the NFT mint
	nft2 := types.NewBaseNFT(tokenID2, tokenNm, address, tokenURI, tokenURIHash2, tokenData)

	denomE := types.Denom{
		Id:               denomID,
		Name:             denomNm,
		Schema:           schema,
		Creator:          address.String(),
		Symbol:           denomSymbol,
		MintRestricted:   true,
		UpdateRestricted: true,
	}

	collection2 := types.Collection{
		Denom: denomE,
		NFTs:  []types.BaseNFT{nft2, nft},
	}

	err := suite.keeper.SaveCollection(suite.ctx, collection2)
	suite.Nil(err)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *KeeperSuite) TestGetCollections() {
	// SaveNFT shouldn't fail when collection does not exist
	err := suite.keeper.SaveNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

func (suite *KeeperSuite) TestGetSupply() {
	// SaveNFT shouldn't fail when collection does not exist
	err := suite.keeper.SaveNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	// SaveNFT shouldn't fail when collection does not exist
	err = suite.keeper.SaveNFT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenURIHash, tokenData, address2)
	suite.NoError(err)

	// SaveNFT shouldn't fail when collection does not exist
	err = suite.keeper.SaveNFT(suite.ctx, denomID2, tokenID, tokenNm2, tokenURI, tokenURIHash, tokenData, address2)
	suite.NoError(err)

	supply := suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(2), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID2)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetBalance(suite.ctx, denomID, address)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetBalance(suite.ctx, denomID, address2)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(2), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID2)
	suite.Equal(uint64(1), supply)

	// burn nft
	err = suite.keeper.RemoveNFT(suite.ctx, denomID, tokenID, address)
	suite.NoError(err)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(1), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(1), supply)

	// burn nft
	err = suite.keeper.RemoveNFT(suite.ctx, denomID, tokenID2, address2)
	suite.NoError(err)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(0), supply)

	supply = suite.keeper.GetTotalSupply(suite.ctx, denomID)
	suite.Equal(uint64(0), supply)
}
