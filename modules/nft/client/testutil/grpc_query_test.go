package testutil_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cosmos/gogoproto/proto"
// 	"github.com/stretchr/testify/suite"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	"github.com/cosmos/cosmos-sdk/testutil"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/irisnet/irismod/simapp"
// 	nftcli "github.com/irisnet/irismod/nft/client/cli"
// 	nfttestutil "github.com/irisnet/irismod/nft/client/testutil"
// 	nfttypes "github.com/irisnet/irismod/nft/types"
// )

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	s.network = simapp.SetupNetwork(s.T())
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestNft() {
// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx
// 	// ---------------------------------------------------------------------------

// 	from := val.Address
// 	tokenName := "Kitty Token"
// 	uri := "uri"
// 	uriHash := "uriHash"
// 	description := "description"
// 	data := "{\"key1\":\"value1\",\"key2\":\"value2\"}"
// 	tokenID := "kitty"
// 	//owner     := "owner"
// 	denomName := "name"
// 	denomID := "denom"
// 	schema := "schema"
// 	symbol := "symbol"
// 	mintRestricted := true
// 	updateRestricted := false
// 	baseURL := val.APIAddress

// 	//------test GetCmdIssueDenom()-------------
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", nftcli.FlagDenomName, denomName),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagSymbol, symbol),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagSchema, schema),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagURI, uri),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagURIHash, uriHash),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagDescription, description),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagData, data),
// 		fmt.Sprintf("--%s=%t", nftcli.FlagMintRestricted, mintRestricted),
// 		fmt.Sprintf("--%s=%t", nftcli.FlagUpdateRestricted, updateRestricted),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	expectedCode := uint32(0)

// 	txResult := nfttestutil.IssueDenomExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), denomID, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	//------test GetCmdQueryDenom()-------------
// 	url := fmt.Sprintf("%s/irismod/nft/denoms/%s", baseURL, denomID)
// 	resp, err := testutil.GetRequest(url)
// 	respType := proto.Message(&nfttypes.QueryDenomResponse{})
// 	s.Require().NoError(err)
// 	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
// 	denomItem := respType.(*nfttypes.QueryDenomResponse)
// 	s.Require().Equal(denomName, denomItem.Denom.Name)
// 	s.Require().Equal(schema, denomItem.Denom.Schema)
// 	s.Require().Equal(symbol, denomItem.Denom.Symbol)
// 	s.Require().Equal(uri, denomItem.Denom.Uri)
// 	s.Require().Equal(uriHash, denomItem.Denom.UriHash)
// 	s.Require().Equal(description, denomItem.Denom.Description)
// 	s.Require().Equal(data, denomItem.Denom.Data)
// 	s.Require().Equal(mintRestricted, denomItem.Denom.MintRestricted)
// 	s.Require().Equal(updateRestricted, denomItem.Denom.UpdateRestricted)

// 	//------test GetCmdQueryDenoms()-------------
// 	url = fmt.Sprintf("%s/irismod/nft/denoms", baseURL)
// 	resp, err = testutil.GetRequest(url)
// 	respType = proto.Message(&nfttypes.QueryDenomsResponse{})
// 	s.Require().NoError(err)
// 	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
// 	denomsResp := respType.(*nfttypes.QueryDenomsResponse)
// 	s.Require().Equal(1, len(denomsResp.Denoms))
// 	s.Require().Equal(denomID, denomsResp.Denoms[0].Id)

// 	//------test GetCmdMintNFT()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", nftcli.FlagData, data),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, from.String()),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagURI, uri),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagURIHash, uriHash),
// 		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = nfttestutil.MintNFTExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), denomID, tokenID, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	//------test GetCmdQuerySupply()-------------
// 	url = fmt.Sprintf("%s/irismod/nft/collections/%s/supply", baseURL, denomID)
// 	resp, err = testutil.GetRequest(url)
// 	respType = proto.Message(&nfttypes.QuerySupplyResponse{})
// 	s.Require().NoError(err)
// 	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
// 	supplyResp := respType.(*nfttypes.QuerySupplyResponse)
// 	s.Require().Equal(uint64(1), supplyResp.Amount)

// 	//------test GetCmdQueryNFT()-------------
// 	url = fmt.Sprintf("%s/irismod/nft/nfts/%s/%s", baseURL, denomID, tokenID)
// 	resp, err = testutil.GetRequest(url)
// 	respType = proto.Message(&nfttypes.QueryNFTResponse{})
// 	s.Require().NoError(err)
// 	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
// 	nftItem := respType.(*nfttypes.QueryNFTResponse)
// 	s.Require().Equal(tokenID, nftItem.NFT.Id)
// 	s.Require().Equal(tokenName, nftItem.NFT.Name)
// 	s.Require().Equal(uri, nftItem.NFT.URI)
// 	s.Require().Equal(uriHash, nftItem.NFT.UriHash)
// 	s.Require().Equal(data, nftItem.NFT.Data)
// 	s.Require().Equal(from.String(), nftItem.NFT.Owner)

// 	//------test GetCmdQueryOwner()-------------
// 	url = fmt.Sprintf("%s/irismod/nft/nfts?owner=%s", baseURL, from.String())
// 	resp, err = testutil.GetRequest(url)
// 	respType = proto.Message(&nfttypes.QueryNFTsOfOwnerResponse{})
// 	s.Require().NoError(err)
// 	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
// 	ownerResp := respType.(*nfttypes.QueryNFTsOfOwnerResponse)
// 	s.Require().Equal(from.String(), ownerResp.Owner.Address)
// 	s.Require().Equal(denomID, ownerResp.Owner.IDCollections[0].DenomId)
// 	s.Require().Equal(tokenID, ownerResp.Owner.IDCollections[0].TokenIds[0])

// 	//------test GetCmdQueryCollection()-------------
// 	url = fmt.Sprintf("%s/irismod/nft/collections/%s", baseURL, denomID)
// 	resp, err = testutil.GetRequest(url)
// 	respType = proto.Message(&nfttypes.QueryCollectionResponse{})
// 	s.Require().NoError(err)
// 	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
// 	collectionResp := respType.(*nfttypes.QueryCollectionResponse)
// 	s.Require().Equal(1, len(collectionResp.Collection.NFTs))
// }
