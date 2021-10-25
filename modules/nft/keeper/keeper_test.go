package keeper_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft/keeper"
	"github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/simapp"
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

	address   = CreateTestAddrs(1)[0]
	address2  = CreateTestAddrs(2)[1]
	address3  = CreateTestAddrs(3)[2]
	tokenURI  = "https://google.com/token-1.json"
	tokenURI2 = "https://google.com/token-2.json"
	tokenData = "{a:a,b:b}"

	isCheckTx = false
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

	app := simapp.Setup(isCheckTx)

	suite.app = app
	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.NFTKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.NFTKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	err := suite.keeper.IssueDenom(suite.ctx, denomID, denomNm, schema, denomSymbol, address, false, false)
	suite.NoError(err)

	// MintNFT shouldn't fail when collection does not exist
	err = suite.keeper.IssueDenom(suite.ctx, denomID2, denomNm2, schema, denomSymbol2, address, false, false)
	suite.NoError(err)

	err = suite.keeper.IssueDenom(suite.ctx, denomID3, denomNm3, schema, denomSymbol3, address3, true, true)
	suite.NoError(err)

	// collections should equal 3
	collections := suite.keeper.GetCollections(suite.ctx)
	suite.NotEmpty(collections)
	suite.Equal(len(collections), 3)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperSuite))
}

func (suite *KeeperSuite) TestMintNFT() {
	// MintNFT shouldn't fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// MintNFT shouldn't fail when collection exists
	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenData, address)
	suite.NoError(err)

}

func (suite *KeeperSuite) TestUpdateNFT() {
	// EditNFT should fail when NFT doesn't exists
	err := suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm3, tokenURI, tokenData, address)
	suite.Error(err)

	// MintNFT shouldn't fail when collection does not exist
	err = suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// EditNFT should fail when NFT doesn't exists
	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID2, tokenNm2, tokenURI, tokenData, address)
	suite.Error(err)

	// EditNFT shouldn't fail when NFT exists
	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI2, tokenData, address)
	suite.NoError(err)

	// EditNFT should fail when NFT failed to authorize
	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI2, tokenData, address2)
	suite.Error(err)

	// GetNFT should get the NFT with new tokenURI
	receivedNFT, err := suite.keeper.GetNFT(suite.ctx, denomID, tokenID)
	suite.NoError(err)
	suite.Equal(receivedNFT.GetURI(), tokenURI2)

	// EditNFT shouldn't fail when NFT exists
	err = suite.keeper.EditNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI2, tokenData, address2)
	suite.Error(err)

	err = suite.keeper.MintNFT(suite.ctx, denomID3, denomID3, tokenID3, tokenURI, tokenData, address3)
	suite.NoError(err)

	// EditNFT should fail if updateRestricted equal to true, nobody can update the NFT under this denom
	err = suite.keeper.EditNFT(suite.ctx, denomID3, denomID3, tokenID3, tokenURI, tokenData, address3)
	suite.Error(err)
}

func (suite *KeeperSuite) TestTransferOwner() {

	// MintNFT shouldn't fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// invalid owner
	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address2, address3)
	suite.Error(err)

	// right
	err = suite.keeper.TransferOwner(suite.ctx, denomID, tokenID, tokenNm2, tokenURI2, tokenData, address, address2)
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

	denom, _ := suite.keeper.GetDenom(suite.ctx, denomID)

	// denom.Creator should equal to address3 after transfer
	suite.Equal(denom.Creator, address3.String())
}

func (suite *KeeperSuite) TestBurnNFT() {
	// MintNFT should not fail when collection does not exist
	err := suite.keeper.MintNFT(suite.ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, address)
	suite.NoError(err)

	// BurnNFT should fail when NFT doesn't exist but collection does exist
	err = suite.keeper.BurnNFT(suite.ctx, denomID, tokenID, address)
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
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, testAddr(buffer.String(), bech))
		buffer.Reset()
	}

	return addresses
}

// for incode address generation
func testAddr(addr string, bech string) sdk.AccAddress {
	res, err := sdk.AccAddressFromHex(addr)
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
