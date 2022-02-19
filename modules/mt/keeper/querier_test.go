package keeper_test

import (
	"encoding/binary"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irismod/modules/mt/exported"
	keep "github.com/irisnet/irismod/modules/mt/keeper"
	"github.com/irisnet/irismod/modules/mt/types"
)

func (suite *KeeperSuite) TestNewQuerier() {
	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err := querier(suite.ctx, []string{"foo", "bar"}, query)
	suite.Error(err)
}

func (suite *KeeperSuite) TestQuerySupply() {
	// MintMT shouldn't fail when collection does not exist
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	query.Path = "/custom/mt/supply"
	query.Data = []byte("?")

	res, err := querier(suite.ctx, []string{"supply"}, query)
	suite.Error(err)
	suite.Nil(res)

	queryCollectionParams := types.NewQuerySupplyParams(denomID2, nil)
	bz, errRes := suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)
	query.Data = bz
	res, err = querier(suite.ctx, []string{"supply"}, query)
	suite.NoError(err)
	supplyResp := binary.LittleEndian.Uint64(res)
	suite.Equal(0, int(supplyResp))

	queryCollectionParams = types.NewQuerySupplyParams(denomID, nil)
	bz, errRes = suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)
	query.Data = bz

	res, err = querier(suite.ctx, []string{"supply"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	supplyResp = binary.LittleEndian.Uint64(res)
	suite.Equal(1, int(supplyResp))
}

func (suite *KeeperSuite) TestQueryCollection() {
	// MintMT shouldn't fail when collection does not exist
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	query.Path = "/custom/mt/collection"

	query.Data = []byte("?")
	res, err := querier(suite.ctx, []string{"collection"}, query)
	suite.Error(err)
	suite.Nil(res)

	queryCollectionParams := types.NewQuerySupplyParams(denomID2, nil)
	bz, errRes := suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)

	query.Data = bz
	_, err = querier(suite.ctx, []string{"collection"}, query)
	suite.NoError(err)

	queryCollectionParams = types.NewQuerySupplyParams(denomID, nil)
	bz, errRes = suite.legacyAmino.MarshalJSON(queryCollectionParams)
	suite.Nil(errRes)

	query.Data = bz
	res, err = querier(suite.ctx, []string{"collection"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	var collection types.Collection
	types.ModuleCdc.MustUnmarshalJSON(res, &collection)
	suite.Len(collection.MTs, 1)
}

func (suite *KeeperSuite) TestQueryOwner() {
	// MintMT shouldn't fail when collection does not exist
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)
	query := abci.RequestQuery{
		Path: "/custom/mt/owner",
		Data: []byte{},
	}

	query.Data = []byte("?")
	_, err = querier(suite.ctx, []string{"owner"}, query)
	suite.Error(err)

	// query the balance using no denomID so that all denoms will be returns
	params := types.NewQuerySupplyParams("", address)
	bz, err2 := suite.legacyAmino.MarshalJSON(params)
	suite.Nil(err2)
	query.Data = bz

	var out types.Owner
	res, err := querier(suite.ctx, []string{"owner"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	suite.legacyAmino.MustUnmarshalJSON(res, &out)

	// build the owner using both denoms
	idCollection1 := types.NewIDCollection(denomID, []string{tokenID})
	idCollection2 := types.NewIDCollection(denomID2, []string{tokenID})
	owner := types.NewOwner(address, idCollection1, idCollection2)

	suite.EqualValues(out.String(), owner.String())
}

func (suite *KeeperSuite) TestQueryMT() {
	// MintMT shouldn't fail when collection does not exist
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	query.Path = "/custom/mt/mt"
	var res []byte

	query.Data = []byte("?")
	res, err = querier(suite.ctx, []string{"mt"}, query)
	suite.Error(err)
	suite.Nil(res)

	params := types.NewQueryMTParams(denomID2, tokenID2)
	bz, err2 := suite.legacyAmino.MarshalJSON(params)
	suite.Nil(err2)

	query.Data = bz
	res, err = querier(suite.ctx, []string{"mt"}, query)
	suite.Error(err)
	suite.Nil(res)

	params = types.NewQueryMTParams(denomID, tokenID)
	bz, err2 = suite.legacyAmino.MarshalJSON(params)
	suite.Nil(err2)

	query.Data = bz
	res, err = querier(suite.ctx, []string{"mt"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	var out exported.MT
	suite.legacyAmino.MustUnmarshalJSON(res, &out)

	suite.Equal(out.GetID(), tokenID)
	suite.Equal(out.GetURI(), tokenURI)
	suite.Equal(out.GetOwner(), address)
}

func (suite *KeeperSuite) TestQueryDenoms() {
	// MintMT shouldn't fail when collection does not exist
	err := suite.keeper.MintMT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	err = suite.keeper.MintMT(suite.ctx, denomID2, tokenID, tokenNm, tokenURI, tokenURIHash, tokenData, address)
	suite.NoError(err)

	querier := keep.NewQuerier(suite.keeper, suite.legacyAmino)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	var res []byte
	query.Path = "/custom/mt/denoms"

	res, err = querier(suite.ctx, []string{"denoms"}, query)
	suite.NoError(err)
	suite.NotNil(res)

	denoms := []string{denomID, denomID2, denomID3}

	var out []types.Denom
	suite.legacyAmino.MustUnmarshalJSON(res, &out)

	for key, denomInQuestion := range out {
		suite.Equal(denomInQuestion.Id, denoms[key])
	}
}
