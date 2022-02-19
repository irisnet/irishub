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

	mtcli "github.com/irisnet/irismod/modules/mt/client/cli"
	mttestutil "github.com/irisnet/irismod/modules/mt/client/testutil"
	mttypes "github.com/irisnet/irismod/modules/mt/types"
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
	uri := "uri"
	uriHash := "uriHash"
	description := "description"
	data := "data"
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
		fmt.Sprintf("--%s=%s", mtcli.FlagName, denomName),
		fmt.Sprintf("--%s=%s", mtcli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", mtcli.FlagSchema, schema),
		fmt.Sprintf("--%s=%s", mtcli.FlagURI, uri),
		fmt.Sprintf("--%s=%s", mtcli.FlagURIHash, uriHash),
		fmt.Sprintf("--%s=%s", mtcli.FlagDescription, description),
		fmt.Sprintf("--%s=%s", mtcli.FlagData, data),
		fmt.Sprintf("--%s=%t", mtcli.FlagMintRestricted, mintRestricted),
		fmt.Sprintf("--%s=%t", mtcli.FlagUpdateRestricted, updateRestricted),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := mttestutil.IssueDenomExec(val.ClientCtx, from.String(), denom, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	denomID := gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()

	//------test GetCmdQueryDenom()-------------
	url := fmt.Sprintf("%s/irismod/mt/denoms/%s", baseURL, denomID)
	resp, err := rest.GetRequest(url)
	respType = proto.Message(&mttypes.QueryDenomResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	denomItem := respType.(*mttypes.QueryDenomResponse)
	s.Require().Equal(denomName, denomItem.Denom.Name)
	s.Require().Equal(schema, denomItem.Denom.Schema)
	s.Require().Equal(symbol, denomItem.Denom.Symbol)
	s.Require().Equal(uri, denomItem.Denom.Uri)
	s.Require().Equal(uriHash, denomItem.Denom.UriHash)
	s.Require().Equal(description, denomItem.Denom.Description)
	s.Require().Equal(data, denomItem.Denom.Data)
	s.Require().Equal(mintRestricted, denomItem.Denom.MintRestricted)
	s.Require().Equal(updateRestricted, denomItem.Denom.UpdateRestricted)

	//------test GetCmdQueryDenoms()-------------
	url = fmt.Sprintf("%s/irismod/mt/denoms", baseURL)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&mttypes.QueryDenomsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	denomsResp := respType.(*mttypes.QueryDenomsResponse)
	s.Require().Equal(1, len(denomsResp.Denoms))
	s.Require().Equal(denomID, denomsResp.Denoms[0].Id)

	//------test GetCmdMintMT()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagData, data),
		fmt.Sprintf("--%s=%s", mtcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", mtcli.FlagURI, uri),
		fmt.Sprintf("--%s=%s", mtcli.FlagURIHash, uriHash),
		fmt.Sprintf("--%s=%s", mtcli.FlagTokenName, tokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = mttestutil.MintMTExec(val.ClientCtx, from.String(), denomID, tokenID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	//------test GetCmdQuerySupply()-------------
	url = fmt.Sprintf("%s/irismod/mt/collections/%s/supply", baseURL, denomID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&mttypes.QuerySupplyResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	supplyResp := respType.(*mttypes.QuerySupplyResponse)
	s.Require().Equal(uint64(1), supplyResp.Amount)

	//------test GetCmdQueryMT()-------------
	url = fmt.Sprintf("%s/irismod/mt/mts/%s/%s", baseURL, denomID, tokenID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&mttypes.QueryMTResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	mtItem := respType.(*mttypes.QueryMTResponse)
	s.Require().Equal(tokenID, mtItem.MT.Id)
	s.Require().Equal(tokenName, mtItem.MT.Name)
	s.Require().Equal(uri, mtItem.MT.URI)
	s.Require().Equal(uriHash, mtItem.MT.UriHash)
	s.Require().Equal(data, mtItem.MT.Data)
	s.Require().Equal(from.String(), mtItem.MT.Owner)

	//------test GetCmdQueryOwner()-------------
	url = fmt.Sprintf("%s/irismod/mt/mts?owner=%s", baseURL, from.String())
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&mttypes.QueryOwnerResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	ownerResp := respType.(*mttypes.QueryOwnerResponse)
	s.Require().Equal(from.String(), ownerResp.Owner.Address)
	s.Require().Equal(denom, ownerResp.Owner.IDCollections[0].DenomId)
	s.Require().Equal(tokenID, ownerResp.Owner.IDCollections[0].TokenIds[0])

	//------test GetCmdQueryCollection()-------------
	url = fmt.Sprintf("%s/irismod/mt/collections/%s", baseURL, denomID)
	resp, err = rest.GetRequest(url)
	respType = proto.Message(&mttypes.QueryCollectionResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	collectionResp := respType.(*mttypes.QueryCollectionResponse)
	s.Require().Equal(1, len(collectionResp.Collection.MTs))

	//------test GetCmdTransferDenom()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = mttestutil.TransferDenomExec(val.ClientCtx, from.String(), recipient.String(), denomID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&mttypes.Denom{})
	bz, err = mttestutil.QueryDenomExec(val.ClientCtx, denomID)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	denomItem2 := respType.(*mttypes.Denom)
	s.Require().Equal(recipient.String(), denomItem2.Creator)
	s.Require().Equal(denomName, denomItem2.Name)
	s.Require().Equal(schema, denomItem2.Schema)
	s.Require().Equal(symbol, denomItem2.Symbol)
	s.Require().Equal(mintRestricted, denomItem2.MintRestricted)
	s.Require().Equal(updateRestricted, denomItem2.UpdateRestricted)
}
