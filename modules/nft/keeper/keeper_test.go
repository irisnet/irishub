package keeper_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"mods.irisnet.org/modules/nft/keeper"
	"mods.irisnet.org/modules/nft/types"
	"mods.irisnet.org/simapp"
)

var (
	denomID     = "denomid"
	denomNm     = "denomnm"
	denomSymbol = "denomSymbol"
	schema      = "{a:a,b:b}"

	denomID2     = "denomid2"
	denomNm2     = "denom2nm"
	denomSymbol2 = "denomSymbol2"

	tokenID  = "tokenid"
	tokenID2 = "tokenid2"
	tokenID3 = "tokenid3"

	tokenNm  = "tokennm"
	tokenNm2 = "tokennm2"
	tokenNm3 = "tokennm3"

	denomID3     = "denomid3"
	denomNm3     = "denom3nm"
	denomSymbol3 = "denomSymbol3"

	address          = CreateTestAddrs(1)[0]
	address2         = CreateTestAddrs(2)[1]
	address3         = CreateTestAddrs(3)[2]
	tokenURI         = "https://google.com/token-1.json"
	tokenURIHash     = "tokenURIHash"
	tokenURI2        = "https://google.com/token-2.json"
	tokenURIHash2    = "tokenURIHash2"
	tokenData        = "{a:a,b:b}"
	denomDescription = "this is a class name of a nft"
	denomUri         = "denom uri"
	denomUriHash     = "denom uri hash"
	denomData        = "{\"key1\":\"value1\",\"key2\":\"value2\"}"
	isCheckTx        = false
)

type KeeperSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	app         *simapp.SimApp

	queryClient types.QueryClient
}

func (suite *KeeperSuite) SetupTest() {
	depInjectOptions := simapp.DepinjectOptions{
		Config:    AppConfig,
		Providers: []interface{}{},
		Consumers: []interface{}{&suite.keeper},
	}

	app := simapp.Setup(suite.T(), isCheckTx, depInjectOptions)

	suite.app = app
	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	err := suite.keeper.SaveDenom(
		suite.ctx,
		denomID,
		denomNm,
		schema,
		denomSymbol,
		address,
		false,
		false,
		denomDescription,
		denomUri,
		denomUriHash,
		denomData,
	)
	suite.NoError(err)

	// SaveNFT shouldn't fail when collection does not exist
	err = suite.keeper.SaveDenom(
		suite.ctx,
		denomID2,
		denomNm2,
		schema,
		denomSymbol2,
		address,
		false,
		false,
		denomDescription,
		denomUri,
		denomUriHash,
		denomData,
	)
	suite.NoError(err)

	err = suite.keeper.SaveDenom(
		suite.ctx,
		denomID3,
		denomNm3,
		schema,
		denomSymbol3,
		address3,
		true,
		true,
		denomDescription,
		denomUri,
		denomUriHash,
		denomData,
	)
	suite.NoError(err)

	// collections should equal 3
	collections, err := suite.keeper.GetCollections(suite.ctx)
	suite.NoError(err)
	suite.NotEmpty(collections)
	suite.Equal(len(collections), 3)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperSuite))
}

func (suite *KeeperSuite) TestMintNFT() {
	// SaveNFT shouldn't fail when collection does not exist
	err := suite.keeper.SaveNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI,
		tokenURIHash,
		tokenData,
		address,
	)
	suite.NoError(err)

	// SaveNFT shouldn't fail when collection exists
	err = suite.keeper.SaveNFT(
		suite.ctx,
		denomID,
		tokenID2,
		tokenNm2,
		tokenURI,
		tokenURIHash,
		tokenData,
		address,
	)
	suite.NoError(err)
}

func (suite *KeeperSuite) TestUpdateNFT() {
	// UpdateNFT should fail when NFT doesn't exists
	err := suite.keeper.UpdateNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm3,
		tokenURI,
		tokenURIHash,
		tokenData,
		address,
	)
	suite.Error(err)

	// SaveNFT shouldn't fail when collection does not exist
	err = suite.keeper.SaveNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI,
		tokenURIHash,
		tokenData,
		address,
	)
	suite.NoError(err)

	// UpdateNFT should fail when NFT doesn't exists
	err = suite.keeper.UpdateNFT(
		suite.ctx,
		denomID,
		tokenID2,
		tokenNm2,
		tokenURI,
		tokenURIHash,
		tokenData,
		address,
	)
	suite.Error(err)

	// UpdateNFT shouldn't fail when NFT exists
	err = suite.keeper.UpdateNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI2,
		tokenURIHash2,
		tokenData,
		address,
	)
	suite.NoError(err)

	// UpdateNFT should fail when NFT failed to authorize
	err = suite.keeper.UpdateNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI2,
		tokenURIHash2,
		tokenData,
		address2,
	)
	suite.Error(err)

	// GetNFT should get the NFT with new tokenURI
	receivedNFT, err := suite.keeper.GetNFT(suite.ctx, denomID, tokenID)
	suite.NoError(err)
	suite.Equal(receivedNFT.GetURI(), tokenURI2)

	// UpdateNFT shouldn't fail when NFT exists
	err = suite.keeper.UpdateNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI2,
		tokenURIHash2,
		tokenData,
		address2,
	)
	suite.Error(err)

	err = suite.keeper.SaveNFT(
		suite.ctx,
		denomID3,
		denomID3,
		tokenID3,
		tokenURI,
		tokenURIHash,
		tokenData,
		address3,
	)
	suite.NoError(err)

	// UpdateNFT should fail if updateRestricted equal to true, nobody can update the NFT under this denom
	err = suite.keeper.UpdateNFT(
		suite.ctx,
		denomID3,
		denomID3,
		tokenID3,
		tokenURI,
		tokenURIHash,
		tokenData,
		address3,
	)
	suite.Error(err)
}

func (suite *KeeperSuite) TestTransferOwnership() {
	// SaveNFT shouldn't fail when collection does not exist
	err := suite.keeper.SaveNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI,
		tokenURIHash,
		tokenData,
		address,
	)
	suite.NoError(err)

	// invalid owner
	err = suite.keeper.TransferOwnership(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI,
		tokenURIHash,
		tokenData,
		address2,
		address3,
	)
	suite.Error(err)

	// right
	err = suite.keeper.TransferOwnership(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm2,
		tokenURI2,
		tokenURIHash2,
		tokenData,
		address,
		address2,
	)
	suite.NoError(err)

	nft, err := suite.keeper.GetNFT(suite.ctx, denomID, tokenID)
	suite.NoError(err)
	suite.Equal(tokenURI2, nft.GetURI())
}

func (suite *KeeperSuite) TestTransferDenom() {
	// invalid owner
	err := suite.keeper.TransferDenomOwner(suite.ctx, denomID, address3, address)
	suite.Error(err)

	// right
	err = suite.keeper.TransferDenomOwner(suite.ctx, denomID, address, address3)
	suite.NoError(err)

	denom, _ := suite.keeper.GetDenomInfo(suite.ctx, denomID)

	// denom.Creator should equal to address3 after transfer
	suite.Equal(denom.Creator, address3.String())
}

func (suite *KeeperSuite) TestBurnNFT() {
	// SaveNFT should not fail when collection does not exist
	err := suite.keeper.SaveNFT(
		suite.ctx,
		denomID,
		tokenID,
		tokenNm,
		tokenURI,
		tokenURIHash,
		tokenData,
		address,
	)
	suite.NoError(err)

	// RemoveNFT should fail when NFT doesn't exist but collection does exist
	err = suite.keeper.RemoveNFT(suite.ctx, denomID, tokenID, address)
	suite.NoError(err)

	// NFT should no longer exist
	isNFT := suite.keeper.HasNFT(suite.ctx, denomID, tokenID)
	suite.False(isNFT)

	msg, fail := keeper.SupplyInvariant(suite.keeper)(suite.ctx)
	suite.False(fail, msg)
}

// CreateTestAddrs creates test addresses
func CreateTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string

		buffer.WriteString(numString) // adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
		bech := res.String()
		addresses = append(addresses, testAddr(buffer.String(), bech))
		buffer.Reset()
	}

	return addresses
}

// for incode address generation
func testAddr(addr, bech string) sdk.AccAddress {
	res, err := sdk.AccAddressFromHexUnsafe(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res
}
