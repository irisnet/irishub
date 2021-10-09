package rest_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	nftcli "github.com/irisnet/irismod/modules/nft/client/cli"
	nfttestutil "github.com/irisnet/irismod/modules/nft/client/testutil"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 2

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestNft() {
	val := s.network.Validators[0]
	recipient := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))
	// ---------------------------------------------------------------------------

	from := val.Address
	tokenName := "Kitty Token"
	tokenURI := "uri"
	tokenData := "data"
	tokenID := "kitty"
	//owner     := "owner"
	denomName := "name"
	denom := "denom"
	schema := "schema"
	symbol := "symbol"
	mintRestricted := true
	updateRestricted := false
	baseURL := val.APIAddress

	//------test GetCmdIssueDenom()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagDenomName, denomName),
		fmt.Sprintf("--%s=%s", nftcli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", nftcli.FlagSchema, schema),
		fmt.Sprintf("--%s=%t", nftcli.FlagMintRestricted, mintRestricted),
		fmt.Sprintf("--%s=%t", nftcli.FlagUpdateRestricted, updateRestricted),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := nfttestutil.IssueDenomExec(val.ClientCtx, from.String(), denom, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	denomID := gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()

	//------test GetCmdQueryDenom()-------------
	url := fmt.Sprintf("%s/irismod/nft/denoms/%s", baseURL, denomID)
	resp, err := rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryDenomResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	denomItem := respType.(*nfttypes.QueryDenomResponse)
	s.Require().Equal(denomName, denomItem.Denom.Name)
	s.Require().Equal(schema, denomItem.Denom.Schema)
	s.Require().Equal(symbol, denomItem.Denom.Symbol)
	s.Require().Equal(mintRestricted, denomItem.Denom.MintRestricted)
	s.Require().Equal(updateRestricted, denomItem.Denom.UpdateRestricted)

	//------test GetCmdQueryDenoms()-------------
	url = fmt.Sprintf("%s/irismod/nft/denoms", baseURL)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryDenomsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	denomsResp := respType.(*nfttypes.QueryDenomsResponse)
	s.Require().Equal(1, len(denomsResp.Denoms))
	s.Require().Equal(denomID, denomsResp.Denoms[0].Id)

	//------test GetCmdMintNFT()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
		fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, tokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.MintNFTExec(val.ClientCtx, from.String(), denomID, tokenID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQuerySupply()-------------
	url = fmt.Sprintf("%s/irismod/nft/collections/%s/supply", baseURL, denomID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QuerySupplyResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	supplyResp := respType.(*nfttypes.QuerySupplyResponse)
	s.Require().Equal(uint64(1), supplyResp.Amount)

	//------test GetCmdQueryNFT()-------------
	url = fmt.Sprintf("%s/irismod/nft/nfts/%s/%s", baseURL, denomID, tokenID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryNFTResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	nftItem := respType.(*nfttypes.QueryNFTResponse)
	s.Require().Equal(tokenID, nftItem.NFT.Id)
	s.Require().Equal(tokenName, nftItem.NFT.Name)
	s.Require().Equal(tokenURI, nftItem.NFT.URI)
	s.Require().Equal(tokenData, nftItem.NFT.Data)
	s.Require().Equal(from.String(), nftItem.NFT.Owner)

	//------test GetCmdQueryOwner()-------------
	url = fmt.Sprintf("%s/irismod/nft/nfts?owner=%s", baseURL, from.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryOwnerResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	ownerResp := respType.(*nfttypes.QueryOwnerResponse)
	s.Require().Equal(from.String(), ownerResp.Owner.Address)
	s.Require().Equal(denom, ownerResp.Owner.IDCollections[0].DenomId)
	s.Require().Equal(tokenID, ownerResp.Owner.IDCollections[0].TokenIds[0])

	//------test GetCmdQueryCollection()-------------
	url = fmt.Sprintf("%s/irismod/nft/collections/%s", baseURL, denomID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&nfttypes.QueryCollectionResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	collectionResp := respType.(*nfttypes.QueryCollectionResponse)
	s.Require().Equal(1, len(collectionResp.Collection.NFTs))

	//------test GetCmdTransferDenom()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = nfttestutil.TransferDenomExec(val.ClientCtx, from.String(), recipient.String(), denomID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&nfttypes.Denom{})
	bz, err = nfttestutil.QueryDenomExec(val.ClientCtx, denomID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	denomItem2 := respType.(*nfttypes.Denom)
	s.Require().Equal(recipient.String(), denomItem2.Creator)
	s.Require().Equal(denomName, denomItem2.Name)
	s.Require().Equal(schema, denomItem2.Schema)
	s.Require().Equal(symbol, denomItem2.Symbol)
	s.Require().Equal(mintRestricted, denomItem2.MintRestricted)
	s.Require().Equal(updateRestricted, denomItem2.UpdateRestricted)
}
